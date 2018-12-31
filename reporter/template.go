package reporter

import "html/template"

var reportTemplate = template.Must(template.ParseFiles("../reporter/template.html"))
