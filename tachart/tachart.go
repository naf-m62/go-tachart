package tachart

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/naf-m62/go-tachart/charts"
	"github.com/naf-m62/go-tachart/components"
	"github.com/naf-m62/go-tachart/opts"
)

const (
	tooltipPositionFunc = `
		function(pos, params, el, elRect, size) {
			var obj = {top: 10};
			if (pos[0] > size.viewSize[0]/2) {
				obj['left'] = 30;
			} else {
				obj['right'] = 30;
			}
			return obj;
		}`
	tooltipFormatterFuncTpl = `
		function(value) {
			var eventMap = JSON.parse('__EVENT_MAP__');
			var title = (sz,txt) => '<span style="display:inline;line-height:'+(sz+2)+'px;font-size:'+sz+'px;font-weight:bold;">'+txt+'</span>';
			var square = (sz,sign,color,txt) => '<span style="display:inline;line-height:'+(sz+2)+'px;font-size:'+sz+'px;"><span style="display:inline-block;height:'+(sz+2)+'px;border-radius:3px;padding:1px 4px 1px 4px;text-align:center;margin-right:10px;background-color:' + color + ';vertical-align:top;">'+sign+'</span>'+txt+'</span>';
			var wrap = (sz,txt,width) => '<span style="display:inline-block;width:'+width+'px;word-break:break-word;word-wrap:break-word;white-space:pre-wrap;line-height:'+(sz+2)+'px;font-size:'+sz+'px;">'+txt+'</span>';
			var nowrap = (sz,txt) => '<span style="display:inline-block;line-height:'+(sz+2)+'px;font-size:'+sz+'px;">'+txt+'</span>';

			value.sort((a, b) => a.seriesIndex -b.seriesIndex);
			var cdl = value[0];
			var ret = title(14, cdl.axisValueLabel)+ '  ['+cdl.dataIndex+']' + '<br/>' +
			square(13,'O',cdl.color,cdl.value[1].toFixed(__DECIMAL_PLACES__)) + '<br/>' +
			square(13,'C',cdl.color,cdl.value[2].toFixed(__DECIMAL_PLACES__)) + '<br/>' +
			square(13,'L',cdl.color,cdl.value[3].toFixed(__DECIMAL_PLACES__)) + '<br/>' +
			square(13,'H',cdl.color,cdl.value[4].toFixed(__DECIMAL_PLACES__)) + '<br/>';
			for (var i = 1; i < value.length; i++) {
				var s = value[i];
				ret += square(13,s.seriesName,s.color,s.value.toFixed(__DECIMAL_PLACES__)) + '<br/>';
			}

			var desc = eventMap[cdl.axisValueLabel];
			if (desc) {
				if (__WRAP_DESC__) {
					ret += '<hr>' + wrap(13,desc,__WRAP_WIDTH__);
				} else {
					ret += '<hr>' + nowrap(13,desc);
				}
			}
			return ret;
		}`
	minRoundFuncTpl = `
		function(value) {
			return (value.min*0.99).toFixed(__DECIMAL_PLACES__);
		}`
	maxRoundFuncTpl = `
		function(value) {
			return (value.max*1.01).toFixed(__DECIMAL_PLACES__);
		}`
	yLabelFormatterFuncTpl = `
		function(value) {
			return value.toFixed(__DECIMAL_PLACES__);
		}`
)

var (
	ErrDuplicateCandleLabel = errors.New("candles with duplicated labels")

	// TODO: complete the map for all themes
	pageBgColorMap = map[Theme]string{
		ThemeWhite:   "#FFFFFF",
		ThemeVintage: "#FEF8EF",
	}

	// left margin
	left = 80
	// right margin
	right   = 40
	sliderH = 80
	// vertical gap between charts
	gap = 15
)

type gridLayout struct {
	top  int
	left int
	w    int
	h    int
}

type TAChart struct {
	// TODO: support dynamic auto-refresh
	cfg            Config
	globalOptsData globalOptsData
	extendedXAxis  []opts.XAxis
	extendedYAxis  []opts.YAxis
	gridLayouts    []gridLayout
}

func New(cfg Config, cdls []Candle) *TAChart {
	decimalPlaces := fmt.Sprintf("%v", cfg.precision)
	minRoundFunc := strings.Replace(minRoundFuncTpl, "__DECIMAL_PLACES__", decimalPlaces, -1)
	maxRoundFunc := strings.Replace(maxRoundFuncTpl, "__DECIMAL_PLACES__", decimalPlaces, -1)
	yLabelFormatterFunc := strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", decimalPlaces, -1)
	tooltipFormatterFunc := strings.Replace(tooltipFormatterFuncTpl, "__DECIMAL_PLACES__", decimalPlaces, -1)
	if cfg.eventDescWrapWidth == 0 {
		tooltipFormatterFunc = strings.Replace(tooltipFormatterFunc, "__WRAP_DESC__", "false", -1)
		tooltipFormatterFunc = strings.Replace(tooltipFormatterFunc, "__WRAP_WIDTH__", "0", -1)
	} else {
		tooltipFormatterFunc = strings.Replace(tooltipFormatterFunc, "__WRAP_DESC__", "true", -1)
		tooltipFormatterFunc = strings.Replace(tooltipFormatterFunc, "__WRAP_WIDTH__", fmt.Sprintf("%v", cfg.eventDescWrapWidth), -1)
	}

	// grid layuout: N = len(indicators) + 1
	// ----------------------------------------
	//   candlestick chart + overlay + events (h/2)
	// ----------------------------------------
	//   		indicator chart               (h/2/N)
	//   			...
	//   		indicator chart               (h/2/N)
	// ----------------------------------------
	//   		  volume chart                (h/2/N)
	// ----------------------------------------

	separator := 4

	h := (cfg.layout.chartHeight - sliderH) / (len(cfg.indicators) + 1 + separator)
	// candlestick+overlay
	cdlChartTop := 20
	// event
	eventChartTop := cdlChartTop + h*2 - 30
	eventChartH := 10

	grids := []opts.Grid{
		opts.Grid{ // candlestick + overlay
			Left:   px(left),
			Right:  px(right),
			Top:    px(cdlChartTop),
			Height: px(h * separator),
		},
		opts.Grid{ // event
			Left:   px(left),
			Right:  px(right),
			Top:    px(eventChartTop),
			Height: px(eventChartH),
		},
	}
	gridLayouts := []gridLayout{
		gridLayout{
			top:  cdlChartTop,
			left: left,
			w:    right - left,
			h:    h * separator,
		},
		gridLayout{
			top:  eventChartTop,
			left: left,
			w:    right - left,
			h:    eventChartH,
		},
	}
	xAxisIndex := []int{0, 1}
	extendedXAxis := []opts.XAxis{
		opts.XAxis{ // event
			Show:      false,
			GridIndex: 1,
		},
	}
	extendedYAxis := []opts.YAxis{
		opts.YAxis{ // event
			Show:      false,
			GridIndex: 1,
		},
	}

	// indicator & vol chart, inddex starting from 2
	top := cdlChartTop + h*separator + gap*2
	for i := 0; i < len(cfg.indicators)+1; i++ {
		gridIndex := i + 2
		grids = append(grids, opts.Grid{
			Left:   px(left),
			Right:  px(right),
			Top:    px(top),
			Height: px(h - gap),
		})
		gridLayouts = append(gridLayouts, gridLayout{
			top:  top,
			left: left,
			w:    right - left,
			h:    h - gap,
		})

		top += h

		xAxisIndex = append(xAxisIndex, gridIndex)

		extendedXAxis = append(extendedXAxis, opts.XAxis{
			Show:        true,
			GridIndex:   gridIndex,
			SplitNumber: 20,
			AxisTick: &opts.AxisTick{
				Show: false,
			},
			AxisLabel: &opts.AxisLabel{
				Show: false,
			},
		})
		// TODO: make this configurable
		min := minRoundFunc
		max := maxRoundFunc
		indYLabelFormatterFunc := yLabelFormatterFunc
		if i == len(cfg.indicators) {
			// volume
			min = "0"
			indYLabelFormatterFunc = strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
		} else {
			v := cfg.indicators[i].yAxisLabel()
			if v != "" {
				indYLabelFormatterFunc = v
			}
			closes := []float64{}
			for _, cdl := range cdls {
				closes = append(closes, cdl.C)
			}
			cfg.indicators[i].calcVals(closes)
			v = cfg.indicators[i].yAxisMin()
			if v != "" {
				min = v
			}
			v = cfg.indicators[i].yAxisMax()
			if v != "" {
				max = v
			}
		}

		extendedYAxis = append(extendedYAxis, opts.YAxis{
			Show:        true,
			GridIndex:   gridIndex,
			Scale:       true,
			SplitNumber: 2,
			SplitLine: &opts.SplitLine{
				Show: true,
			},
			AxisLabel: &opts.AxisLabel{
				Show:         true,
				ShowMinLabel: true,
				ShowMaxLabel: true,
				Formatter:    opts.FuncOpts(indYLabelFormatterFunc),
			},
			Min: opts.FuncOpts(min),
			Max: opts.FuncOpts(max),
		})
	}

	globalOptsData := globalOptsData{
		init: opts.Initialization{
			Theme:      string(cfg.theme),
			Width:      px(cfg.layout.chartWidth),
			Height:     px(cfg.layout.chartHeight),
			AssetsHost: cfg.assetsHost,
		},
		tooltip: opts.Tooltip{
			Show:      true,
			Trigger:   "axis",
			TriggerOn: "mousemove|click",
			Position:  opts.FuncOpts(tooltipPositionFunc),
			Formatter: opts.FuncOpts(tooltipFormatterFunc),
		},
		axisPointer: opts.AxisPointer{
			Type: "line",
			Snap: true,
			Link: opts.AxisPointerLink{
				XAxisIndex: "all",
			},
		},
		grids: grids,
		xAxis: opts.XAxis{ // candlestick+overlay
			Show:        true,
			GridIndex:   0,
			SplitNumber: 20,
		},
		yAxis: opts.YAxis{ // candlestick+overlay
			Show:      true,
			GridIndex: 0,
			Scale:     true,
			SplitArea: &opts.SplitArea{
				Show: true,
			},
			Min: opts.FuncOpts(minRoundFunc),
			Max: opts.FuncOpts(maxRoundFunc),
			AxisLabel: &opts.AxisLabel{
				Show:         true,
				ShowMinLabel: true,
				ShowMaxLabel: true,
				Formatter:    opts.FuncOpts(yLabelFormatterFunc),
			},
		},
		// dataZooms: []opts.DataZoom{
		// 	opts.DataZoom{
		// 		Type:       "slider",
		// 		Start:      50,
		// 		End:        100,
		// 		XAxisIndex: xAxisIndex,
		// 	},
		// },
	}
	if cfg.draggable {
		globalOptsData.dataZooms = append(globalOptsData.dataZooms,
			opts.DataZoom{
				Type:       "inside",
				Start:      50,
				End:        100,
				XAxisIndex: xAxisIndex,
			})
	}

	layout := gridLayouts[0]
	top = layout.top - 5
	for _, ol := range cfg.overlays {
		if ol == nil {
			continue
		}
		globalOptsData.titles = append(globalOptsData.titles, ol.getTitleOpts(top, layout.left+5)...)
		top += chartLabelFontHeight
	}
	for i, ind := range cfg.indicators {
		layout := gridLayouts[i+2]
		globalOptsData.titles = append(globalOptsData.titles, ind.getTitleOpts(layout.top-5, layout.left+5)...)
	}
	layout = gridLayouts[len(gridLayouts)-1]
	globalOptsData.titles = append(globalOptsData.titles, opts.Title{
		TitleStyle: &opts.TextStyle{
			FontSize: chartLabelFontSize,
		},
		Title: "Vol",
		Left:  px(layout.left + 5),
		Top:   px(layout.top - 5),
	})

	return &TAChart{
		cfg:            cfg,
		globalOptsData: globalOptsData,
		extendedXAxis:  extendedXAxis,
		extendedYAxis:  extendedYAxis,
		gridLayouts:    gridLayouts,
	}
}

// GenImage generates and returns chart as a PNG image byte slice
func (c TAChart) GenImage(cdls []Candle) ([]byte, error) {
	// Создаем холст изображения с размерами, соответствующими размерам графика
	width := c.cfg.layout.chartWidth
	height := c.cfg.layout.chartHeight
	dc := gg.NewContext(width, height)

	// Заполняем фон
	pageBgColor := pageBgColorMap[c.cfg.theme]
	if pageBgColor == "" {
		pageBgColor = "#FFFFFF"
	}
	bgColor, _ := parseHexColor(pageBgColor)
	dc.SetColor(bgColor)
	dc.Clear()

	// Рассчитываем высоты блоков графиков
	const topMargin = 0.0     // Отступ сверху
	const bottomMargin = 30.0 // Отступ снизу для дат
	const leftMargin = 75.0   // Отступ слева для меток оси Y
	const rightMargin = 25.0  // Отступ справа

	numIndicators := len(c.cfg.indicators) + 1 // volume indicator
	totalParts := 4 + numIndicators            // 4/7 - свечной график, по 1/7 - каждый индикатор
	usableHeight := float64(height) - bottomMargin
	partHeight := usableHeight / float64(totalParts)

	// Размеры свечного графика
	candleChartHeight := partHeight * 4
	candleChartTop := topMargin // Отступ сверху
	candleChartBottom := candleChartTop + candleChartHeight

	// Рисуем свечной график
	xAxis := make([]string, 0)
	klineSeries := []opts.KlineData{}
	opens := []float64{}
	highs := []float64{}
	lows := []float64{}
	closes := []float64{}
	vols := []float64{}
	for _, cdl := range cdls {
		xAxis = append(xAxis, cdl.Label)
		klineSeries = append(klineSeries, opts.KlineData{Value: []float64{cdl.O, cdl.C, cdl.L, cdl.H}})
		opens = append(opens, cdl.O)
		highs = append(highs, cdl.H)
		lows = append(lows, cdl.L)
		closes = append(closes, cdl.C)
		vols = append(vols, cdl.V)
	}

	// Находим минимум и максимум для масштабирования свечного графика
	min, max := findMinMax(lows, highs)
	canvasWidth := float64(width) - 100.0    // Оставляем место для осей
	canvasHeight := candleChartHeight - 25.0 // Высота свечного графика с учетом отступов

	// Рисуем рамку для свечного графика
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(1)
	dc.DrawLine(leftMargin, candleChartTop, leftMargin, candleChartBottom)                    // Вертикальная ось Y
	dc.DrawLine(leftMargin, candleChartBottom, float64(width)-rightMargin, candleChartBottom) // Горизонтальная ось X
	dc.Stroke()

	// Рисуем свечи
	candleWidth := canvasWidth / float64(len(cdls))
	candleBarWidth := candleWidth * 0.6

	// каждые 120px рисуем метку на оси X
	// высчитываем сколько свечей в 120 px
	candleCount := int(120 / candleWidth)

	for i, cdl := range cdls {
		x := leftMargin + float64(i)*candleWidth + candleWidth/2.0
		// Добавляем отступ сверху к координатам
		yOpen := mapValueToCanvas(cdl.O, min, max, canvasHeight) + candleChartTop
		yClose := mapValueToCanvas(cdl.C, min, max, canvasHeight) + candleChartTop
		yLow := mapValueToCanvas(cdl.L, min, max, canvasHeight) + candleChartTop
		yHigh := mapValueToCanvas(cdl.H, min, max, canvasHeight) + candleChartTop

		if cdl.O > cdl.C {
			// Медвежья свеча (красная)
			dc.SetRGB(1, 0, 0)
		} else {
			// Бычья свеча (зеленая)
			dc.SetRGB(0, 1, 0)
		}

		// Рисуем линию от минимума до максимума
		dc.DrawLine(x, yLow, x, yHigh)
		dc.Stroke()

		// Рисуем тело свечи
		yBodyTop := yOpen
		yBodyBottom := yClose
		if yOpen > yClose {
			yBodyTop = yClose
			yBodyBottom = yOpen
		}

		if cdl.O == cdl.C {
			// Рисуем линию
			dc.DrawLine(x-candleBarWidth/2.0, yBodyTop, x+candleBarWidth/2.0, yBodyTop)
			dc.Stroke()
		} else {
			dc.DrawRectangle(x-candleBarWidth/2.0, yBodyTop, candleBarWidth, yBodyBottom-yBodyTop)
			dc.Fill()
		}

		// Каждые candleCount свечей рисуем метку на оси X
		if i%candleCount == 0 && i < len(cdls)-1 {
			// рисуем линию над меткой
			dc.SetRGB(0, 0, 0)
			dc.DrawLine(x, float64(height)-bottomMargin, x, float64(height)-bottomMargin+10)
			dc.DrawString(cdls[i].Label, x, float64(height)-bottomMargin+20)
			dc.Stroke()
		}
	}

	// Рисуем значения на оси Y для свечного графика
	steps := 5
	for i := 0; i <= steps; i++ {
		value := min + (max-min)*float64(i)/float64(steps)
		y := mapValueToCanvas(value, min, max, canvasHeight)
		y += candleChartTop // Прибавляем отступ сверху
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(fmt.Sprintf("%.2f", value), leftMargin-5, y, 1.0, 0.5)

		// Рисуем горизонтальные линии сетки
		dc.SetRGBA(0, 0, 0, 0.2)
		dc.DrawLine(leftMargin, y, float64(width)-rightMargin, y)
		dc.Stroke()
	}

	// Рисуем overlays
	for _, ol := range c.cfg.overlays {
		if ol == nil {
			continue
		}

		// Получаем значения индикатора
		vals := ol.calcVals(closes) // Используем closes как базовый набор данных
		if len(vals) == 0 {
			continue // Пропускаем, если нет данных
		}

		// Определяем тип отрисовки
		drawType := ol.getDrawType()

		// Настраиваем цвет для индикатора

		color, err := parseHexColor(getColor(ol.getColor()))
		if err != nil {
			// Цвет по умолчанию - синий
			dc.SetRGB(0.0, 0.5, 1.0)
		} else {
			dc.SetColor(color)
		}

		// Отрисовываем в зависимости от типа
		switch drawType {
		case "line":
			// Рисуем линию
			dc.SetLineWidth(1.5) // Делаем линию чуть толще, чем у свечей

			// Находим min и max для индикатора, если нужно масштабирование
			// В этом примере используем тот же масштаб, что и для свечей

			// Рисуем линию, соединяя действительные точки (не NaN)
			for i := 0; i < len(vals) && i < len(cdls); i++ {
				var lastValidX, lastValidY float64
				var hasLastValid bool

				// Проходим по всем точкам линии
				for j := 0; j < len(vals[i]); j++ {
					// Текущая координата X
					x := leftMargin + float64(j)*candleWidth + candleWidth/2.0

					// Проверяем, является ли текущее значение действительным (не NaN)
					value := vals[i][j]
					if !math.IsNaN(value) {
						// Рассчитываем координату Y для действительного значения
						y := mapValueToCanvas(value, min, max, canvasHeight) + candleChartTop

						// Если есть предыдущая действительная точка, соединяем с ней линией
						if hasLastValid {
							dc.DrawLine(lastValidX, lastValidY, x, y)
						}

						// Запоминаем эту точку как последнюю действительную
						lastValidX = x
						lastValidY = y
						hasLastValid = true
					}
				}
			}
			dc.Stroke()

		case "bar":
			// Рисуем бары (столбцы)
			dc.SetLineWidth(1)

			// Для баров используем меньшую ширину, чем для свечей
			barWidth := candleWidth * 0.4

			for i := 0; i < len(vals) && i < len(cdls); i++ {
				x := leftMargin + float64(i)*candleWidth + candleWidth/2.0
				y := mapValueToCanvas(vals[0][i], min, max, canvasHeight) + candleChartTop

				// Рисуем бар от оси Y до значения
				baseY := mapValueToCanvas(0, min, max, canvasHeight) + candleChartTop
				dc.DrawRectangle(x-barWidth/2.0, math.Min(y, baseY), barWidth, math.Abs(y-baseY))
				dc.Fill()
			}

		case "macd":
			// Для MACD нужна специальная обработка, так как у него несколько линий
			// В этом примере просто рисуем основную линию
			dc.SetLineWidth(1.5)

			for i := 0; i < len(vals) && i < len(cdls); i++ {
				for j := 1; j < len(vals[i]); j++ {
					if j != 2 {
						x1 := leftMargin + float64(j-1)*candleWidth + candleWidth/2.0
						x2 := leftMargin + float64(j)*candleWidth + candleWidth/2.0

						y1 := mapValueToCanvas(vals[i][j-1], min, max, canvasHeight)
						y2 := mapValueToCanvas(vals[i][j], min, max, canvasHeight)
						dc.DrawLine(x1, y1, x2, y2)
					} else {
						// histogram
						dc.SetLineWidth(1)
						barWidth := candleWidth * 0.6
						x := leftMargin + float64(i)*candleWidth + candleWidth/2.0
						y := mapValueToCanvas(vals[i][j], min, max, canvasHeight) + candleChartTop
						baseY := mapValueToCanvas(0, min, max, canvasHeight) + candleChartTop
						dc.DrawRectangle(x-barWidth/2.0, math.Min(y, baseY), barWidth, math.Abs(y-baseY))
						dc.Fill()
					}
				}
			}
			dc.Stroke()

		default:
			// По умолчанию просто рисуем линию
			dc.SetLineWidth(1.5)

			for i := 1; i < len(vals) && i < len(cdls); i++ {
				for j := 0; j < len(vals[i]); j++ {
					x1 := leftMargin + float64(i-1)*candleWidth + candleWidth/2.0
					x2 := leftMargin + float64(i)*candleWidth + candleWidth/2.0

					y1 := mapValueToCanvas(vals[i-1][j], min, max, canvasHeight)
					y2 := mapValueToCanvas(vals[i][j], min, max, canvasHeight)
					dc.DrawLine(x1, y1, x2, y2)
				}
			}
			dc.Stroke()
		}
	}

	// Рисуем индикаторы на отдельных графиках ниже основного
	for indIdx, ind := range c.cfg.indicators {
		if ind == nil {
			continue
		}

		// Вычисляем координаты и размеры графика индикатора
		indTop := candleChartBottom + 20.0 + float64(indIdx)*partHeight
		indBottom := indTop + partHeight - 20.0
		indHeight := indBottom - indTop

		// Рисуем рамку для графика индикатора
		dc.SetRGB(0, 0, 0)
		dc.SetLineWidth(1)
		dc.DrawLine(leftMargin, indTop, leftMargin, indBottom)                    // Вертикальная ось Y
		dc.DrawLine(leftMargin, indBottom, float64(width)-rightMargin, indBottom) // Горизонтальная ось X
		dc.Stroke()

		// Добавляем название индикатора
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(ind.name(), leftMargin+10, indTop+10, 0, 0.5)

		// Получаем значения индикатора
		vals := ind.calcVals(closes) // Используем closes как базовый набор данных
		if len(vals) == 0 {
			continue // Пропускаем, если нет данных
		}

		// Находим минимум и максимум для масштабирования индикатора
		// Рисуем значения на оси Y для индикатора
		var indMin, indMax float64
		switch {
		case strings.HasPrefix(ind.name(), "RSI"):
			indMin = 0
			indMax = 100
			for _, i := range []int{30, 70} {
				y := indTop + indHeight - (indHeight * float64(i) / float64(indMax))
				dc.SetRGB(0, 0, 0)
				dc.DrawStringAnchored(fmt.Sprintf("%v", i), leftMargin-5, y, 1.0, 0.5)

				// Рисуем горизонтальные линии сетки
				dc.SetRGBA(0, 0, 0, 0.2)
				dc.DrawLine(leftMargin, y, float64(width)-rightMargin, y)
				dc.Stroke()
			}
		case strings.HasPrefix(ind.name(), "MACD"):
			indMin, indMax = findMinMax(vals...)

			// Минимальное значение
			y := indBottom
			dc.SetRGB(0, 0, 0)
			dc.DrawStringAnchored(fmt.Sprintf("%.2f", indMin), leftMargin-5, y, 1.0, 0.5)
			dc.SetRGBA(0, 0, 0, 0.2)
			dc.DrawLine(leftMargin, y, float64(width)-rightMargin, y)
			dc.Stroke()

			// Нулевое значение
			value := 0.0
			y = indTop + indHeight - (indHeight * (value - indMin) / (indMax - indMin))
			dc.SetRGB(0, 0, 0)
			dc.DrawStringAnchored(fmt.Sprintf("%.2f", value), leftMargin-5, y, 1.0, 0.5)
			dc.SetRGBA(0, 0, 0, 0.2)
			dc.DrawLine(leftMargin, y, float64(width)-rightMargin, y)
			dc.Stroke()

			// Максимальное значение
			y = indTop
			dc.SetRGB(0, 0, 0)
			dc.DrawStringAnchored(fmt.Sprintf("%.2f", indMax), leftMargin-5, y, 1.0, 0.5)
			dc.SetRGBA(0, 0, 0, 0.2)
			dc.DrawLine(leftMargin, y, float64(width)-rightMargin, y)
			dc.Stroke()
		default:
			indMin, indMax = findMinMax(vals...)
			steps := 3
			for i := 0; i <= steps; i++ {
				value := indMin + (indMax-indMin)*float64(i)/float64(steps)
				y := indTop + indHeight - (indHeight * float64(i) / float64(steps))
				dc.SetRGB(0, 0, 0)
				dc.DrawStringAnchored(fmt.Sprintf("%.2f", value), leftMargin-5, y, 1.0, 0.5)

				// Рисуем горизонтальные линии сетки
				dc.SetRGBA(0, 0, 0, 0.2)
				dc.DrawLine(leftMargin, y, float64(width)-rightMargin, y)
				dc.Stroke()
			}
		}

		// Определяем тип отрисовки
		drawType := ind.getDrawType()

		// Настраиваем цвет для индикатора
		// Цвет по умолчанию - зеленый
		dc.SetRGB(0.0, 0.8, 0.0)

		// Отрисовываем в зависимости от типа
		switch drawType {
		case "line":
			// Рисуем линию
			dc.SetLineWidth(1.5)

			for i := 0; i < len(vals) && i < len(cdls); i++ {
				for j := 1; j < len(vals[i]); j++ {
					x1 := leftMargin + float64(j-1)*candleWidth + candleWidth/2.0
					x2 := leftMargin + float64(j)*candleWidth + candleWidth/2.0

					// Отображаем значения в пределах графика индикатора
					y1 := indTop + indHeight - (indHeight * (vals[i][j-1] - indMin) / (indMax - indMin))
					y2 := indTop + indHeight - (indHeight * (vals[i][j] - indMin) / (indMax - indMin))

					dc.DrawLine(x1, y1, x2, y2)
				}
			}
			dc.Stroke()

		case "bar":
			// Рисуем бары (столбцы)
			dc.SetLineWidth(1)

			// Для баров используем меньшую ширину, чем для свечей
			barWidth := candleWidth * 0.4

			for i := 0; i < len(vals) && i < len(cdls); i++ {
				x := leftMargin + float64(i)*candleWidth + candleWidth/2.0

				// Отображаем значения в пределах графика индикатора
				y := indTop + indHeight - (indHeight * (vals[0][i] - indMin) / (indMax - indMin))
				baseY := indBottom

				dc.DrawRectangle(x-barWidth/2.0, math.Min(y, baseY), barWidth, math.Abs(y-baseY))
				dc.Fill()
			}

		case "macd":
			// Для MACD нужна специальная обработка, так как у него несколько линий

			// Сначала отрисовываем гистограмму (индекс 2)
			if len(vals) >= 3 && len(vals[2]) > 0 { // Проверяем, есть ли данные гистограммы
				dc.SetLineWidth(1)
				barWidth := candleWidth * 0.4

				// Вычисляем базовую линию для нуля
				baseY := indTop + indHeight - (indHeight * (0 - indMin) / (indMax - indMin))

				for j := 0; j < len(vals[2]) && j < len(cdls); j++ {
					x := leftMargin + float64(j)*candleWidth + candleWidth/2.0
					y := indTop + indHeight - (indHeight * (vals[2][j] - indMin) / (indMax - indMin))

					// Цвет баров в зависимости от значения
					if vals[2][j] >= 0 {
						dc.SetRGB(0, 0.7, 0) // Зеленый для положительных значений
					} else {
						dc.SetRGB(0.7, 0, 0) // Красный для отрицательных
					}

					dc.DrawRectangle(x-barWidth/2.0, math.Min(y, baseY), barWidth, math.Abs(y-baseY))
					dc.Fill()
				}
			}

			// Затем отрисовываем линии MACD и сигнальную линию
			// MACD линия - синяя
			if len(vals) >= 1 && len(vals[0]) > 1 {
				dc.SetLineWidth(1.5)
				dc.SetRGB(0.0, 0.0, 0.8) // Синий цвет для MACD

				for j := 1; j < len(vals[0]) && j < len(cdls); j++ {
					x1 := leftMargin + float64(j-1)*candleWidth + candleWidth/2.0
					x2 := leftMargin + float64(j)*candleWidth + candleWidth/2.0

					y1 := indTop + indHeight - (indHeight * (vals[0][j-1] - indMin) / (indMax - indMin))
					y2 := indTop + indHeight - (indHeight * (vals[0][j] - indMin) / (indMax - indMin))

					dc.DrawLine(x1, y1, x2, y2)
				}
				dc.Stroke()
			}

			// Сигнальная линия - красная
			if len(vals) >= 2 && len(vals[1]) > 1 {
				dc.SetLineWidth(1.5)
				dc.SetRGB(0.8, 0.0, 0.0) // Красный цвет для сигнальной линии

				for j := 1; j < len(vals[1]) && j < len(cdls); j++ {
					x1 := leftMargin + float64(j-1)*candleWidth + candleWidth/2.0
					x2 := leftMargin + float64(j)*candleWidth + candleWidth/2.0

					y1 := indTop + indHeight - (indHeight * (vals[1][j-1] - indMin) / (indMax - indMin))
					y2 := indTop + indHeight - (indHeight * (vals[1][j] - indMin) / (indMax - indMin))

					dc.DrawLine(x1, y1, x2, y2)
				}
				dc.Stroke()
			}
			dc.Stroke()

		default:
			// По умолчанию просто рисуем линию
			dc.SetLineWidth(1.5)

			for i := 1; i < len(vals) && i < len(cdls); i++ {
				x1 := leftMargin + float64(i-1)*candleWidth + candleWidth/2.0
				x2 := leftMargin + float64(i)*candleWidth + candleWidth/2.0

				// Отображаем значения в пределах графика индикатора
				y1 := indTop + indHeight - (indHeight * (vals[0][i-1] - indMin) / (indMax - indMin))
				y2 := indTop + indHeight - (indHeight * (vals[0][i] - indMin) / (indMax - indMin))

				dc.DrawLine(x1, y1, x2, y2)
			}
			dc.Stroke()
		}
	}

	// Рисуем индикатор Volume в конце
	// Вычисляем координаты и размеры графика Volume
	// Объявляем основные переменные для Volume
	var (
		volIdx    int
		volTop    float64
		volBottom float64
		volHeight float64
		volMinVal float64
		volMaxVal float64
		barWidth  float64 // Ширина баров объема
	)

	volIdx = len(c.cfg.indicators)
	volTop = candleChartBottom + 5.0 + float64(volIdx)*partHeight
	volBottom = volTop + partHeight - 5.0
	volHeight = volBottom - volTop

	// Рисуем рамку для графика Volume
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(1)
	dc.DrawLine(leftMargin, volTop, leftMargin, volBottom)                    // Вертикальная ось Y
	dc.DrawLine(leftMargin, volBottom, float64(width)-rightMargin, volBottom) // Горизонтальная ось X
	dc.Stroke()

	// Добавляем название Volume
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Volume", leftMargin+10, volTop+10, 0, 0.5)

	// Находим минимум и максимум для масштабирования Volume
	volMinVal, volMaxVal = findMinMax([]float64(vols))

	// Рисуем значения на оси Y для Volume
	volSteps := 3
	for i := 0; i <= volSteps; i++ {
		value := volMinVal + (volMaxVal-volMinVal)*float64(i)/float64(volSteps)
		y := volTop + volHeight - (volHeight * float64(i) / float64(volSteps))
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(fmt.Sprintf("%.0f", value), leftMargin-5, y, 1.0, 0.5)

		// Рисуем горизонтальные линии сетки
		dc.SetRGBA(0, 0, 0, 0.2)
		dc.DrawLine(leftMargin, y, float64(width)-rightMargin, y)
		dc.Stroke()
	}

	// Рисуем бары объема
	barWidth = candleWidth * 0.6

	for i := 0; i < len(vols) && i < len(cdls); i++ {
		x := leftMargin + float64(i)*candleWidth + candleWidth/2.0

		// Отображаем значения в пределах графика Volume
		var y float64 = volTop + volHeight - (volHeight * (vols[i] - volMinVal) / (volMaxVal - volMinVal))
		baseY := volBottom

		// Цвет бара зависит от направления свечи (растущая или падающая)
		if i > 0 && cdls[i].C > cdls[i].O {
			// Бычья свеча (зеленая)
			dc.SetRGB(0, 0.8, 0)
		} else {
			// Медвежья свеча (красная)
			dc.SetRGB(0.8, 0, 0)
		}

		dc.DrawRectangle(x-barWidth/2.0, math.Min(y, baseY), barWidth, math.Abs(y-baseY))
		dc.Fill()
	}

	// Возвращаем изображение как слайс байтов
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, dc.Image())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// parseHexColor парсит HTML-цветовой код в color.RGBA
func parseHexColor(hexColor string) (color.RGBA, error) {
	var r, g, b uint8
	hexColor = strings.TrimPrefix(hexColor, "#")

	if len(hexColor) == 6 {
		n, err := fmt.Sscanf(hexColor, "%02x%02x%02x", &r, &g, &b)
		if err != nil || n != 3 {
			return color.RGBA{}, err
		}
	} else {
		return color.RGBA{}, fmt.Errorf("invalid hex color: %s", hexColor)
	}

	return color.RGBA{r, g, b, 255}, nil
}

// findMinMax находит минимальное и максимальное значение в массивах
func findMinMax(arrays ...[]float64) (min, max float64) {
	min = float64(^uint(0) >> 1) // Max value for int
	max = -min

	for _, array := range arrays {
		for _, v := range array {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}

	// Добавляем немного запаса для лучшего отображения
	padding := (max - min) * 0.05
	min -= padding
	max += padding

	return
}

// mapValueToCanvas преобразует значение из диапазона данных в координату на холсте
func mapValueToCanvas(value, min, max, canvasHeight float64) float64 {
	// Обратите внимание, что координата Y в canvas начинается сверху,
	// поэтому нам нужно инвертировать значение
	return 25.0 + canvasHeight - (value-min)/(max-min)*canvasHeight
}

func (c TAChart) GenStatic(cdls []Candle, events []Event, path string) error {
	xAxis := make([]string, 0)
	klineSeries := []opts.KlineData{}
	volSeries := []opts.BarData{}
	opens := []float64{}
	highs := []float64{}
	lows := []float64{}
	closes := []float64{}
	vols := []float64{}
	cdlMap := map[string]*Candle{}
	for _, cdl := range cdls {
		xAxis = append(xAxis, cdl.Label)
		// open,close,low,high
		klineSeries = append(klineSeries, opts.KlineData{Value: []float64{cdl.O, cdl.C, cdl.L, cdl.H}})
		opens = append(opens, cdl.O)
		highs = append(highs, cdl.H)
		lows = append(lows, cdl.L)
		closes = append(closes, cdl.C)
		vols = append(vols, cdl.V)

		style := &opts.ItemStyle{
			Color:   colorUpBar,
			Opacity: opacityHeavy,
		}
		if cdl.O > cdl.C {
			style = &opts.ItemStyle{
				Color:   colorDownBar,
				Opacity: opacityHeavy,
			}
		}
		volSeries = append(volSeries, opts.BarData{
			Value:     cdl.V,
			ItemStyle: style,
		})

		if cdlMap[cdl.Label] != nil {
			return ErrDuplicateCandleLabel
		}
		c := cdl
		cdlMap[cdl.Label] = &c
	}

	// candlestick+overlay
	chart := charts.NewKLine().SetXAxis(xAxis).AddSeries("kline",
		klineSeries,
		charts.WithKlineChartOpts(opts.KlineChart{
			BarWidth:   "60%",
			XAxisIndex: 0,
			YAxisIndex: 0,
		}),
		charts.WithItemStyleOpts(opts.ItemStyle{
			Color:        colorUpBar,
			Color0:       colorDownBar,
			BorderColor:  colorUpBar,
			BorderColor0: colorDownBar,
			Opacity:      opacityHeavy,
		}),
	)

	eventDescMap := map[string]string{}
	for _, e := range events {
		eventDescMap[e.Label] = e.Description
	}

	chart.SetGlobalOptions(c.globalOptsData.genOpts(c.cfg, len(cdls), eventDescMap)...)

	for _, ol := range c.cfg.overlays {
		if ol == nil {
			continue
		}
		chart.Overlap(ol.genChart(opens, highs, lows, closes, vols, xAxis, 0))
	}

	for i := 0; i < len(c.extendedXAxis); i++ {
		c.extendedXAxis[i].Data = xAxis
	}
	chart.ExtendXAxis(c.extendedXAxis...)
	chart.ExtendYAxis(c.extendedYAxis...)

	evtOpts := []charts.SeriesOpts{
		charts.WithBarChartOpts(opts.BarChart{
			BarWidth:   "60%",
			XAxisIndex: 1,
			YAxisIndex: 1,
		}),
	}
	for _, e := range events {
		es := eventLabelMap[e.Type]
		if e.Type == CustomEvent {
			es = e.EventMark.toEventStyle()
		}
		evtOpts = append(evtOpts, charts.WithMarkPointNameCoordItemOpts(opts.MarkPointNameCoordItem{
			Symbol:     "roundRect",
			SymbolSize: es.symbolSize,
			Coordinate: []interface{}{e.Label, 0},
			Label:      es.label,
			ItemStyle:  es.style,
		}))
	}
	event := charts.NewBar().AddSeries("events", []opts.BarData{}, evtOpts...)
	chart.Overlap(event)

	// grid index starting from 2 (candlestick+event)
	for i, ind := range c.cfg.indicators {
		chart.Overlap(ind.genChart(opens, highs, lows, closes, vols, xAxis, i+2))
	}

	bar := charts.NewBar().
		SetXAxis(xAxis).
		AddSeries("Vol", volSeries, charts.WithBarChartOpts(opts.BarChart{
			BarWidth:   "60%",
			XAxisIndex: len(c.cfg.indicators) + 2,
			YAxisIndex: len(c.cfg.indicators) + 2,
		}))
	chart.Overlap(bar)
	chart.AddJSFuncs(c.cfg.jsFuncs...)

	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	layout := components.Layout{
		TemplateColumns: template.CSS(fmt.Sprintf("%vpx %vpx %vpx", c.cfg.layout.leftWidth, c.cfg.layout.chartWidth, c.cfg.layout.rightWidth)),
		TopHeight:       template.CSS(px(c.cfg.layout.topHeight)),
		BottomHeight:    template.CSS(px(c.cfg.layout.bottomHeight)),
		TopContent:      template.HTML(c.cfg.layout.topContent),
		BottomContent:   template.HTML(c.cfg.layout.bottomContent),
		LeftContent:     template.HTML(c.cfg.layout.leftContent),
		RightContent:    template.HTML(c.cfg.layout.rightContent),
	}

	pageBgColor := pageBgColorMap[c.cfg.theme]
	if pageBgColor == "" {
		pageBgColor = "#FFFFFF"
	}

	return components.NewPage(c.cfg.assetsHost).
		SetLayout(layout).
		SetBackgroundColor(pageBgColor).
		AddCharts(chart).
		Render(fp)
}

func px(v int) string {
	return fmt.Sprintf("%vpx", v)
}
