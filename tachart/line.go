package tachart

import (
	"fmt"
	"math"
	"strings"

	"github.com/naf-m62/go-tachart/charts"
	"github.com/naf-m62/go-tachart/opts"
)

type line struct {
	nms     []string
	valsArr [][]float64
	nc      int
	ci      int
	dp      int
}

func NewLine(name string, vals []float64, color int) Indicator {
	vals = zeroValsToNaN(vals)
	return &line{
		nms:     []string{name},
		valsArr: [][]float64{vals},
		nc:      1,
		ci:      color,
		dp:      decimals(vals),
	}
}

func NewLine2(n0 string, vals0 []float64, n1 string, vals1 []float64, color int) Indicator {
	vals0 = zeroValsToNaN(vals0)
	vals1 = zeroValsToNaN(vals1)
	return &line{
		nms:     []string{n0, n1},
		valsArr: [][]float64{vals0, vals1},
		nc:      2,
		ci:      color,
		dp:      decimals(vals0, vals1),
	}
}

func NewLine3(n0 string, vals0 []float64, n1 string, vals1 []float64, n2 string, vals2 []float64, color int) Indicator {
	vals0 = zeroValsToNaN(vals0)
	vals1 = zeroValsToNaN(vals1)
	vals2 = zeroValsToNaN(vals2)
	return &line{
		nms:     []string{n0, n1, n2},
		valsArr: [][]float64{vals0, vals1, vals2},
		nc:      3,
		ci:      color,
		dp:      decimals(vals0, vals1, vals2),
	}
}

func (b line) name() string {
	return strings.Join(b.nms, ", ")
}

func (b line) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", b.dp), -1)
}

func (b line) yAxisMin() string {
	return strings.Replace(minRoundFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", b.dp), -1)
}

func (b line) yAxisMax() string {
	return strings.Replace(maxRoundFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", b.dp), -1)
}

func (b line) getColor() int {
	return b.ci
}

func (b *line) getTitleOpts(top, left int) []opts.Title {
	var tls []opts.Title
	for i, nm := range b.nms {
		tls = append(tls, opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    getColor(b.ci + i),
				FontSize: chartLabelFontSize,
			},
			Title: nm,
			Left:  px(left),
			Top:   px(top + i*chartLabelFontHeight),
		})
	}
	return tls
}

func (b line) genChart(_, _, _, _, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	lineItems := []opts.LineData{}
	for _, v := range b.valsArr[0] {
		if v == 0 || math.IsNaN(v) {
			lineItems = append(lineItems, opts.LineData{})
			continue
		}
		lineItems = append(lineItems, opts.LineData{Value: v})
	}

	c := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(b.nms[0], lineItems,
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

	for i := 1; i < len(b.nms); i++ {
		lineItems := []opts.LineData{}
		for _, v := range b.valsArr[i] {
			if v == 0 || math.IsNaN(v) {
				lineItems = append(lineItems, opts.LineData{})
				continue
			}
			lineItems = append(lineItems, opts.LineData{Value: v})
		}

		line := charts.NewLine().
			SetXAxis(xAxis).
			AddSeries(b.nms[i], lineItems,
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
		c.Overlap(line)
	}

	return c
}

// calcVals implements Indicator.
func (b *line) calcVals(vals []float64) [][]float64 {
	return b.valsArr
}

func (b *line) getDrawType() string {
	return "line"
}

func zeroValsToNaN(vals []float64) []float64 {
	for i, v := range vals {
		if v == 0 {
			vals[i] = math.NaN()
		}
	}
	return vals
}
