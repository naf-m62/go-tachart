package tachart

import (
	"fmt"
	"strings"

	"github.com/naf-m62/go-tachart/charts"
	"github.com/naf-m62/go-tachart/opts"
)

type bar struct {
	nm   string
	vals []float64
	ci   int
	dp   int
}

func NewBar(name string, vals []float64, color int) Indicator {
	return &bar{
		nm:   name,
		vals: vals,
		ci:   color,
		dp:   decimals(vals),
	}
}

func (b bar) name() string {
	return b.nm
}

func (b bar) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", b.dp), -1)
}

func (b bar) yAxisMin() string {
	return strings.Replace(minRoundFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", b.dp), -1)
}

func (b bar) yAxisMax() string {
	return strings.Replace(maxRoundFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", b.dp), -1)
}

func (b bar) getColor() int {
	return b.ci
}

func (b *bar) getTitleOpts(top, left int) []opts.Title {
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

func (b bar) genChart(_, _, _, _, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	barItems := []opts.BarData{}
	for _, v := range b.vals {
		style := &opts.ItemStyle{
			Color:   getColor(b.ci),
			Opacity: opacityHeavy,
		}
		barItems = append(barItems, opts.BarData{Value: v, ItemStyle: style})
	}
	return charts.NewBar().
		SetXAxis(xAxis).
		AddSeries(b.nm, barItems, charts.WithBarChartOpts(opts.BarChart{
			BarWidth:   "60%",
			XAxisIndex: gridIndex,
			YAxisIndex: gridIndex,
			ZLevel:     100,
		}))
}

func (b *bar) calcVals(values []float64) [][]float64 {
	panic("not implemented")
}

func (b bar) getDrawType() string {
	return "bar"
}
