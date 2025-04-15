package main

import (
	"fmt"
	"math"
	"time"
)

// SMA (Simple Moving Average) implementation
type SMA struct {
	window int
	values []float64
	output []float64
}

func NewSMA(window int) *SMA {
	return &SMA{
		window: window,
		values: []float64{},
		output: []float64{},
	}
}

func (sma *SMA) AddValue(value float64) {
	sma.values = append(sma.values, value)
	
	// Calculate SMA when we have enough values
	if len(sma.values) >= sma.window {
		sum := 0.0
		for i := len(sma.values) - sma.window; i < len(sma.values); i++ {
			sum += sma.values[i]
		}
		sma.output = append(sma.output, sum/float64(sma.window))
	}
}

func (sma *SMA) GetOutput() []float64 {
	return sma.output
}

// EMA (Exponential Moving Average) implementation
type EMA struct {
	window     int
	alpha      float64
	lastValue  float64
	values     []float64
	output     []float64
	initialized bool
}

func NewEMA(window int) *EMA {
	return &EMA{
		window:     window,
		alpha:      2.0 / float64(window + 1),
		values:     []float64{},
		output:     []float64{},
		initialized: false,
	}
}

func (ema *EMA) AddValue(value float64) {
	ema.values = append(ema.values, value)
	
	// Initialize with SMA if we have enough values
	if !ema.initialized && len(ema.values) >= ema.window {
		// Calculate SMA for the first window
		sum := 0.0
		for i := 0; i < ema.window; i++ {
			sum += ema.values[i]
		}
		sma := sum / float64(ema.window)
		ema.output = append(ema.output, sma)
		ema.lastValue = sma
		ema.initialized = true
	} else if ema.initialized {
		// Calculate EMA: EMA_t = α × value + (1 - α) × EMA_t-1
		newEma := ema.alpha*value + (1-ema.alpha)*ema.lastValue
		ema.output = append(ema.output, newEma)
		ema.lastValue = newEma
	}
}

func (ema *EMA) GetOutput() []float64 {
	return ema.output
}

// RSI (Relative Strength Index) implementation
type RSI struct {
	window      int
	gains       []float64
	losses      []float64
	prevValue   float64
	avgGain     float64
	avgLoss     float64
	output      []float64
	initialized bool
}

func NewRSI(window int) *RSI {
	return &RSI{
		window:     window,
		gains:      []float64{},
		losses:     []float64{},
		output:     []float64{},
		initialized: false,
	}
}

func (rsi *RSI) AddValue(value float64) {
	// First value is just stored
	if !rsi.initialized && len(rsi.gains) == 0 {
		rsi.prevValue = value
		rsi.initialized = true
		return
	}
	
	// Calculate gain or loss
	change := value - rsi.prevValue
	gain := math.Max(0, change)
	loss := math.Max(0, -change)
	
	rsi.gains = append(rsi.gains, gain)
	rsi.losses = append(rsi.losses, loss)
	
	// Calculate RSI when we have enough values
	if len(rsi.gains) >= rsi.window {
		// For the first window, use simple average
		if len(rsi.gains) == rsi.window {
			sumGains := 0.0
			sumLosses := 0.0
			for i := 0; i < rsi.window; i++ {
				sumGains += rsi.gains[i]
				sumLosses += rsi.losses[i]
			}
			rsi.avgGain = sumGains / float64(rsi.window)
			rsi.avgLoss = sumLosses / float64(rsi.window)
		} else {
			// Use smoothed averages for subsequent values
			rsi.avgGain = (rsi.avgGain*float64(rsi.window-1) + gain) / float64(rsi.window)
			rsi.avgLoss = (rsi.avgLoss*float64(rsi.window-1) + loss) / float64(rsi.window)
		}
		
		// Calculate RSI
		var rs float64
		if rsi.avgLoss == 0 {
			rs = 100.0
		} else {
			rs = rsi.avgGain / rsi.avgLoss
		}
		newRSI := 100.0 - (100.0 / (1.0 + rs))
		rsi.output = append(rsi.output, newRSI)
	}
	
	rsi.prevValue = value
}

func (rsi *RSI) GetOutput() []float64 {
	return rsi.output
}

// MACD (Moving Average Convergence Divergence) implementation
type MACDOutput struct {
	MACD      float64
	Signal    float64
	Histogram float64
}

type MACD struct {
	fastEMA    *EMA
	slowEMA    *EMA
	signalEMA  *EMA
	values     []float64
	macdValues []float64
	output     []MACDOutput
}

func NewMACD(fastPeriod, slowPeriod, signalPeriod int) *MACD {
	return &MACD{
		fastEMA:    NewEMA(fastPeriod),
		slowEMA:    NewEMA(slowPeriod),
		signalEMA:  NewEMA(signalPeriod),
		values:     []float64{},
		macdValues: []float64{},
		output:     []MACDOutput{},
	}
}

func (macd *MACD) AddValue(value float64) {
	macd.values = append(macd.values, value)
	
	// Update both EMAs
	macd.fastEMA.AddValue(value)
	macd.slowEMA.AddValue(value)
	
	// Calculate MACD line if both EMAs have outputs
	fastOutput := macd.fastEMA.GetOutput()
	slowOutput := macd.slowEMA.GetOutput()
	
	if len(fastOutput) > 0 && len(slowOutput) > 0 {
		// MACD Line = Fast EMA - Slow EMA
		macdValue := fastOutput[len(fastOutput)-1] - slowOutput[len(slowOutput)-1]
		macd.macdValues = append(macd.macdValues, macdValue)
		
		// Update signal line (EMA of MACD line)
		macd.signalEMA.AddValue(macdValue)
		
		// Get signal value if available
		signalOutput := macd.signalEMA.GetOutput()
		if len(signalOutput) > 0 {
			signalValue := signalOutput[len(signalOutput)-1]
			
			// Histogram = MACD Line - Signal Line
			histValue := macdValue - signalValue
			
			// Add to output
			macd.output = append(macd.output, MACDOutput{
				MACD:      macdValue,
				Signal:    signalValue,
				Histogram: histValue,
			})
		}
	}
}

func (macd *MACD) GetMACDOutput() []MACDOutput {
	return macd.output
}

// BBands (Bollinger Bands) implementation
type BBandsOutput struct {
	Upper  float64
	Middle float64
	Lower  float64
}

type BBands struct {
	sma      *SMA
	window   int
	stdDevK  float64
	values   []float64
	output   []BBandsOutput
}

func NewBBands(window int, stdDevK float64) *BBands {
	return &BBands{
		sma:      NewSMA(window),
		window:   window,
		stdDevK:  stdDevK,
		values:   []float64{},
		output:   []BBandsOutput{},
	}
}

func (bb *BBands) AddValue(value float64) {
	bb.values = append(bb.values, value)
	bb.sma.AddValue(value)
	
	// Calculate bands if we have enough values
	if len(bb.values) >= bb.window {
		smaOutput := bb.sma.GetOutput()
		if len(smaOutput) > 0 {
			// Get the latest SMA value
			middle := smaOutput[len(smaOutput)-1]
			
			// Calculate standard deviation
			sumSquaredDeviations := 0.0
			startIdx := len(bb.values) - bb.window
			for i := startIdx; i < len(bb.values); i++ {
				deviation := bb.values[i] - middle
				sumSquaredDeviations += deviation * deviation
			}
			
			stdDev := math.Sqrt(sumSquaredDeviations / float64(bb.window))
			
			// Calculate upper and lower bands
			upper := middle + bb.stdDevK*stdDev
			lower := middle - bb.stdDevK*stdDev
			
			// Add to output
			bb.output = append(bb.output, BBandsOutput{
				Upper:  upper,
				Middle: middle,
				Lower:  lower,
			})
		}
	}
}

func (bb *BBands) GetBBandsOutput() []BBandsOutput {
	return bb.output
}

// OHLCV structure for candlestick data
type OHLCV struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

func NewOHLCV(timestamp time.Time, open, high, low, close, volume float64) *OHLCV {
	return &OHLCV{
		Timestamp: timestamp,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    volume,
	}
}

// Simple test function to show the indicators in action
func main() {
	// Example 1: Simple Moving Average
	fmt.Println("SMA Example:")
	sma := NewSMA(3)
	
	// Add some values
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	for _, value := range values {
		sma.AddValue(value)
		fmt.Printf("Added value: %.1f\n", value)
		
		// Print output only if we have any
		output := sma.GetOutput()
		if len(output) > 0 {
			fmt.Printf("Current SMA values: %v\n", output)
		}
	}
	fmt.Println()
	
	// Example 2: Relative Strength Index
	fmt.Println("RSI Example:")
	rsi := NewRSI(3)
	
	// Add some values
	values = []float64{10.0, 15.0, 12.0, 15.0, 17.0, 14.0}
	for _, value := range values {
		rsi.AddValue(value)
		fmt.Printf("Added value: %.1f\n", value)
		
		// Print output only if we have any
		output := rsi.GetOutput()
		if len(output) > 0 {
			fmt.Printf("Current RSI values: %v\n", output)
		}
	}
	fmt.Println()
	
	// Example 3: MACD
	fmt.Println("MACD Example:")
	macd := NewMACD(3, 5, 2)
	
	// Add some values (more needed for MACD to initialize)
	values = []float64{10.0, 11.0, 12.0, 13.0, 14.0, 13.0, 12.0, 11.0, 10.0}
	for _, value := range values {
		macd.AddValue(value)
	}
	
	// Get the MACD output
	macdOutput := macd.GetMACDOutput()
	if len(macdOutput) > 0 {
		fmt.Println("MACD outputs:")
		for i, output := range macdOutput {
			fmt.Printf("%d: MACD=%.4f, Signal=%.4f, Histogram=%.4f\n", 
				i, output.MACD, output.Signal, output.Histogram)
		}
	}
	fmt.Println()
	
	// Example 4: Bollinger Bands
	fmt.Println("Bollinger Bands Example:")
	bb := NewBBands(5, 2.0)
	
	// Add some values
	values = []float64{10.0, 12.0, 15.0, 14.0, 13.0, 16.0, 17.0, 18.0}
	for _, value := range values {
		bb.AddValue(value)
	}
	
	// Get the Bollinger Bands output
	bbOutput := bb.GetBBandsOutput()
	if len(bbOutput) > 0 {
		fmt.Println("Bollinger Bands outputs:")
		for i, output := range bbOutput {
			fmt.Printf("%d: Upper=%.4f, Middle=%.4f, Lower=%.4f\n", 
				i, output.Upper, output.Middle, output.Lower)
		}
	}
	
	// Example 5: Working with OHLCV data
	fmt.Println("\nOHLCV Stream Example:")
	
	// Create a simple dataset of candles
	now := time.Now()
	candles := []*OHLCV{
		NewOHLCV(now.Add(-5*time.Minute), 100.0, 105.0, 98.0, 103.0, 1000.0),
		NewOHLCV(now.Add(-4*time.Minute), 103.0, 107.0, 102.0, 105.0, 1200.0),
		NewOHLCV(now.Add(-3*time.Minute), 105.0, 108.0, 103.0, 107.0, 1100.0),
		NewOHLCV(now.Add(-2*time.Minute), 107.0, 110.0, 106.0, 109.0, 1300.0),
		NewOHLCV(now.Add(-1*time.Minute), 109.0, 112.0, 107.0, 111.0, 1400.0),
	}
	
	// Create indicators to process OHLCV data
	closeSMA := NewSMA(3)
	closeRSI := NewRSI(3)
	
	// Process each candle
	for _, candle := range candles {
		// Update indicators with close price
		closeSMA.AddValue(candle.Close)
		closeRSI.AddValue(candle.Close)
		
		// Print candle info
		fmt.Printf("Candle at %s: Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f\n", 
			candle.Timestamp.Format("15:04:05"), candle.Open, candle.High, candle.Low, candle.Close)
		
		// Print indicator values if available
		smaOutput := closeSMA.GetOutput()
		if len(smaOutput) > 0 {
			fmt.Printf("  SMA(3): %.4f\n", smaOutput[len(smaOutput)-1])
		}
		
		rsiOutput := closeRSI.GetOutput()
		if len(rsiOutput) > 0 {
			fmt.Printf("  RSI(3): %.4f\n", rsiOutput[len(rsiOutput)-1])
		}
	}
}