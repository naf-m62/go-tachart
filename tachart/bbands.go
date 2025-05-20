package tachart

import (
	"fmt"

	"github.com/iamjinlei/go-tart"

	"github.com/naf-m62/go-tachart/charts"
	"github.com/naf-m62/go-tachart/opts"
)

type bbands struct {
	nm      string
	n       int64
	nStdDev float64
	isSma   bool
	ci      int
}

func NewBBandsSMA(n int, nStdDev float64, color int) Indicator {
	return &bbands{
		nm:      fmt.Sprintf("BBANDS(SMA, %v)", n),
		n:       int64(n),
		nStdDev: nStdDev,
		isSma:   true,
		ci:      color,
	}
}

func NewBBandsEMA(n int, nStdDev float64, color int) Indicator {
	return &bbands{
		nm:      fmt.Sprintf("BBANDS(EMA, %v)", n),
		n:       int64(n),
		nStdDev: nStdDev,
		isSma:   false,
		ci:      color,
	}
}

func (b bbands) name() string {
	return b.nm
}

func (b bbands) yAxisLabel() string {
	return ""
}

func (b bbands) yAxisMin() string {
	return ""
}

func (b bbands) yAxisMax() string {
	return ""
}

func (b bbands) getColor() int {
	return b.ci
}

func (b *bbands) getTitleOpts(top, left int) []opts.Title {
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    getColor(b.ci),
				FontSize: chartLabelFontSize,
			},
			Title: b.nm + "-Ma",
			Left:  px(left),
			Top:   px(top),
		},
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    getColor(b.ci + 1),
				FontSize: chartLabelFontSize,
			},
			Title: b.nm + "-Upper",
			Left:  px(left),
			Top:   px(top + chartLabelFontHeight),
		},
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    getColor(b.ci + 1),
				FontSize: chartLabelFontSize,
			},
			Title: b.nm + "-Lower",
			Left:  px(left),
			Top:   px(top + 2*chartLabelFontHeight),
		},
	}
}

func (b bbands) genChart(_, _, _, closes, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	var u, m, l []float64
	if b.isSma {
		u, m, l = tart.BBandsArr(tart.SMA, closes, b.n, b.nStdDev, b.nStdDev)
	} else {
		u, m, l = tart.BBandsArr(tart.EMA, closes, b.n, b.nStdDev, b.nStdDev)
	}

	uItems := []opts.LineData{}
	mItems := []opts.LineData{}
	lItems := []opts.LineData{}
	for i := 0; i < len(m); i++ {
		uItems = append(uItems, opts.LineData{Value: u[i]})
		mItems = append(mItems, opts.LineData{Value: m[i]})
		lItems = append(lItems, opts.LineData{Value: l[i]})
	}

	ml := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(b.nm+"-Ma", mItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   getColor(b.ci),
				Opacity: opacityMed,
			}))
	ul := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(b.nm+"-Upper", uItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   getColor(b.ci + 1),
				Opacity: opacityMed,
			}))
	ll := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(b.nm+"-Lower", lItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   getColor(b.ci + 1),
				Opacity: opacityMed,
			}))

	ml.Overlap(ul, ll)
	return ml
}

func (b *bbands) calcVals(closes []float64) [][]float64 {
	var upper, middle, lower []float64
	if b.isSma {
		upper, middle, lower = tart.BBandsArr(tart.SMA, closes, b.n, b.nStdDev, b.nStdDev)
	} else {
		upper, middle, lower = tart.BBandsArr(tart.EMA, closes, b.n, b.nStdDev, b.nStdDev)
	}

	return [][]float64{middle, upper, lower}
}

func (b bbands) getDrawType() string {
	return "bbands"
}
