package tachart

import (
	"fmt"
	"strings"

	"github.com/naf-m62/go-tachart/charts"
	"github.com/naf-m62/go-tachart/opts"
)

type boundedLine struct {
	nm          string
	vals        []float64
	min         float64
	max         float64
	lowerMarker float64
	upperMarker float64
	ci          int
}

func NewBoundedLine(name string, vals []float64, min, max, lowerMarker, upperMarker float64, color int) Indicator {
	return &boundedLine{
		nm:          name,
		vals:        vals,
		min:         min,
		max:         max,
		lowerMarker: lowerMarker,
		upperMarker: upperMarker,
		ci:          color,
	}
}

func (b boundedLine) name() string {
	return b.nm
}

func (b boundedLine) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (b boundedLine) yAxisMin() string {
	return fmt.Sprintf("function(value) { return %v }", b.min)
}

func (b boundedLine) yAxisMax() string {
	return fmt.Sprintf("function(value) { return %v }", b.max)
}

func (b boundedLine) getColor() int {
	return b.ci
}

func (b *boundedLine) getTitleOpts(top, left int) []opts.Title {
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    getColor(b.ci),
				FontSize: chartLabelFontSize,
			},
			Title: b.nm,
			Left:  px(left),
			Top:   px(top),
		},
	}
}

func (b boundedLine) genChart(_, _, _, _, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	lineItems := []opts.LineData{}
	for _, v := range b.vals {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}

	return charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(b.nm, lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   getColor(b.ci),
				Opacity: opacityMed,
			}),
			charts.WithMarkLineNameYAxisItemOpts(
				opts.MarkLineNameYAxisItem{
					Name:  "lower_bound ",
					YAxis: b.lowerMarker,
				},
				opts.MarkLineNameYAxisItem{
					Name:  "upper_bound",
					YAxis: b.upperMarker,
				},
			),
			charts.WithMarkLineStyleOpts(
				opts.MarkLineStyle{
					Symbol: []string{"none", "none"},
					LineStyle: &opts.LineStyle{
						Color:   colorDownBar,
						Opacity: opacityMed,
					},
				},
			),
		)
}

func (b *boundedLine) calcVals(values []float64) [][]float64 {
	return [][]float64{b.vals}
}

func (b boundedLine) getDrawType() string {
	return "boundedline"
}
