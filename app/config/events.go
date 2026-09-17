package config

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Event kinds.
const (
	// EventHoliday is a day with no school: public holidays, "días no
	// lectivos" and school vacations. Day cells are shaded and the daily
	// page shows no classes or period boxes.
	EventHoliday = "holiday"
	// EventInfo is an informational date (start/end of the school year).
	// The day number is boxed in every calendar; otherwise a normal day.
	EventInfo = "info"
	// EventNote is like EventInfo but is only labelled on the monthly,
	// weekly and daily pages, not marked in the small annual/quarterly
	// calendars (e.g. exam periods, which would clutter them).
	EventNote = "note"
)

// Events is the school calendar: holidays and notable dates.
type Events []Event

// Event is one calendar entry, a single day (From only) or a range of
// days (From..To, both inclusive).
type Event struct {
	// From and To are dates as "YYYY-MM-DD". To defaults to From.
	From string
	To   string
	// Name is the full text shown on the daily page.
	Name string
	// Short is the label used where space is tight (monthly and weekly
	// pages). Defaults to Name.
	Short string
	// Kind is EventHoliday (default) or EventInfo.
	Kind string
}

// DayEvent is an event as seen from one particular day.
type DayEvent struct {
	Name    string
	Short   string
	Holiday bool
	// Note marks an EventNote entry (labelled, but not marked in the
	// small calendars).
	Note bool
}

const dateLayout = "2006-01-02"

// parseDate parses "YYYY-MM-DD" in the local time zone, matching how the
// cal package builds its days.
func parseDate(s string) (time.Time, error) {
	t, err := time.ParseInLocation(dateLayout, strings.TrimSpace(s), time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("want YYYY-MM-DD, got %q", s)
	}

	return t, nil
}

// DayKey is the map key used for a day: its date as "YYYY-MM-DD".
func DayKey(t time.Time) string {
	return t.Format(dateLayout)
}

// ByDay expands the events into a map from day key to the events on that
// day, in the order they appear in the config.
func (e Events) ByDay() (map[string][]DayEvent, error) {
	out := map[string][]DayEvent{}

	for _, ev := range e {
		from, err := parseDate(ev.From)
		if err != nil {
			return nil, fmt.Errorf("event %q from: %w", ev.Name, err)
		}

		to := from

		if ev.To != "" {
			if to, err = parseDate(ev.To); err != nil {
				return nil, fmt.Errorf("event %q to: %w", ev.Name, err)
			}
		}

		if to.Before(from) {
			return nil, fmt.Errorf("event %q: to %s is before from %s", ev.Name, ev.To, ev.From)
		}

		kind := ev.Kind
		if kind == "" {
			kind = EventHoliday
		}

		if kind != EventHoliday && kind != EventInfo && kind != EventNote {
			return nil, fmt.Errorf("event %q: unknown kind %q (want %s, %s or %s)", ev.Name, kind, EventHoliday, EventInfo, EventNote)
		}

		short := ev.Short
		if short == "" {
			short = ev.Name
		}

		for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
			key := DayKey(d)
			out[key] = append(out[key], DayEvent{Name: ev.Name, Short: short, Holiday: kind == EventHoliday, Note: kind == EventNote})
		}
	}

	return out, nil
}

// Sorted returns the events ordered by start date (for listings).
func (e Events) Sorted() Events {
	out := make(Events, len(e))
	copy(out, e)
	sort.SliceStable(out, func(i, j int) bool { return out[i].From < out[j].From })

	return out
}
