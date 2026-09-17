package config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Schedule is a weekly class timetable. When Enabled, the daily pages
// draw the hour grid as a fixed-height column with one block per
// timetable slot of that weekday (see tpls/schedule_classes.tpl):
//
//   - a shaded block with the group name for each Class,
//   - an empty outlined box for each Period without a class (the teacher
//     is at school but has no lesson), and
//   - a "break" block for periods flagged as Break (e.g. the recess).
type Schedule struct {
	Enabled bool
	// Days are the weekdays that get the period boxes (default Mon-Fri).
	// Weekdays not listed here only show classes explicitly assigned to
	// them, which normally means an empty grid (weekends, holidays).
	Days []time.Weekday
	// Periods is the daily frame shared by all school days: the lesson
	// slots and the breaks, in any order.
	Periods []Period
	Classes []Class
}

// Period is one slot of the daily frame.
type Period struct {
	// Start and End are "HH:MM" 24h times.
	Start string
	End   string
	// Name is printed inside the box when there is no class in it. For a
	// lesson slot it is usually empty; for a break it is e.g. "RECREO".
	Name string
	// Break marks a non-teaching slot (recess). It is drawn differently
	// and never gets a class matched to it.
	Break bool
}

// Class is one recurring lesson in the weekly timetable.
type Class struct {
	// Day is the weekday the class happens on: 1 = Monday ... 5 = Friday
	// (time.Weekday numbering, so 0 = Sunday and 6 = Saturday also work).
	Day time.Weekday
	// Start and End are "HH:MM" 24h times, e.g. "08:30" and "09:20".
	Start string
	End   string
	// Name is the short label printed inside the block.
	Name string
	// Group is the full group name (optional; kept for later linking the
	// block to the group's student roster page).
	Group string
}

// Block kinds, as seen by the template.
const (
	BlockClass = "class"
	BlockFree  = "free"
	BlockBreak = "break"
)

// ScheduleBlock is one box on the daily hour grid. Top and Bottom are the
// slot start/end expressed in hours since the grid's bottom hour (e.g.
// with a grid starting at 08:00, a 08:30-09:20 slot has Top 0.5 and
// Bottom 1.333). They are consumed by the TikZ template, which uses the
// hour height as its y unit.
type ScheduleBlock struct {
	Kind   string
	Name   string
	Group  string
	Start  string
	End    string
	Top    float64
	Bottom float64
}

// Mid is the vertical centre of the block, in the same units as Top/Bottom.
func (b ScheduleBlock) Mid() float64 {
	return (b.Top + b.Bottom) / 2
}

// ForWeekday returns the blocks to draw on the given weekday, sorted by
// start time and positioned relative to bottomHour: every class of that
// day, plus (on school days) one free/break box per period that no class
// occupies. Blocks entirely outside [bottomHour, topHour+1) are dropped;
// partial overlaps are clipped so a block never leaves the grid.
func (s Schedule) ForWeekday(wd time.Weekday, bottomHour, topHour int) ([]ScheduleBlock, error) {
	gridTop := float64(topHour + 1 - bottomHour)
	blocks := make([]ScheduleBlock, 0, len(s.Periods)+len(s.Classes))

	add := func(kind, name, group, start, end string) error {
		from, err := parseHHMM(start)
		if err != nil {
			return fmt.Errorf("%s %q start: %w", kind, name, err)
		}

		to, err := parseHHMM(end)
		if err != nil {
			return fmt.Errorf("%s %q end: %w", kind, name, err)
		}

		if to <= from {
			return fmt.Errorf("%s %q: end %s is not after start %s", kind, name, end, start)
		}

		top := from - float64(bottomHour)
		bottom := to - float64(bottomHour)

		if bottom <= 0 || top >= gridTop {
			return nil
		}

		if top < 0 {
			top = 0
		}

		if bottom > gridTop {
			bottom = gridTop
		}

		blocks = append(blocks, ScheduleBlock{Kind: kind, Name: name, Group: group, Start: start, End: end, Top: top, Bottom: bottom})

		return nil
	}

	for _, c := range s.Classes {
		if c.Day != wd {
			continue
		}

		if err := add(BlockClass, c.Name, c.Group, c.Start, c.End); err != nil {
			return nil, err
		}
	}

	if s.isSchoolDay(wd) {
		for _, p := range s.Periods {
			kind := BlockFree
			if p.Break {
				kind = BlockBreak
			}

			if kind == BlockFree && s.periodTaken(wd, p) {
				continue
			}

			if err := add(kind, p.Name, "", p.Start, p.End); err != nil {
				return nil, err
			}
		}
	}

	sort.SliceStable(blocks, func(i, j int) bool { return blocks[i].Top < blocks[j].Top })

	return blocks, nil
}

// isSchoolDay reports whether wd gets the period boxes: the days listed
// in Days, or Monday to Friday when Days is empty.
func (s Schedule) isSchoolDay(wd time.Weekday) bool {
	if len(s.Days) == 0 {
		return wd >= time.Monday && wd <= time.Friday
	}

	for _, d := range s.Days {
		if d == wd {
			return true
		}
	}

	return false
}

// periodTaken reports whether some class on wd overlaps the period.
func (s Schedule) periodTaken(wd time.Weekday, p Period) bool {
	pFrom, err1 := parseHHMM(p.Start)
	pTo, err2 := parseHHMM(p.End)

	if err1 != nil || err2 != nil {
		return false
	}

	for _, c := range s.Classes {
		if c.Day != wd {
			continue
		}

		cFrom, err1 := parseHHMM(c.Start)
		cTo, err2 := parseHHMM(c.End)

		if err1 != nil || err2 != nil {
			continue
		}

		if cFrom < pTo && cTo > pFrom {
			return true
		}
	}

	return false
}

// parseHHMM converts "HH:MM" into fractional hours (e.g. "09:20" -> 9.333).
func parseHHMM(hhmm string) (float64, error) {
	parts := strings.Split(strings.TrimSpace(hhmm), ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("want HH:MM, got %q", hhmm)
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil || h < 0 || h > 23 {
		return 0, fmt.Errorf("bad hour in %q", hhmm)
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 0 || m > 59 {
		return 0, fmt.Errorf("bad minute in %q", hhmm)
	}

	return float64(h) + float64(m)/60, nil
}
