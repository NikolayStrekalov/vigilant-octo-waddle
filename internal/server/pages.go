package server

import (
	"bytes"
	"errors"
	"html/template"
	"log"
)

type templateArgs struct {
	Gauge   []GaugeListItem
	Counter []CounterListItem
}

var indexTemplate = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Metrics list</title>
  </head>
  <body>
	<ul>{{ range .Gauge}}
	<li>gauge {{ .Name }} {{ .Value }}</li>{{ end }}{{ range .Counter}}
	<li>counter {{ .Name }} {{ .Value }}</li>{{ end }}
	</ul>
  </body>
</html>
`

func renderIndexPage(counters []CounterListItem, gauges []GaugeListItem) (*bytes.Buffer, error) {
	indexTemplate := template.Must(template.New("metrics").Parse(indexTemplate))
	buf := new(bytes.Buffer)
	if err := indexTemplate.Execute(buf, templateArgs{Gauge: gauges, Counter: counters}); err != nil {
		log.Println(err)
		return nil, errors.Unwrap(err)
	}
	return buf, nil
}
