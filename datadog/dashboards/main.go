package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type Dashboard struct {
	Description       string
	Filename          string
	LayoutType        string
	Model             string
	NotifyList        []string
	ReflowType        string
	TemplateVariables []string
	Title             string
	Widgets           []Widget
}

type Widget struct {
	ID         int
	Definition WidgetDef
	Layout     WidgetLayout
}

type WidgetDef struct {
	Title         string
	TitleSize     string
	TitleAlign    string
	ShowLegend    bool
	LegendLayout  string
	LegendColumns []string
	Type          string
	Request       WidgetRequest
}

type WidgetRequest struct {
	ResponseFormat string
	Queries        []Query
	Formulas       []Formula
	DisplayType    string
}

type Query struct {
	Name       string
	DataSource string
	Query      string
}

type Formula struct {
	Formula string
	Style   Style
}

type Style struct {
	Palette      string
	PaletteIndex int
}

type WidgetLayout struct {
	X      int
	Y      int
	Width  int
	Height int
}

func readMetrics(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var groups [][]string
	var currentGroup []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
				currentGroup = []string{}
			}
		} else {
			currentGroup = append(currentGroup, fmt.Sprintf("avg:%s{*}", line))
		}
	}

	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func main() {
	groups, err := readMetrics("metrics.txt")
	if err != nil {
		fmt.Println("Error reading metrics file:", err)
		os.Exit(1)
	}

	var widgets []Widget
	id := 150000000000000
	x, y := 0, 0
	paletteIndex := 0

	for _, group := range groups {
		queries := []Query{}
		formulas := []Formula{}
		for i, metric := range group {
			queryName := fmt.Sprintf("query%d", i)
			queries = append(queries, Query{
				Name:       queryName,
				DataSource: "metrics",
				Query:      metric,
			})
			formulas = append(formulas, Formula{
				Formula: queryName,
				Style: Style{
					Palette:      "dd20",
					PaletteIndex: paletteIndex % 20,
				},
			})
			paletteIndex++
		}

		widget := Widget{
			ID: id,
			Definition: WidgetDef{
				Title:         "",
				TitleSize:     "16",
				TitleAlign:    "left",
				ShowLegend:    false,
				LegendLayout:  "auto",
				LegendColumns: []string{"avg", "min", "max", "value", "sum"},
				Type:          "timeseries",
				Request: WidgetRequest{
					ResponseFormat: "timeseries",
					Queries:        queries,
					Formulas:       formulas,
					DisplayType:    "area",
				},
			},
			Layout: WidgetLayout{
				X:      x,
				Y:      y,
				Width:  4,
				Height: 2,
			},
		}

		widgets = append(widgets, widget)
		id++
		x += 4
		if x >= 12 {
			x = 0
			y += 2
		}
	}

	models := []string{"BGW210-700", "BGW320-500", "BGW320-505"}

	for _, model := range models {
		metricModel := strings.ToLower(model)
		metricModel = strings.Replace(metricModel, "-", "", 1)
		filename := model + ".json"
		title := model + " dashboard"

		dashboard := Dashboard{
			Description:       "",
			Filename:          filename,
			LayoutType:        "ordered",
			Model:             metricModel,
			NotifyList:        []string{},
			ReflowType:        "fixed",
			TemplateVariables: []string{},
			Title:             title,
			Widgets:           widgets,
		}

		tmpl := template.Must(template.New("dashboard").Parse(`{
  "title": "{{.Title}}",
  "description": "{{.Description}}",
  "widgets": [
    {{- range $i, $widget := .Widgets -}}
    {{if $i}},{{end}}
    {
      "id": {{$widget.ID}},
      "definition": {
        "title": "{{$widget.Definition.Title}}",
        "title_size": "{{$widget.Definition.TitleSize}}",
        "title_align": "{{$widget.Definition.TitleAlign}}",
        "show_legend": {{$widget.Definition.ShowLegend}},
        "legend_layout": "{{$widget.Definition.LegendLayout}}",
        "legend_columns": [
          {{- range $j, $col := $widget.Definition.LegendColumns}}
          {{- if $j}},{{end}}
          "{{$col}}"
          {{- end}}
        ],
        "type": "{{$widget.Definition.Type}}",
        "requests": [
          {
            "response_format": "{{$widget.Definition.Request.ResponseFormat}}",
            "queries": [
              {{- range $j, $query := $widget.Definition.Request.Queries -}}
              {{if $j}},{{end}}
              {
                "name": "{{$query.Name}}",
                "data_source": "{{$query.DataSource}}",
                "query": "{{$query.Query}}"
              }
              {{- end}}
            ],
            "formulas": [
              {{- range $j, $formula := $widget.Definition.Request.Formulas -}}
              {{if $j}},{{end}}
              {
                "formula": "{{$formula.Formula}}",
                "style": {
                  "palette": "{{$formula.Style.Palette}}",
                  "palette_index": {{$formula.Style.PaletteIndex}}
                }
              }
              {{- end}}
            ],
            "display_type": "{{$widget.Definition.Request.DisplayType}}"
          }
        ]
      },
      "layout": {
        "x": {{$widget.Layout.X}},
        "y": {{$widget.Layout.Y}},
        "width": {{$widget.Layout.Width}},
        "height": {{$widget.Layout.Height}}
      }
    }
    {{- end}}
  ],
  "template_variables": [],
  "layout_type": "{{.LayoutType}}",
  "notify_list": [],
  "reflow_type": "{{.ReflowType}}"
}`))

		re := regexp.MustCompile(`bgw\w{6}`)

		for wIndex, widget := range dashboard.Widgets {
			for qIndex, query := range widget.Definition.Request.Queries {
				query.Query = re.ReplaceAllString(query.Query, metricModel)
				// Write back the modified query into the original Queries slice
				dashboard.Widgets[wIndex].Definition.Request.Queries[qIndex] = query
			}

		}

		outputDashboard := struct {
			Title             string
			Description       string
			Widgets           []Widget
			TemplateVariables []string
			LayoutType        string
			NotifyList        []string
			ReflowType        string
		}{
			Title:             dashboard.Title,
			Description:       dashboard.Description,
			Widgets:           dashboard.Widgets,
			TemplateVariables: dashboard.TemplateVariables,
			LayoutType:        dashboard.LayoutType,
			NotifyList:        dashboard.NotifyList,
			ReflowType:        dashboard.ReflowType,
		}

		var output bytes.Buffer
		if err := tmpl.Execute(&output, outputDashboard); err != nil {
			fmt.Println("Error executing template:", err)
			os.Exit(1)
		}

		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			os.Exit(1)
		}
		defer file.Close()

		file.Write(output.Bytes())
		fmt.Printf("Dashboard JSON written to %s\n", filename)

	}
}
