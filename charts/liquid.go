package charts

import (
	"github.com/naf-m62/go-tachart/opts"
	"github.com/naf-m62/go-tachart/render"
	"github.com/naf-m62/go-tachart/types"
)

// Liquid represents a liquid chart.
type Liquid struct {
	BaseConfiguration
}

// Type returns the chart type.
func (Liquid) Type() string { return types.ChartLiquid }

// NewLiquid creates a new liquid chart.
func NewLiquid() *Liquid {
	c := &Liquid{}
	c.initBaseConfiguration()
	c.Renderer = render.NewChartRender(c, c.Validate)
	c.JSAssets.Add("echarts-liquidfill.min.js")
	return c
}

// AddSeries adds new data sets.
func (c *Liquid) AddSeries(name string, data []opts.LiquidData, options ...SeriesOpts) *Liquid {
	series := SingleSeries{Name: name, Type: types.ChartLiquid, Data: data}
	series.configureSeriesOpts(options...)
	c.MultiSeries = append(c.MultiSeries, series)
	return c
}

// SetGlobalOptions sets options for the Liquid instance.
func (c *Liquid) SetGlobalOptions(options ...GlobalOpts) *Liquid {
	c.BaseConfiguration.setBaseGlobalOptions(options...)
	return c
}

// Validate
func (c *Liquid) Validate() {
	c.Assets.Validate(c.AssetsHost)
}
