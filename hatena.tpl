---
Title: slack one day({{range .}}{{ .Postdate }}{{end}})
Category:
- slack
Date: {{range .}}{{ .Postdate }}{{end}}T06:00:00+09:00
---
## information
{{range .}}
### {{ .Posttime }}

{{ if ne .User ""}}{{ .User }} posted{{ end }}

{{ .Message }}

{{ .Text }}

Postes at {{ .PosttimeDetail }}
{{end}}