package cal

import (
	"math"
	"strconv"
	"time"

	"github.com/kudrykv/latex-yearly-planner/app/components/header"
	"github.com/kudrykv/latex-yearly-planner/app/tex"
)

type Years []*Year
type Year struct {
	Number     int // starting calendar year of the range (e.g. 2025)
	StartMonth time.Month
	NumMonths  int

	Quarters Quarters
	Weeks    Weeks
}

// NewYear builds a Year covering numMonths consecutive months starting at
// (year, startMonth). Passing startMonth=January and numMonths=12
// reproduces the original plain-calendar-year behavior; any other start
// month/count builds a "school year" style range, which may span two
// calendar years (e.g. September 2025 - June 2026).
func NewYear(wd time.Weekday, year int, startMonth time.Month, numMonths int) *Year {
	if startMonth == 0 {
		startMonth = time.January
	}

	if numMonths <= 0 {
		numMonths = 12
	}

	out := &Year{Number: year, StartMonth: startMonth, NumMonths: numMonths}

	setRange(year, startMonth, numMonths)

	out.Weeks = NewWeeksForYear(wd, out)

	numQuarters := int(math.Ceil(float64(numMonths) / 3.))
	for q := 1; q <= numQuarters; q++ {
		out.Quarters = append(out.Quarters, NewQuarter(wd, out, q))
	}

	return out
}

func (y Year) Breadcrumb() string {
	return header.Items{
		header.NewIntItem(y.Number).Ref(),
		header.NewItemsGroup(
			header.NewTextItem("Q1"),
			header.NewTextItem("Q2"),
			header.NewTextItem("Q3"),
			header.NewTextItem("Q4"),
		),
	}.Table(true)
}

func (y Year) SideQuarters(sel ...int) []header.CellItem {
	out := make([]header.CellItem, 0, len(y.Quarters))

	for i := len(y.Quarters) - 1; i >= 0; i-- {
		mark := false
		for _, oneof := range sel {
			if oneof == y.Quarters[i].Number {
				mark = true

				break
			}
		}

		out = append(out, header.NewCellItem(y.Quarters[i].Name()).Selected(mark))
	}

	return out
}

func (y Year) SideMonths(sel ...time.Month) []header.CellItem {
	out := make([]header.CellItem, 0, 12)

	for i := len(y.Quarters) - 1; i >= 0; i-- {
		for j := len(y.Quarters[i].Months) - 1; j >= 0; j-- {
			mon := y.Quarters[i].Months[j]
			mark := false

			for _, month := range sel {
				if month == mon.Month {
					mark = true

					break
				}
			}

			cell := header.NewCellItem(mon.ShortName()).Refer(mon.Month.String()).Selected(mark)
			out = append(out, cell)
		}
	}

	return out
}

// Label is the visible name of the range: "2026" for a calendar year,
// "2026-2027" when the range spans two calendar years.
func (y Year) Label() string {
	endYear, _ := monthAt(y.NumMonths - 1)
	if endYear == y.Number {
		return strconv.Itoa(y.Number)
	}

	return strconv.Itoa(y.Number) + "-" + strconv.Itoa(endYear)
}

func (y Year) HeadingMOS() string {
	return tex.ResizeBoxW(`\myLenHeaderResizeBox`, tex.Hypertarget("Calendar", y.Label()))
}
