package charts

import (
	"github.com/naf-m62/go-tachart/opts"
	"github.com/naf-m62/go-tachart/render"
	"github.com/naf-m62/go-tachart/types"
)

// Sankey represents a sankey chart.
type Sankey struct {
	BaseConfiguration
}

// Type returns the chart type.
func (Sankey) Type() string { return types.ChartSankey }

// NewSankey creates a new sankey chart.
func NewSankey() *Sankey {
	c := &Sankey{}
	c.initBaseConfiguration()
	c.Renderer = render.NewChartRender(c, c.Validate)
	return c
}

// AddSeries adds new data sets.
func (c *Sankey) AddSeries(name string, nodes []opts.SankeyNode, links []opts.SankeyLink, options ...SeriesOpts) *Sankey {
	series := SingleSeries{Name: name, Type: types.ChartSankey, Data: nodes, Links: links}
	series.configureSeriesOpts(options...)
	c.MultiSeries = append(c.MultiSeries, series)
	return c
}

// SetGlobalOptions sets options for the Sankey instance.
func (c *Sankey) SetGlobalOptions(options ...GlobalOpts) *Sankey {
	c.BaseConfiguration.setBaseGlobalOptions(options...)
	return c
}

// Validate
func (c *Sankey) Validate() {
	c.Assets.Validate(c.AssetsHost)
}
