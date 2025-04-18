package charts

import (
	"github.com/naf-m62/go-tachart/opts"
	"github.com/naf-m62/go-tachart/render"
	"github.com/naf-m62/go-tachart/types"
)

// Radar represents a radar chart.
type Radar struct {
	BaseConfiguration
}

// Type returns the chart type.
func (Radar) Type() string { return types.ChartRadar }

// NewRadar creates a new radar chart.
func NewRadar() *Radar {
	c := &Radar{}
	c.initBaseConfiguration()
	c.Renderer = render.NewChartRender(c, c.Validate)
	c.hasRadar = true
	return c
}

// AddSeries adds new data sets.
func (c *Radar) AddSeries(name string, data []opts.RadarData, options ...SeriesOpts) *Radar {
	series := SingleSeries{Name: name, Type: types.ChartRadar, Data: data}
	series.configureSeriesOpts(options...)
	c.MultiSeries = append(c.MultiSeries, series)
	c.legends = append(c.legends, name)
	return c
}

// SetGlobalOptions sets options for the Radar instance.
func (c *Radar) SetGlobalOptions(options ...GlobalOpts) *Radar {
	c.BaseConfiguration.setBaseGlobalOptions(options...)
	return c
}

// Validate
func (c *Radar) Validate() {
	c.Legend.Data = c.legends
	c.Assets.Validate(c.AssetsHost)
}
