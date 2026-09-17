package config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Schedule is a weekly class timetable. When Enabled, the daily pages
// draw the hour grid as a fixed-height column with one shaded block per
// class that falls on that weekday (see tpls/schedule_classes.tpl).
type Schedule struct {
	Enabled bool
	Classes []Class
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

// ClassBlock is a Class placed on the daily hour grid: Top and Bottom are
// the class start/end expressed in hours since the grid's bottom hour
// (e.g. with a grid starting at 08:00, a 08:30-09:20 class has Top 0.5
// and Bottom 1.333). They are consumed by the TikZ template, which uses
// the hour height as its y unit.
type ClassBlock struct {
	Class
	Top    float64
	Bottom float64
}

// ForWeekday returns the classes on the given weekday, sorted by start
// time, positioned relative to bottomHour. Classes that fall entirely
// outside [bottomHour, topHour+1) are dropped; partial overlaps are
// clipped so the block never leaves the grid.
func (s Schedule) ForWeekday(wd time.Weekday, bottomHour, topHour int) ([]ClassBlock, error) {
	gridTop := float64(topHour + 1 - bottomHour)
	blocks := make([]ClassBlock, 0, len(s.Classes))

	for _, c := range s.Classes {
		if c.Day != wd {
			continue
		}

		start, err := parseHHMM(c.Start)
		if err != nil {
			return nil, fmt.Errorf("class %q start: %w", c.Name, err)
		}

		end, err := parseHHMM(c.End)
		if err != nil {
			return nil, fmt.Errorf("class %q end: %w", c.Name, err)
		}

		if end <= start {
			return nil, fmt.Errorf("class %q: end %s is not after start %s", c.Name, c.End, c.Start)
		}

		top := start - float64(bottomHour)
		bottom := end - float64(bottomHour)

		if bottom <= 0 || top >= gridTop {
			continue
		}

		if top < 0 {
			top = 0
		}

		if bottom > gridTop {
			bottom = gridTop
		}

		blocks = append(blocks, ClassBlock{Class: c, Top: top, Bottom: bottom})
	}

	sort.Slice(blocks, func(i, j int) bool { return blocks[i].Top < blocks[j].Top })

	return blocks, nil
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
