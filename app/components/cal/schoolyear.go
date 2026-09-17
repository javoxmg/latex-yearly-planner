package cal

import (
	"strconv"
	"time"
)

// The functions and package-level state in this file generalize the
// original "calendar year" assumption (always January through December)
// to an arbitrary contiguous range of months: a calendar year (start
// month January, 12 months) or an academic/school year (e.g. start
// month September, 10 months), possibly spanning two calendar years.
//
// NewYear sets these once per run (this program generates one range per
// invocation), and Day/Week/Month helpers read them so they don't each
// need to carry a pointer back to the *Year.
var (
	rangeStartCalYear int
	rangeStartMonth   time.Month
	rangeNumMonths    int
	rangeStart        time.Time // first day of the range
	rangeEnd          time.Time // last day of the range (inclusive)
)

// setRange records the boundaries of the range currently being generated.
func setRange(calYear int, startMonth time.Month, numMonths int) {
	rangeStartCalYear = calYear
	rangeStartMonth = startMonth
	rangeNumMonths = numMonths

	rangeStart = time.Date(calYear, startMonth, 1, 0, 0, 0, 0, time.Local)

	endYear, endMonth := monthAt(numMonths - 1)
	rangeEnd = time.Date(endYear, endMonth, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, -1)
}

// monthIndex returns the 0-based index of (year, month) relative to the
// start of the range, where index 0 is (rangeStartCalYear, rangeStartMonth).
func monthIndex(year int, month time.Month) int {
	return (year-rangeStartCalYear)*12 + int(month) - int(rangeStartMonth)
}

// monthAt is the inverse of monthIndex: given a 0-based index into the
// range, it returns the real calendar (year, month).
func monthAt(idx int) (int, time.Month) {
	total := int(rangeStartMonth) - 1 + idx
	year := rangeStartCalYear + total/12
	month := time.Month(total%12 + 1)

	return year, month
}

// quarterNumber returns the 1-based "quarter" that (year, month) falls
// into, grouping the range's months into consecutive groups of up to 3,
// starting from the range's first month.
func quarterNumber(year int, month time.Month) int {
	return monthIndex(year, month)/3 + 1
}

// weekRefString builds the unique reference string used to link to a
// week, from its ISO week number and ISO year. A plain "Week N" is not
// safe to use as a unique LaTeX anchor name here: a school-year range
// (say, September to August) covers close to a full year and can end on
// the very same ISO week number it started on, in a different ISO year.
// Including the ISO year in the ref keeps every week's anchor unique
// regardless of how the range is configured, while the visible page
// title (Week.Name/day breadcrumb text) stays just "Week N".
func weekRefString(isoWeek, isoYear int) string {
	return "Week " + strconv.Itoa(isoWeek) + "@" + strconv.Itoa(isoYear)
}
