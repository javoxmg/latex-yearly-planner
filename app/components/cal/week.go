package cal

import (
	"strconv"
	"strings"
	"time"

	"github.com/kudrykv/latex-yearly-planner/app/components/header"
	"github.com/kudrykv/latex-yearly-planner/app/components/hyper"
	"github.com/kudrykv/latex-yearly-planner/app/tex"
)

type Weeks []*Week
type Week struct {
	Days [7]Day

	Weekday  time.Weekday
	Year     *Year
	Months   Months
	Quarters Quarters
}

func NewWeeksForMonth(wd time.Weekday, year *Year, qrtr *Quarter, month *Month) Weeks {
	ptr := time.Date(month.CalendarYear, month.Month, 1, 0, 0, 0, 0, time.Local)
	weekday := ptr.Weekday()
	shift := (7 + weekday - wd) % 7

	week := &Week{Weekday: wd, Year: year, Months: Months{month}, Quarters: Quarters{qrtr}}

	for i := shift; i < 7; i++ {
		week.Days[i] = Day{Time: ptr}
		ptr = ptr.AddDate(0, 0, 1)
	}

	weeks := Weeks{}
	weeks = append(weeks, week)

	for ptr.Month() == month.Month {
		week = &Week{Weekday: weekday, Year: year, Months: Months{month}, Quarters: Quarters{qrtr}}

		for i := 0; i < 7; i++ {
			if ptr.Month() != month.Month {
				break
			}

			week.Days[i] = Day{ptr}
			ptr = ptr.AddDate(0, 0, 1)
		}

		weeks = append(weeks, week)
	}

	return weeks
}

func NewWeeksForYear(wd time.Weekday, year *Year) Weeks {
	ptr := selectStartWeek(wd)

	firstCalYear, firstMonth := monthAt(0)
	qrtr1 := NewQuarter(wd, year, quarterNumber(firstCalYear, firstMonth))
	mon1 := NewMonth(wd, year, qrtr1, firstMonth, firstCalYear)
	week := &Week{Weekday: wd, Year: year, Quarters: Quarters{qrtr1}, Months: Months{mon1}}
	weeks := make(Weeks, 0, 53)

	for i := 0; i < 7; i++ {
		week.Days[i] = ptr
		ptr = ptr.Add(1)
	}

	weeks = append(weeks, week)

	for !ptr.Time.After(rangeEnd) {
		weeks = append(weeks, fillWeekly(wd, year, ptr))
		ptr = ptr.Add(7)
	}

	weeks[len(weeks)-1].Quarters = weeks[len(weeks)-1].Quarters[:1]
	weeks[len(weeks)-1].Months = weeks[len(weeks)-1].Months[:1]

	return weeks
}

func fillWeekly(wd time.Weekday, year *Year, ptr Day) *Week {
	qrtr := NewQuarter(wd, year, quarterNumber(ptr.Time.Year(), ptr.Time.Month()))
	month := NewMonth(wd, year, qrtr, ptr.Time.Month(), ptr.Time.Year())

	week := &Week{Weekday: wd, Year: year, Quarters: Quarters{qrtr}, Months: Months{month}}

	for i := 0; i < 7; i++ {
		week.Days[i] = ptr
		ptr = ptr.Add(1)
	}

	if week.quarterOverlap() {
		qrtr = NewQuarter(wd, year, week.rightQuarter())
		week.Quarters = append(week.Quarters, qrtr)
	}

	if week.monthOverlap() {
		month = NewMonth(wd, year, qrtr, week.rightMonth(), week.rightYear())
		week.Months = append(week.Months, month)
	}

	return week
}

func selectStartWeek(weekStart time.Weekday) Day {
	sow := rangeStart

	for sow.Weekday() != weekStart {
		sow = sow.AddDate(0, 0, 1)
	}

	if sow.After(rangeStart) {
		sow = sow.AddDate(0, 0, -7)
	}

	return Day{Time: sow}
}

func (w *Week) WeekNumber(large interface{}) string {
	wn := w.weekNumber()
	larg, _ := large.(bool)

	itoa := strconv.Itoa(wn)
	ref := w.ref()
	if !larg {
		return hyper.Link(ref, itoa)
	}

	text := `\rotatebox[origin=tr]{90}{\makebox[\myLenMonthlyCellHeight][c]{Week ` + itoa + `}}`

	return hyper.Link(ref, text)
}

func (w *Week) weekNumber() int {
	_, wn := w.Days[0].Time.ISOWeek()

	for _, t := range w.Days {
		if _, cwn := t.Time.ISOWeek(); !t.Time.IsZero() && cwn != wn {
			return cwn
		}
	}

	return wn
}

func (w *Week) Breadcrumb() string {
	return header.Items{
		header.NewIntItem(w.Year.Number),
		w.QuartersBreadcrumb(),
		w.MonthsBreadcrumb(),
		header.NewTextItem("Week " + strconv.Itoa(w.weekNumber())).RefText(w.ref()).Ref(true),
	}.Table(true)
}

func (w *Week) monthOverlap() bool {
	return w.Days[0].Time.Month() != w.Days[6].Time.Month()
}

func (w *Week) quarterOverlap() bool {
	return w.leftQuarter() != w.rightQuarter()
}

func (w *Week) leftQuarter() int {
	return quarterNumber(w.leftYear(), w.leftMonth())
}

func (w *Week) rightQuarter() int {
	return quarterNumber(w.rightYear(), w.rightMonth())
}

func (w *Week) rightMonth() time.Month {
	for i := 6; i >= 0; i-- {
		if w.Days[i].Time.IsZero() {
			continue
		}

		return w.Days[i].Time.Month()
	}

	return -1
}

func (w *Week) PrevNext() header.Items {
	items := header.Items{}

	if w.PrevExists() {
		wn := w.Prev().weekNumber()
		items = append(items, header.NewTextItem("Week "+strconv.Itoa(wn)))
	}

	if w.NextExists() {
		wn := w.Next().weekNumber()
		items = append(items, header.NewTextItem("Week "+strconv.Itoa(wn)))
	}

	return items
}

func (w *Week) NextExists() bool {
	stillInRange := !w.Days[6].Time.After(rangeEnd)
	isntTheLastDayOfRange := !w.Days[0].Time.Equal(rangeEnd)
	return stillInRange && isntTheLastDayOfRange
}

func (w *Week) PrevExists() bool {
	stillInRange := !w.Days[0].Time.Before(rangeStart)
	isntTheFirstDayOfRange := !w.Days[0].Time.Equal(rangeStart)
	return stillInRange && isntTheFirstDayOfRange
}

func (w *Week) Next() *Week {
	return fillWeekly(w.Weekday, w.Year, w.Days[0].Add(7))
}

func (w *Week) Prev() *Week {
	return fillWeekly(w.Weekday, w.Year, w.Days[0].Add(-7))
}

func (w *Week) QuartersBreadcrumb() header.ItemsGroup {
	group := header.ItemsGroup{}.Delim(" / ")

	for _, quarter := range w.Quarters {
		group.Items = append(group.Items, header.NewTextItem("Q"+strconv.Itoa(quarter.Number)))
	}

	return group
}

func (w *Week) MonthsBreadcrumb() header.ItemsGroup {
	group := header.ItemsGroup{}.Delim(" / ")

	for _, month := range w.Months {
		group.Items = append(group.Items, header.NewMonthItem(month.Month))
	}

	return group
}

func (w *Week) ref() string {
	prefix := ""
	wn := w.weekNumber()
	rm := w.rightMonth()
	ry := w.rightYear()

	if isClassicCalendarYear() && wn > 50 && rm == time.January && ry == rangeStartCalYear {
		prefix = "fw"
	}

	return prefix + "Week " + strconv.Itoa(wn)
}

func (w *Week) leftMonth() time.Month {
	for _, day := range w.Days {
		if day.Time.IsZero() {
			continue
		}

		return day.Time.Month()
	}

	return -1
}

func (w *Week) leftYear() int {
	for _, day := range w.Days {
		if day.Time.IsZero() {
			continue
		}

		return day.Time.Year()
	}

	return -1
}

func (w *Week) rightYear() int {
	for i := 6; i >= 0; i-- {
		if w.Days[i].Time.IsZero() {
			continue
		}

		return w.Days[i].Time.Year()
	}

	return -1
}

func (w *Week) HeadingMOS() string {
	var contents []string

	if w.PrevExists() {
		leftNavBox := tex.ResizeBoxW(`\myLenHeaderResizeBox`, `$\langle$`)
		contents = append(contents, tex.Hyperlink(w.Prev().ref(), leftNavBox))
	}

	contents = append(contents, tex.ResizeBoxW(`\myLenHeaderResizeBox`, w.Target()))

	if w.NextExists() {
		rightNavBox := tex.ResizeBoxW(`\myLenHeaderResizeBox`, `$\rangle$`)
		contents = append(contents, tex.Hyperlink(w.Next().ref(), rightNavBox))
	}

	return tex.Tabular("@{}"+strings.Repeat("l", len(contents)), strings.Join(contents, ` & `))
}

func (w *Week) Name() string {
	return "Week " + strconv.Itoa(w.weekNumber())
}

func (w *Week) Target() string {
	return tex.Hypertarget(w.ref(), w.Name())
}
