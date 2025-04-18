package tachart

import (
	"fmt"

	"github.com/iamjinlei/go-tart"

	"github.com/naf-m62/go-tachart/charts"
	"github.com/naf-m62/go-tachart/opts"
)

type ma struct {
	nm string
	n  int64
	fn func([]float64, int64) []float64
	ci int
}

func NewSMA(n int) Indicator {
	return &ma{
		nm: fmt.Sprintf("SMA(%v)", n),
		n:  int64(n),
		fn: tart.SmaArr,
	}
}

func NewEMA(n int) Indicator {
	return &ma{
		nm: fmt.Sprintf("EMA(%v)", n),
		n:  int64(n),
		fn: tart.EmaArr,
	}
}

func (c ma) name() string {
	return c.nm
}

func (c ma) yAxisLabel() string {
	return ""
}

func (c ma) yAxisMin() string {
	return ""
}

func (c ma) yAxisMax() string {
	return ""
}

func (c ma) getNumColors() int {
	return 1
}

func (c *ma) getTitleOpts(top, left int, colorIndex int) []opts.Title {
	c.ci = colorIndex
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    getColor(c.ci),
				FontSize: chartLabelFontSize,
			},
			Title: c.nm,
			Left:  px(left),
			Top:   px(top),
		},
	}
}

func (c ma) genChart(_, _, _, closes, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	// check len
	if len(closes) < int(c.n) {
		return charts.NewLine()
	}
	ma := c.fn(closes, c.n)
	for i := 0; i < int(c.n); i++ {
		ma[i] = ma[c.n]
	}

	// ma = ma[c.countExtraCandles()-1:]

	items := []opts.LineData{}
	for i, v := range ma {
		if i < int(c.n) {
			items = append(items, opts.LineData{})
			continue
		}
		items = append(items, opts.LineData{Value: v})
	}

	return charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(c.nm, items,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   getColor(c.ci),
				Opacity: opacityMed,
			}))
}

// calcVals implements Indicator. Need for cal yMin and yMax
func (c *ma) calcVals(vals []float64) {
	panic("unimplemented")
}
