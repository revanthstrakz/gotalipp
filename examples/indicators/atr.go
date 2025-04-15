package indicators

import (
	"math"
)

// ATR represents an Average True Range indicator
type ATR struct {
	*BaseIndicator
	windowSize     int
	prevClose      float64
	firstValueSet  bool
	trueRanges     []float64
	smoothedValue  float64
}

// NewATR creates a new Average True Range indicator with specified window size
func NewATR(windowSize int) *ATR {
	if windowSize <= 0 {
		panic("Window size must be greater than 0")
	}

	return &ATR{
		BaseIndicator: NewBaseIndicator("ATR"),
		windowSize:    windowSize,
		prevClose:     0.0,
		firstValueSet: false,
		trueRanges:    make([]float64, 0),
		smoothedValue: 0.0,
	}
}

// AddOHLCValue adds a new OHLC candle data to the ATR calculation
func (atr *ATR) AddOHLCValue(high, low, close float64) {
	var trueRange float64
	
	if !atr.firstValueSet {
		// For the first value, true range is simply High - Low
		trueRange = high - low
		atr.prevClose = close
		atr.firstValueSet = true
	} else {
		// Calculate the true range
		tr1 := high - low                           // Current high - current low
		tr2 := math.Abs(high - atr.prevClose)       // Current high - previous close
		tr3 := math.Abs(low - atr.prevClose)        // Current low - previous close
		
		// True range is the maximum of the three
		trueRange = math.Max(tr1, math.Max(tr2, tr3))
		
		// Update previous close
		atr.prevClose = close
	}
	
	// Store the true range
	atr.trueRanges = append(atr.trueRanges, trueRange)
	
	// Calculate ATR
	if len(atr.trueRanges) == atr.windowSize {
		// For the first complete window, calculate a simple average
		var sum float64
		for _, tr := range atr.trueRanges {
			sum += tr
		}
		atr.smoothedValue = sum / float64(atr.windowSize)
		atr.AddOutput(atr.smoothedValue)
		
		// Remove the oldest true range
		atr.trueRanges = atr.trueRanges[1:]
	} else if len(atr.trueRanges) > atr.windowSize {
		// For subsequent values, use smoothed method
		atr.smoothedValue = ((atr.smoothedValue * float64(atr.windowSize-1)) + trueRange) / float64(atr.windowSize)
		atr.AddOutput(atr.smoothedValue)
		
		// Remove the oldest true range
		atr.trueRanges = atr.trueRanges[1:]
	}
}

// AddValue is not the preferred method for ATR, but can be used for compatibility
// with the Indicator interface. It will use the value as both high and low.
func (atr *ATR) AddValue(value float64) {
	atr.AddOHLCValue(value, value, value)
}

// GetWindowSize returns the window size of the ATR
func (atr *ATR) GetWindowSize() int {
	return atr.windowSize
}
