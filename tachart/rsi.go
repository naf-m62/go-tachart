package tachart

import (
	"fmt"
	"strings"

	"github.com/iamjinlei/go-tart"

	"github.com/naf-m62/go-tachart/charts"
	"github.com/naf-m62/go-tachart/opts"
)

type rsi struct {
	nm         string
	n          int64
	oversold   float64
	overbought float64
	ci         int
	vals       []float64
}

func NewRSI(n int, oversold, overbought float64) Indicator {
	return &rsi{
		nm:         fmt.Sprintf("RSI(%v)", n),
		n:          int64(n),
		oversold:   oversold,
		overbought: overbought,
	}
}

func (r rsi) name() string {
	return r.nm
}

func (r rsi) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (r rsi) yAxisMin() string {
	// get min value from vals[r.n:]
	if len(r.vals) < int(r.n) {
		return fmt.Sprintf(`function(value) { return %v }`, 0)
	}
	minVal := 100.0
	for _, v := range r.vals[r.n:] {
		if v < minVal {
			minVal = v
		}
	}
	minVal = minVal - 10
	if minVal < 0 {
		minVal = 0
	}
	return fmt.Sprintf(`function(value) { return %v }`, minVal)
}

func (r rsi) yAxisMax() string {
	// get max value from vals[r.n:]
	if len(r.vals) < int(r.n) {
		return fmt.Sprintf(`function(value) { return %v }`, 100)
	}
	maxVal := 0.0
	for _, v := range r.vals[r.n:] {
		if v > maxVal {
			maxVal = v
		}
	}
	maxVal = maxVal + 10
	if maxVal > 100 {
		maxVal = 100
	}
	return fmt.Sprintf(`function(value) { return %v }`, maxVal)
}

func (r rsi) getNumColors() int {
	return 1
}

func (r *rsi) getTitleOpts(top, left int, colorIndex int) []opts.Title {
	r.ci = colorIndex
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    getColor(r.ci),
				FontSize: chartLabelFontSize,
			},
			Title: r.nm,
			Left:  px(left),
			Top:   px(top),
		},
	}
}

func (r *rsi) calcVals(closes []float64) {
	r.vals = tart.RsiArr(closes, r.n)
}

func (r rsi) genChart(_, _, _, closes, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	vals := r.vals

	// vals = vals[r.countExtraCandles()-1:]

	lineItems := []opts.LineData{}
	for _, v := range vals {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}

	return charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(r.nm, lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   getColor(r.ci),
				Opacity: opacityMed,
			}),
			charts.WithMarkLineNameYAxisItemOpts(
				opts.MarkLineNameYAxisItem{
					Name:  "oversold",
					YAxis: r.oversold,
				},
				opts.MarkLineNameYAxisItem{
					Name:  "overbought",
					YAxis: r.overbought,
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
