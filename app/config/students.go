package config

// StudentGroups is the class rosters: one entry per teaching group,
// keyed by the same "group" string used in Schedule.Classes (see
// cfg/teacher_schedule.yaml), so a class block on the daily schedule can
// be linked to its student list and, from there, to its attendance
// sheets (see cfg/teacher_students.yaml).
type StudentGroups []StudentGroup

// StudentGroup is one group's roster, in the order attendance is taken.
type StudentGroup struct {
	Group    string
	Students []string
}

// ByGroup indexes the rosters by group name for quick lookup when
// building the attendance pages.
func (g StudentGroups) ByGroup() map[string][]string {
	out := make(map[string][]string, len(g))

	for _, sg := range g {
		out[sg.Group] = sg.Students
	}

	return out
}
