package main

import (
	"github.com/naf-m62/go-tachart/tachart"
)

var (
	cdls = []tachart.Candle{
		{Label: "2018/1/24", O: 2320.26, C: 2320.26, L: 2287.3, H: 2362.94, V: 149092},
		{Label: "2018/1/25", O: 2300, C: 2291.3, L: 2288.26, H: 2308.38, V: 189092},
		{Label: "2018/1/28", O: 2295.35, C: 2346.5, L: 2295.35, H: 2346.92, V: 159034},
		{Label: "2018/1/29", O: 2347.22, C: 2358.98, L: 2337.35, H: 2363.8, V: 249910},
		{Label: "2018/1/30", O: 2360.75, C: 2382.48, L: 2347.89, H: 2383.76, V: 119910},
		{Label: "2018/1/31", O: 2383.43, C: 2385.42, L: 2371.23, H: 2391.82, V: 89940},
		{Label: "2018/2/1", O: 2377.41, C: 2419.02, L: 2369.57, H: 2421.15, V: 192941},
		{Label: "2018/2/4", O: 2425.92, C: 2428.15, L: 2417.58, H: 2440.38, V: 249410},
		{Label: "2018/2/5", O: 2411, C: 2433.13, L: 2403.3, H: 2437.42, V: 149410},
		{Label: "2018/2/6", O: 2432.68, C: 2434.48, L: 2427.7, H: 2441.73, V: 149910},
		{Label: "2018/2/7", O: 2430.69, C: 2418.53, L: 2394.22, H: 2433.89, V: 249910},
		{Label: "2018/2/8", O: 2416.62, C: 2432.4, L: 2414.4, H: 2443.03, V: 149410},
		{Label: "2018/2/18", O: 2441.91, C: 2421.56, L: 2415.43, H: 2444.8, V: 249910},
		{Label: "2018/2/19", O: 2420.26, C: 2382.91, L: 2373.53, H: 2427.07, V: 149910},
		{Label: "2018/2/20", O: 2383.49, C: 2397.18, L: 2370.61, H: 2397.94, V: 149910},
		{Label: "2018/2/21", O: 2378.82, C: 2325.95, L: 2309.17, H: 2378.82, V: 449910},
		{Label: "2018/2/22", O: 2322.94, C: 2314.16, L: 2308.76, H: 2330.88, V: 149910},
		{Label: "2018/2/25", O: 2320.62, C: 2325.82, L: 2315.01, H: 2338.78, V: 249910},
		{Label: "2018/2/26", O: 2313.74, C: 2293.34, L: 2289.89, H: 2340.71, V: 249940},
		{Label: "2018/2/27", O: 2297.77, C: 2313.22, L: 2292.03, H: 2324.63, V: 149944},
		{Label: "2018/2/28", O: 2322.32, C: 2365.59, L: 2308.92, H: 2366.16, V: 249910},
		{Label: "2018/3/1", O: 2364.54, C: 2359.51, L: 2330.86, H: 2369.65, V: 249914},
		{Label: "2018/3/4", O: 2332.08, C: 2273.4, L: 2259.25, H: 2333.54, V: 241910},
		{Label: "2018/3/5", O: 2274.81, C: 2326.31, L: 2270.1, H: 2328.14, V: 249910},
		{Label: "2018/3/6", O: 2333.61, C: 2347.18, L: 2321.6, H: 2351.44, V: 241911},
		{Label: "2018/3/7", O: 2340.44, C: 2324.29, L: 2304.27, H: 2352.02, V: 249910},
		{Label: "2018/3/8", O: 2326.42, C: 2318.61, L: 2314.59, H: 2333.67, V: 249110},
		{Label: "2018/3/11", O: 2314.68, C: 2310.59, L: 2296.58, H: 2320.96, V: 249910},
		{Label: "2018/3/12", O: 2309.16, C: 2286.6, L: 2264.83, H: 2333.29, V: 241940},
		{Label: "2018/3/13", O: 2282.17, C: 2263.97, L: 2253.25, H: 2286.33, V: 249911},
		{Label: "2018/3/14", O: 2255.77, C: 2270.28, L: 2253.31, H: 2276.22, V: 149110},
		{Label: "2018/3/15", O: 2269.31, C: 2278.4, L: 2250, H: 2312.08, V: 149911},
		{Label: "2018/3/18", O: 2267.29, C: 2240.02, L: 2239.21, H: 2276.05, V: 249410},
		{Label: "2018/3/19", O: 2244.26, C: 2257.43, L: 2232.02, H: 2261.31, V: 149910},
		{Label: "2018/3/20", O: 2257.74, C: 2317.37, L: 2257.42, H: 2317.86, V: 249910},
		{Label: "2018/3/21", O: 2318.21, C: 2324.24, L: 2311.6, H: 2330.81, V: 219910},
		{Label: "2018/3/22", O: 2321.4, C: 2328.28, L: 2314.97, H: 2332, V: 249911},
		{Label: "2018/3/25", O: 2334.74, C: 2326.72, L: 2319.91, H: 2344.89, V: 249410},
		{Label: "2018/3/26", O: 2318.58, C: 2297.67, L: 2281.12, H: 2319.99, V: 149910},
		{Label: "2018/3/27", O: 2299.38, C: 2301.26, L: 2289, H: 2323.48, V: 89910},
		{Label: "2018/3/28", O: 2273.55, C: 2236.3, L: 2232.91, H: 2273.55, V: 240910},
		{Label: "2018/3/29", O: 2238.49, C: 2236.62, L: 2228.81, H: 2246.87, V: 249410},
		{Label: "2018/4/1", O: 2229.46, C: 2234.4, L: 2227.31, H: 2243.95, V: 249010},
		{Label: "2018/4/2", O: 2234.9, C: 2227.74, L: 2220.44, H: 2253.42, V: 249910},
		{Label: "2018/4/3", O: 2232.69, C: 2225.29, L: 2217.25, H: 2241.34, V: 219910},
		{Label: "2018/4/8", O: 2196.24, C: 2211.59, L: 2180.67, H: 2212.59, V: 149940},
		{Label: "2018/4/9", O: 2215.47, C: 2225.77, L: 2215.47, H: 2234.73, V: 191410},
		{Label: "2018/4/10", O: 2224.93, C: 2226.13, L: 2212.56, H: 2233.04, V: 149910},
		{Label: "2018/4/11", O: 2236.98, C: 2219.55, L: 2217.26, H: 2242.48, V: 219410},
		{Label: "2018/4/12", O: 2218.09, C: 2206.78, L: 2204.44, H: 2226.26, V: 149910},
		{Label: "2018/4/15", O: 2199.91, C: 2181.94, L: 2177.39, H: 2204.99, V: 249911},
		{Label: "2018/4/16", O: 2169.63, C: 2194.85, L: 2165.78, H: 2196.43, V: 279417},
		{Label: "2018/4/17", O: 2195.03, C: 2193.8, L: 2178.47, H: 2197.51, V: 179940},
		{Label: "2018/4/18", O: 2181.82, C: 2197.6, L: 2175.44, H: 2206.03, V: 249940},
		{Label: "2018/4/19", O: 2201.12, C: 2244.64, L: 2200.58, H: 2250.11, V: 249410},
		{Label: "2018/4/22", O: 2236.4, C: 2242.17, L: 2232.26, H: 2245.12, V: 149910},
		{Label: "2018/4/23", O: 2242.62, C: 2184.54, L: 2182.81, H: 2242.62, V: 119910},
		{Label: "2018/4/24", O: 2187.35, C: 2218.32, L: 2184.11, H: 2226.12, V: 239910},
		{Label: "2018/4/25", O: 2213.19, C: 2199.31, L: 2191.85, H: 2224.63, V: 214410},
		{Label: "2018/4/26", O: 2203.89, C: 2177.91, L: 2173.86, H: 2210.58, V: 234910},
		{Label: "2018/5/2", O: 2170.78, C: 2174.12, L: 2161.14, H: 2179.65, V: 299310},
		{Label: "2018/5/3", O: 2179.05, C: 2205.5, L: 2179.05, H: 2222.81, V: 149410},
		{Label: "2018/5/6", O: 2212.5, C: 2231.17, L: 2212.5, H: 2236.07, V: 299910},
		{Label: "2018/5/7", O: 2227.86, C: 2235.57, L: 2219.44, H: 2240.26, V: 139910},
		{Label: "2018/5/8", O: 2242.39, C: 2246.3, L: 2235.42, H: 2255.21, V: 119910},
		{Label: "2018/5/9", O: 2246.96, C: 2232.97, L: 2221.38, H: 2247.86, V: 164910},
		{Label: "2018/5/10", O: 2228.82, C: 2246.83, L: 2225.81, H: 2247.67, V: 96910},
		{Label: "2018/5/13", O: 2247.68, C: 2241.92, L: 2231.36, H: 2250.85, V: 149410},
		{Label: "2018/5/14", O: 2238.9, C: 2217.01, L: 2205.87, H: 2239.93, V: 119010},
		{Label: "2018/5/15", O: 2217.09, C: 2224.8, L: 2213.58, H: 2225.19, V: 139140},
		{Label: "2018/5/16", O: 2221.34, C: 2251.81, L: 2210.77, H: 2252.87, V: 194510},
		{Label: "2018/5/17", O: 2249.81, C: 2282.87, L: 2248.41, H: 2288.09, V: 189950},
		{Label: "2018/5/20", O: 2286.33, C: 2299.99, L: 2281.9, H: 2309.39, V: 289910},
		{Label: "2018/5/21", O: 2297.11, C: 2305.11, L: 2290.12, H: 2305.3, V: 271310},
		{Label: "2018/5/22", O: 2303.75, C: 2302.4, L: 2292.43, H: 2314.18, V: 249310},
		{Label: "2018/5/23", O: 2293.81, C: 2275.67, L: 2274.1, H: 2304.95, V: 239010},
		{Label: "2018/5/24", O: 2281.45, C: 2288.53, L: 2270.25, H: 2292.59, V: 240910},
		{Label: "2018/5/27", O: 2286.66, C: 2293.08, L: 2283.94, H: 2301.7, V: 109910},
		{Label: "2018/5/28", O: 2293.4, C: 2321.32, L: 2281.47, H: 2322.1, V: 149340},
		{Label: "2018/5/29", O: 2323.54, C: 2324.02, L: 2321.17, H: 2334.33, V: 239910},
		{Label: "2018/5/30", O: 2316.25, C: 2317.75, L: 2310.49, H: 2325.72, V: 143910},
		{Label: "2018/5/31", O: 2320.74, C: 2300.59, L: 2299.37, H: 2325.53, V: 109110},
		{Label: "2018/6/3", O: 2300.21, C: 2299.25, L: 2294.11, H: 2313.43, V: 147910},
		{Label: "2018/6/4", O: 2297.1, C: 2272.42, L: 2264.76, H: 2297.1, V: 179914},
		{Label: "2018/6/5", O: 2270.71, C: 2270.93, L: 2260.87, H: 2276.86, V: 194110},
		{Label: "2018/6/6", O: 2264.43, C: 2242.11, L: 2240.07, H: 2266.69, V: 149710},
		{Label: "2018/6/7", O: 2242.26, C: 2210.9, L: 2205.07, H: 2250.63, V: 279914},
		{Label: "2018/6/13", O: 2190.1, C: 2148.35, L: 2126.22, H: 2190.1, V: 239510},
	}

	events = []tachart.Event{
		{
			Type:        tachart.Long,
			Label:       cdls[55].Label,
			Description: "go long on " + cdls[55].Label,
		},
		{
			Type:        tachart.Open,
			Label:       cdls[60].Label,
			Description: "This is a demo event description. Randomly pick this candle to open position on " + cdls[60].Label,
		},
		{
			Type:        tachart.CustomEvent,
			Label:       cdls[65].Label,
			Description: "This is a user defined event demo, which custom event mark and color",
			EventMark: tachart.EventMark{
				Name:       "CM",
				FontColor:  "#000000",
				BgColor:    "#AAAAAA",
				SymbolSize: 32,
			},
		},
		{
			Type:        tachart.Short,
			Label:       cdls[71].Label,
			Description: "go short on " + cdls[71].Label,
		},
	}
)

func main() {
	top := `
<div style="border:2px solid blue;text-align:center;font-size:40px;height:80px;line-height:80px;">
Candlestick Chart Demo
</div>
`
	left := `
<div style="border:1px solid black;text-align:center;font-size:20px;width:60px;height:600px;">
left column
</div>
`
	right := `
<div style="border:1px solid black;text-align:center;font-size:20px;width:300px;height:600px;">
right column
</div>
`
	bottom := `
<div style="border:1px solid black;text-align:center;font-size:20px;margin:10px 0px;height:30px;line-height:30px;">
bottom bar
</div>
`
	vals0 := []float64{}
	for _, cdl := range cdls {
		vals0 = append(vals0, float64(int64(cdl.C)%100))
	}
	vals1 := []float64{}
	for _, cdl := range cdls {
		vals1 = append(vals1, float64(int64(cdl.C)%130))
	}

	cfg := tachart.NewConfig().
		SetTheme(tachart.ThemeVintage).
		SetChartWidth(900).
		SetChartHeight(800).
		SetTopRowContent(top, 100).
		SetBottomRowContent(bottom, 50).
		SetLeftColContent(left, 70).
		SetRightColContent(right, 300).
		SetDraggable(true).
		AddOverlay(
			//			tachart.NewSMA(5),
			//			tachart.NewSMA(20),
			tachart.NewBBandsSMA(20, 2),
		).
		AddIndicator(
			tachart.NewMACD(12, 26, 9),
			tachart.NewRSI(14, 30, 70),
			tachart.NewATR(5),
			tachart.NewBoundedLine("custom_bounded_line", vals0, 0, 100, 20, 80),
			tachart.NewLine2("double_line0", vals0, "double_line1", vals1),
			tachart.NewBar("bars", vals0),
		).
		UseRepoAssets() // serving assets file from current repo, avoid network access

	c := tachart.New(*cfg, cdls)
	c.GenStatic(cdls, events, "/Volumes/tmpfs/tmp/kline.html")
}
