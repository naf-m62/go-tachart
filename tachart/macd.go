package tachart

import (
	"fmt"
	"strings"

	"github.com/iamjinlei/go-tart"

	"github.com/naf-m62/go-tachart/charts"
	"github.com/naf-m62/go-tachart/opts"
)

type macd struct {
	nm                  string
	fast                int64
	slow                int64
	signal              int64
	ci                  int
	macd, signalv, hist []float64
}

func NewMACD(fast, slow, signal int) Indicator {
	return &macd{
		nm:     fmt.Sprintf("MACD(%v,%v,%v)", fast, slow, signal),
		fast:   int64(fast),
		slow:   int64(slow),
		signal: int64(signal),
	}
}

func (c macd) name() string {
	return c.nm
}

func (c macd) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (c macd) yAxisMin() string {
	// looking for min value from signalv
	minVal := c.signalv[0]
	for _, v := range c.signalv {
		if v < minVal {
			minVal = v
		}
	}
	minVal *= 1.1
	return fmt.Sprintf(`function(value) { return %v }`, minVal)
}

func (c macd) yAxisMax() string {
	// looking for max value from signalv
	maxVal := c.signalv[0]
	for _, v := range c.signalv {
		if v > maxVal {
			maxVal = v
		}
	}
	maxVal *= 1.1
	return fmt.Sprintf(`function(value) { return %v }`, maxVal)
}

func (c macd) getNumColors() int {
	return 2
}

func (c *macd) calcVals(closes []float64) {
	c.macd, c.signalv, c.hist = tart.MacdArr(closes, c.fast, c.slow, c.signal)
	// умножаем на 100
	for i := range c.macd {
		c.macd[i] *= 100
	}
	for i := range c.signalv {
		c.signalv[i] *= 100
	}
	for i := range c.hist {
		c.hist[i] *= 100
	}
}

func (c *macd) getTitleOpts(top, left int, colorIndex int) []opts.Title {
	c.ci = colorIndex
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    getColor(c.ci),
				FontSize: chartLabelFontSize,
			},
			Title: c.nm + "-Diff",
			Left:  px(left),
			Top:   px(top),
		},
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    getColor(c.ci + 1),
				FontSize: chartLabelFontSize,
			},
			Title: c.nm + "-Sig",
			Left:  px(left),
			Top:   px(top + chartLabelFontHeight),
		},
	}
}

func (c macd) genChart(_, _, _, closes, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	macd, signal, hist := c.macd, c.signalv, c.hist

	// macd = macd[c.countExtraCandles()-1:]
	// hist = hist[c.countExtraCandles()-1:]
	// signal = signal[c.countExtraCandles()-1:]

	lineItems := []opts.LineData{}
	for _, v := range macd {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}
	macdLine := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(c.nm+"-Diff", lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   getColor(c.ci),
				Opacity: opacityMed,
			}),
		)

	lineItems = []opts.LineData{}
	for _, v := range signal {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}
	signalLine := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(c.nm+"-Sig", lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   getColor(c.ci + 1),
				Opacity: opacityMed,
			}),
		)

	barItems := []opts.BarData{}
	for _, v := range hist {
		style := &opts.ItemStyle{
			Color:   colorUpBar,
			Opacity: opacityHeavy,
		}
		if v < 0 {
			style = &opts.ItemStyle{
				Color:   colorDownBar,
				Opacity: opacityHeavy,
			}
		}
		barItems = append(barItems, opts.BarData{Value: v, ItemStyle: style})
	}
	histBar := charts.NewBar().
		SetXAxis(xAxis).
		AddSeries(c.nm+"-Hist", barItems, charts.WithBarChartOpts(opts.BarChart{
			BarWidth:   "60%",
			XAxisIndex: gridIndex,
			YAxisIndex: gridIndex,
			ZLevel:     100,
		}))

	macdLine.Overlap(signalLine, histBar)

	return macdLine
}
