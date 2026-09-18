{{- /*
  Daily schedule column with the teacher's timetable.

  The whole column is one TikZ picture whose y unit is one hour
  (\myLenDailyHourHeight, pointing downwards), so a slot from 08:30 to
  09:20 on a grid starting at 08:00 is simply a rectangle from y=0.5 to
  y=1.333. Hour labels sit in a left gutter (\myLenScheduleGutter) so the
  boxes never cover them.

  Three kinds of box (see config.Schedule.ForWeekday):
    class  shaded, group name in bold + time
    free   outlined only (a period without a lesson), time in gray
    break  shaded darker, name centred (e.g. RECREO)
*/ -}}
{{- $bottom := .Cfg.Layout.Numbers.DailyBottomHour -}}
{{- $top := .Cfg.Layout.Numbers.DailyTopHour -}}
{{- $hours := .Day.Hours $bottom $top -}}
{{- $n := len $hours -}}
\myUnderline{Schedule\textcolor{white}{g}}\vskip-\myLenLineThicknessDefault
{{- if .Day.HasEvents}}
{{- if .Day.IsHoliday}}
\noindent{\setlength{\fboxsep}{2pt}\colorbox{\myColorHolidayFill}{\parbox{\dimexpr\myLenTriCol-2\fboxsep\relax}{\centering\small\textbf{ {{- .Day.EventsLong -}} }}}}\par\vskip1pt
{{- else}}
\noindent\parbox{\myLenTriCol}{\centering\small\textbf{ {{- .Day.EventsLong -}} }}\par\vskip1pt
{{- end}}
{{- end}}
\noindent\begin{tikzpicture}[x=1mm, y={(0,-\myLenDailyHourHeight)}, inner sep=0pt, outer sep=0pt, line width=\myLenLineThicknessDefault]
  \useasboundingbox (0,0) rectangle (\myLenTriCol,{{$n}});
{{- range $i, $hour := $hours}}
  \draw[color=\myColorLightGray] (0,{{$i}}.5) -- (\myLenTriCol,{{$i}}.5);
  \draw[color=\myColorGray] (0,{{incr $i}}) -- (\myLenTriCol,{{incr $i}});
  \node[anchor=north west, inner sep=1pt] at (0,{{$i}}) {\small {{- $hour.FormatHour $.Cfg.AMPMTime -}} };
{{- end}}
{{- range $b := .Cfg.Schedule.ForDate .Day.Time $bottom $top .Day.IsHoliday}}
{{- if eq $b.Kind "class"}}
{{- $ref := $b.AttendanceRef $.Day.Time}}
  \filldraw[fill=\myColorClassFill, draw=\myColorGray] (\myLenScheduleGutter,{{printf "%.4f" $b.Top}}) rectangle (\myLenTriCol,{{printf "%.4f" $b.Bottom}});
  \node[anchor=north west, inner sep=2pt, align=left, text width=\dimexpr\myLenTriCol-\myLenScheduleGutter-4pt\relax] at (\myLenScheduleGutter,{{printf "%.4f" $b.Top}}) {\small\textbf{ {{- $b.Name -}} }\\[-1pt]\scriptsize\textcolor{\myColorGray}{ {{- $b.Start}}--{{$b.End -}} {{if $ref}}\ \ \hyperlink{ {{- $ref -}} }{Attendance}{{end -}} }};
{{- else if eq $b.Kind "break"}}
  \filldraw[fill=\myColorBreakFill, draw=\myColorGray] (\myLenScheduleGutter,{{printf "%.4f" $b.Top}}) rectangle (\myLenTriCol,{{printf "%.4f" $b.Bottom}});
  \node[anchor=center, inner sep=0pt] at ({0.5*(\myLenScheduleGutter+\myLenTriCol)},{{printf "%.4f" $b.Mid}}) {\scriptsize\textcolor{\myColorGray}{\textbf{ {{- $b.Name -}} }\enspace {{$b.Start}}--{{$b.End}}}};
{{- else}}
  \filldraw[fill=white, draw=\myColorGray] (\myLenScheduleGutter,{{printf "%.4f" $b.Top}}) rectangle (\myLenTriCol,{{printf "%.4f" $b.Bottom}});
  \node[anchor=north west, inner sep=2pt] at (\myLenScheduleGutter,{{printf "%.4f" $b.Top}}) {\scriptsize\textcolor{\myColorGray}{ {{- if $b.Name}}\textbf{ {{- $b.Name -}} }\enspace{{end -}} {{$b.Start}}--{{$b.End -}} }};
{{- end}}
{{- end}}
\end{tikzpicture}\par
