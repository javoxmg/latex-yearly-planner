package config

import (
	"sort"
	"strings"
	"time"
)

// SessionDates returns, in chronological order, every calendar date on
// which the given group actually has class: a Class entry exists for
// that weekday, the date falls within that class's active range
// (Schedule.From/To or the class's own From/Until, see ForDate), and the
// date is not a school holiday. isHoliday may be nil (no holidays
// excluded).
//
// This is the base data for the attendance sheets: each returned date
// becomes one column on the roster page for that group.
func (s Schedule) SessionDates(group string, from, to time.Time, isHoliday func(time.Time) bool) ([]time.Time, error) {
	var out []time.Time

	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		if isHoliday != nil && isHoliday(d) {
			continue
		}

		wd := d.Weekday()
		met := false

		for _, c := range s.Classes {
			if c.Group != group || c.Day != wd {
				continue
			}

			ok, err := s.classActive(c, d)
			if err != nil {
				return nil, err
			}

			if ok {
				met = true

				break
			}
		}

		if met {
			out = append(out, d)
		}
	}

	return out, nil
}

// Groups returns the distinct, non-empty Class.Group values, in the
// order they first appear in Classes (the order the timetable was
// written in cfg/teacher_schedule.yaml).
func (s Schedule) Groups() []string {
	seen := map[string]bool{}

	var out []string

	for _, c := range s.Classes {
		if c.Group == "" || seen[c.Group] {
			continue
		}

		seen[c.Group] = true

		out = append(out, c.Group)
	}

	return out
}

// ShortName returns the short Class.Name first associated with group
// (e.g. "Mat II · 2ºBCB" for "Matemáticas II - 2BCB"), for use as a
// compact page heading. Falls back to the group string itself.
func (s Schedule) ShortName(group string) string {
	for _, c := range s.Classes {
		if c.Group == group {
			return c.Name
		}
	}

	return group
}

// AttendanceRef is the anchor to link to from a class block on the daily
// schedule: the attendance page for this block's group and the given
// day's month. Empty when the block has no group (a free/break box) or
// isn't a class block at all, so the template can skip the "Attendance"
// link in that case.
func (b ScheduleBlock) AttendanceRef(t time.Time) string {
	if b.Kind != BlockClass || b.Group == "" {
		return ""
	}

	return AttendanceAnchor(b.Group, t)
}

// AttendanceAnchor is the hyperref anchor name for the attendance page
// of the given group and month (any date within that month). It is
// built from ASCII only, so accents or punctuation in the group name
// never risk producing an invalid PDF name: the same group and calendar
// month always yield the same anchor, which is all that is needed to
// link a class block on the daily schedule to its roster page.
func AttendanceAnchor(group string, t time.Time) string {
	return "attendance-" + slugify(group) + "-" + t.Format("2006-01")
}

func slugify(s string) string {
	var b strings.Builder

	dash := true // start "true" so a leading run of separators is dropped

	for _, r := range s {
		r = unaccent(r)

		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			dash = false
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r - 'A' + 'a')
			dash = false
		default:
			if !dash {
				b.WriteByte('-')
				dash = true
			}
		}
	}

	return strings.Trim(b.String(), "-")
}

func unaccent(r rune) rune {
	switch r {
	case 'á', 'à', 'ä', 'â':
		return 'a'
	case 'é', 'è', 'ë', 'ê':
		return 'e'
	case 'í', 'ì', 'ï', 'î':
		return 'i'
	case 'ó', 'ò', 'ö', 'ô', 'º':
		return 'o'
	case 'ú', 'ù', 'ü', 'û':
		return 'u'
	case 'ñ':
		return 'n'
	case 'Á':
		return 'A'
	case 'É':
		return 'E'
	case 'Í':
		return 'I'
	case 'Ó':
		return 'O'
	case 'Ú':
		return 'U'
	case 'Ñ':
		return 'N'
	default:
		return r
	}
}

// SortedMonthKeys sorts "YYYY-MM" keys chronologically (plain string
// sort already works since the format is zero-padded, but this documents
// the intent at the call site).
func SortedMonthKeys(keys []string) []string {
	out := make([]string, len(keys))
	copy(out, keys)
	sort.Strings(out)

	return out
}
