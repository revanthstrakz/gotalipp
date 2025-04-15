package indicators

import (
	"math"
)

// RSI represents a Relative Strength Index indicator
type RSI struct {
	*BaseIndicator
	windowSize    int
	lastValue     float64
	avgGain       float64
	avgLoss       float64
	prevValue     float64
	firstValueSet bool
}

// NewRSI creates a new Relative Strength Index indicator with specified window size
func NewRSI(windowSize int) *RSI {
	if windowSize <= 0 {
		panic("Window size must be greater than 0")
	}

	return &RSI{
		BaseIndicator: NewBaseIndicator("RSI"),
		windowSize:    windowSize,
		lastValue:     0.0,
		avgGain:       0.0,
		avgLoss:       0.0,
		prevValue:     0.0,
		firstValueSet: false,
	}
}

// AddValue adds a new value to the RSI calculation
func (rsi *RSI) AddValue(value float64) {
	rsi.input = append(rsi.input, value)

	if !rsi.firstValueSet {
		rsi.prevValue = value
		rsi.firstValueSet = true
		return
	}

	// Calculate gain and loss
	change := value - rsi.prevValue
	var gain, loss float64
	if change > 0 {
		gain = change
		loss = 0.0
	} else {
		gain = 0.0
		loss = -change
	}

	if len(rsi.input) <= rsi.windowSize {
		// Accumulate initial values
		rsi.avgGain += gain
		rsi.avgLoss += loss

		// When we have enough data, calculate the first RSI
		if len(rsi.input) == rsi.windowSize {
			// Average initial gains and losses
			rsi.avgGain /= float64(rsi.windowSize)
			rsi.avgLoss /= float64(rsi.windowSize)

			// Calculate RSI
			if rsi.avgLoss == 0 {
				rsi.lastValue = 100.0
			} else {
				rs := rsi.avgGain / rsi.avgLoss
				rsi.lastValue = 100.0 - (100.0 / (1.0 + rs))
			}
			rsi.AddOutput(rsi.lastValue)
		}
	} else {
		// Use smoothed method for subsequent values
		rsi.avgGain = ((rsi.avgGain * float64(rsi.windowSize-1)) + gain) / float64(rsi.windowSize)
		rsi.avgLoss = ((rsi.avgLoss * float64(rsi.windowSize-1)) + loss) / float64(rsi.windowSize)

		// Calculate RSI
		if rsi.avgLoss == 0 {
			rsi.lastValue = 100.0
		} else {
			rs := rsi.avgGain / rsi.avgLoss
			rsi.lastValue = 100.0 - (100.0 / (1.0 + rs))
		}
		rsi.AddOutput(rsi.lastValue)
	}

	rsi.prevValue = value
}

// GetWindowSize returns the window size of the RSI
func (rsi *RSI) GetWindowSize() int {
	return rsi.windowSize
}
