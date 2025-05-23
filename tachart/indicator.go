package tachart

import (
	"github.com/naf-m62/go-tachart/charts"
	"github.com/naf-m62/go-tachart/opts"
)

const (
	chartLabelFontSize   = 11
	chartLabelFontHeight = 13
)

type Indicator interface {
	// indicator name
	name() string
	// y axis label formatter
	yAxisLabel() string
	// y axis min label formatter
	yAxisMin() string
	// y axis max label formatter
	yAxisMax() string
	// # of colors needed
	getColor() int
	// indicator chart legend config
	getTitleOpts(top, left int) []opts.Title
	// indicator chart config
	genChart(opens, highs, lows, closes, vols []float64, xAxis interface{}, gridIndex int) charts.Overlaper
	// calculate indicator values
	calcVals(vals []float64) [][]float64
	// get draw type (line, bars, etc)
	getDrawType() string
}
