package config

type Lengths struct {
	TabColSep               string
	LineThicknessDefault    string
	LineThicknessThick      string
	LineHeightButLine       string
	TwoColSep               string
	TriColSep               string
	FiveColSep              string
	MonthlyCellHeight       string
	NotesIndexCellHeight    string
	HeaderResizeBox         string
	HeaderSideCellHeight    string
	HeaderSideQuartersWidth string
	HeaderSideMonthsWidth   string
	QuarterlySpring         string
	MonthlySpring           string
	DailySpring             string

	// DailyHourHeight is the vertical space given to one hour on the daily
	// schedule column when Schedule.Enabled is set (e.g. "2.3cm").
	DailyHourHeight string
	// ScheduleGutter is the space reserved on the left of the daily
	// schedule column for the hour labels; class blocks start after it.
	ScheduleGutter string
}
