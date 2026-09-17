package cal

import (
	"strings"
	"time"

	"github.com/kudrykv/latex-yearly-planner/app/tex"
)

// DayEvent is a calendar entry (holiday, vacation, notable date) attached
// to one day. The map is filled once per run by SetDayEvents, from the
// events in the config, and read by the Day helpers when rendering.
type DayEvent struct {
	Name    string
	Short   string
	Holiday bool
	Note    bool // labelled on monthly/weekly/daily pages only
}

var dayEvents map[string][]DayEvent

// SetDayEvents installs the day -> events map (keys are "YYYY-MM-DD").
func SetDayEvents(events map[string][]DayEvent) {
	dayEvents = events
}

// Events returns the events on this day, in config order.
func (d Day) Events() []DayEvent {
	if d.Time.IsZero() {
		return nil
	}

	return dayEvents[d.Time.Format("2006-01-02")]
}

// IsHoliday reports whether any event on this day is a holiday.
func (d Day) IsHoliday() bool {
	for _, ev := range d.Events() {
		if ev.Holiday {
			return true
		}
	}

	return false
}

// HasEvents reports whether the day has any event at all.
func (d Day) HasEvents() bool {
	return len(d.Events()) > 0
}

// EventsShort returns the short label to print where space is tight
// (monthly cells, weekly day headers): one label only, the day's most
// relevant event — a holiday first, then an info date, then a note. The
// daily page prints all of them (EventsLong).
func (d Day) EventsShort() string {
	evs := d.Events()
	if len(evs) == 0 {
		return ""
	}

	best := evs[0]

	for _, ev := range evs[1:] {
		if rank(ev) > rank(best) {
			best = ev
		}
	}

	return best.Short
}

func rank(ev DayEvent) int {
	switch {
	case ev.Holiday:
		return 2
	case !ev.Note:
		return 1
	default:
		return 0
	}
}

// EventsLong joins the full names of the day's events, for the daily page.
func (d Day) EventsLong() string {
	parts := make([]string, 0, 2)

	for _, ev := range d.Events() {
		parts = append(parts, ev.Name)
	}

	return strings.Join(parts, " / ")
}

// IsMarked reports whether the day gets a mark in the small calendars:
// it has a holiday or an info event (notes are not marked there).
func (d Day) IsMarked() bool {
	for _, ev := range d.Events() {
		if ev.Holiday || !ev.Note {
			return true
		}
	}

	return false
}

// markCell decorates a calendar cell for the day's events: holidays get
// a shaded background, informational dates a box around the number.
// Returns the cell content unchanged when the day has nothing to mark.
func (d Day) markCell(content string) string {
	if d.IsHoliday() {
		return tex.CellColor(`\myColorHolidayFill`, content)
	}

	if d.IsMarked() {
		return `\fbox{` + content + `}`
	}

	return content
}

// weekdayIsWeekend is a small helper for templates that want to know
// whether a day is Saturday or Sunday.
func (d Day) IsWeekend() bool {
	return d.Time.Weekday() == time.Saturday || d.Time.Weekday() == time.Sunday
}
