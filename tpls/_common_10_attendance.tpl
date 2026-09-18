{{- /*
  Attendance/roster page: one per (teaching group, calendar month) with
  at least one class session (see app/compose/attendance.go, which
  builds the whole table in Go and hands it over ready-made, since its
  column count varies with the month). The table is shorter than the
  page, so \vfill on both sides centers it in the space left below the
  header instead of leaving it stuck to the top with a blank gap below.
*/ -}}
\vfill
{{- .Body.Table }}
\vfill
