{{- /*
  Daily schedule column with the teacher's class timetable.

  The whole column is one TikZ picture whose y unit is one hour
  (\myLenDailyHourHeight, pointing downwards), so a class from 08:30 to
  09:20 on a grid starting at 08:00 is simply a rectangle from y=0.5 to
  y=1.333. Hour labels sit in a left gutter (\myLenScheduleGutter) so the
  class blocks never cover them.
*/ -}}
{{- $bottom := .Cfg.Layout.Numbers.DailyBottomHour -}}
{{- $top := .Cfg.Layout.Numbers.DailyTopHour -}}
{{- $hours := .Day.Hours $bottom $top -}}
{{- $n := len $hours -}}
\myUnderline{Schedule\textcolor{white}{g}}\vskip-\myLenLineThicknessDefault
\noindent\begin{tikzpicture}[x=1mm, y={(0,-\myLenDailyHourHeight)}, inner sep=0pt, outer sep=0pt, line width=\myLenLineThicknessDefault]
  \useasboundingbox (0,0) rectangle (\myLenTriCol,{{$n}});
{{- range $i, $hour := $hours}}
  \draw[color=\myColorLightGray] (0,{{$i}}.5) -- (\myLenTriCol,{{$i}}.5);
  \draw[color=\myColorGray] (0,{{incr $i}}) -- (\myLenTriCol,{{incr $i}});
  \node[anchor=north west, inner sep=1pt] at (0,{{$i}}) {\small {{- $hour.FormatHour $.Cfg.AMPMTime -}} };
{{- end}}
{{- range $c := .Cfg.Schedule.ForWeekday .Day.Time.Weekday $bottom $top}}
  \filldraw[fill=\myColorClassFill, draw=\myColorGray] (\myLenScheduleGutter,{{printf "%.4f" $c.Top}}) rectangle (\myLenTriCol,{{printf "%.4f" $c.Bottom}});
  \node[anchor=north west, inner sep=2pt, align=left, text width=\dimexpr\myLenTriCol-\myLenScheduleGutter-4pt\relax] at (\myLenScheduleGutter,{{printf "%.4f" $c.Top}}) {\small\textbf{ {{- $c.Name -}} }\\[-1pt]\scriptsize\textcolor{\myColorGray}{ {{- $c.Start}}--{{$c.End -}} }};
{{- end}}
\end{tikzpicture}\par
