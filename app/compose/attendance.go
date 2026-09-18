package compose

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kudrykv/latex-yearly-planner/app/components/cal"
	"github.com/kudrykv/latex-yearly-planner/app/components/page"
	"github.com/kudrykv/latex-yearly-planner/app/config"
	"github.com/kudrykv/latex-yearly-planner/app/tex"
)

// attendanceMinRows is the number of roster rows always drawn on an
// attendance page, even when a group's list in cfg/teacher_students.yaml
// has fewer names: the rest are printed as blank, ruled rows to fill in
// by hand.
const attendanceMinRows = 25

// Attendance builds one page per (teaching group, calendar month) that
// actually has at least one class session, in the order the groups
// first appear in cfg/teacher_schedule.yaml and chronologically within
// each group. Produces no pages at all when the schedule has no groups
// (e.g. a classic, non-teacher config).
func Attendance(cfg config.Config, tpls []string) (page.Modules, error) {
	modules := make(page.Modules, 0)

	groups := cfg.Schedule.Groups()
	if len(groups) == 0 {
		return modules, nil
	}

	rosters := cfg.Students.ByGroup()

	from := time.Date(cfg.Year, cfg.StartMonth, 1, 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, cfg.NumMonths, -1)

	byDay, err := cfg.Events.ByDay()
	if err != nil {
		return nil, fmt.Errorf("events: %w", err)
	}

	isHoliday := func(t time.Time) bool {
		for _, ev := range byDay[config.DayKey(t)] {
			if ev.Holiday {
				return true
			}
		}

		return false
	}

	year := cal.NewYear(cfg.WeekStart, cfg.Year, cfg.StartMonth, cfg.NumMonths)

	for _, group := range groups {
		dates, err := cfg.Schedule.SessionDates(group, from, to, isHoliday)
		if err != nil {
			return nil, fmt.Errorf("session dates %q: %w", group, err)
		}

		byMonth := map[string][]time.Time{}

		var monthKeys []string

		for _, d := range dates {
			key := d.Format("2006-01")
			if _, ok := byMonth[key]; !ok {
				monthKeys = append(monthKeys, key)
			}

			byMonth[key] = append(byMonth[key], d)
		}

		for _, key := range config.SortedMonthKeys(monthKeys) {
			monthDates := byMonth[key]
			first := monthDates[0]

			anchor := config.AttendanceAnchor(group, first)
			short := cfg.Schedule.ShortName(group)

			table := attendanceTable(monthDates, rosters[group])
			heading := attendanceHeading(anchor, short, first)

			modules = append(modules, page.Module{
				Cfg: cfg,
				Tpl: tpls[0],
				Body: map[string]interface{}{
					"Year":         year,
					"SideQuarters": year.SideQuarters(cal.Day{Time: first}.Quarter()),
					"SideMonths":   year.SideMonths(first.Month()),
					"Extra2":       extra2(cfg.ClearTopRightCorner, false, false, nil, 0),
					"HeadingMOS":   heading,
					"Table":        table,
				},
			})
		}
	}

	return modules, nil
}

// attendanceHeading is the page title: a hypertarget bound to this
// group+month (the daily schedule's "Attendance" link jumps here), the
// group's short name, and the month/year.
func attendanceHeading(anchor, groupShort string, first time.Time) string {
	label := `{\Large\textbf{` + groupShort + `}}\hfill{\large Attendance ` +
		first.Month().String() + " " + strconv.Itoa(first.Year()) + `}`

	return tex.Hypertarget(anchor, "") + label
}

// attendanceTable builds the roster/attendance grid: a name column, then
// one column per session date. Row 1 is the rotated, clickable date
// header (linking back to that day's page); row 2 is a single "Roll
// called" checkbox per date, to tell "nobody was marked absent" apart
// from "attendance wasn't taken"; the rest are one blank, ruled row per
// student (pad to attendanceMinRows when the roster is shorter).
func attendanceTable(dates []time.Time, students []string) string {
	n := len(dates)
	if n == 0 {
		return ""
	}

	rows := attendanceMinRows
	if len(students) > rows {
		rows = len(students)
	}

	nameCol := `\myLenAttNameCol`
	dateColSpec := `|>{\centering\arraybackslash}p{\myLenAttCol}`
	colFmt := fmt.Sprintf(`|p{%s}`, nameCol) + repeat(dateColSpec, n) + `|`

	// Column width: the (n+1) columns must add up to \linewidth, but a
	// plain \dimexpr(\linewidth-nameCol)/n ignores what tabular itself
	// adds around every column (2*\tabcolsep of padding, plus the n+2
	// vertical rules), which on a table with up to ~19 columns is enough
	// to overflow the page. \tabcolsep is set to a small fixed value
	// first (scoped to this table only) so both computations agree.
	const colSep = `1pt`

	var b string

	b += `{\setlength{\tabcolsep}{` + colSep + `}%` + "\n"
	b += fmt.Sprintf("\\setlength{\\myLenAttCol}{\\dimexpr(\\linewidth-%s-%d\\tabcolsep-%dpt)/%d\\relax}%%\n",
		nameCol, 2*(n+1), n+2, n)
	b += `\renewcommand{\arraystretch}{1.6}` + "\n"
	b += `\begin{tabular}{` + colFmt + "}\n\\hline\n"

	// Row 1: rotated, clickable date headers (weekday name + day number,
	// e.g. "Thursday 6" -> translated to "Jueves 6").
	b += ` `
	for _, d := range dates {
		label := d.Weekday().String() + " " + strconv.Itoa(d.Day())
		b += ` & \rotatebox{45}{\small ` + tex.Hyperlink(d.Format(time.RFC3339), label) + `}`
	}

	b += "\\\\\n\\hline\n"

	// Row 2: one "roll called" cell per date (blank; the cell's own
	// border is the checkbox, no glyph needed inside).
	b += `Roll called`

	for range dates {
		b += ` & `
	}

	b += "\\\\\n\\hline\n"

	// One row per student (padded with blank rows up to attendanceMinRows).
	for i := 0; i < rows; i++ {
		name := ""
		if i < len(students) {
			name = students[i]
		}

		b += name

		for range dates {
			b += ` & `
		}

		b += "\\\\\n\\hline\n"
	}

	b += `\end{tabular}}` + "\n"

	return b
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}

	return out
}
