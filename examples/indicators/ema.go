package indicators

// EMA represents an Exponential Moving Average indicator
type EMA struct {
	*BaseIndicator
	windowSize int
	alpha      float64
	lastValue  float64
}

// NewEMA creates a new Exponential Moving Average indicator with specified window size
func NewEMA(windowSize int) *EMA {
	if windowSize <= 0 {
		panic("Window size must be greater than 0")
	}

	return &EMA{
		BaseIndicator: NewBaseIndicator("EMA"),
		windowSize:    windowSize,
		alpha:         2.0 / float64(windowSize+1),
		lastValue:     0.0,
	}
}

// AddValue adds a new value to the EMA calculation
func (ema *EMA) AddValue(value float64) {
	ema.input = append(ema.input, value)

	if !ema.IsInitialized() {
		// For the first windowSize values, we'll use SMA
		if len(ema.input) == ema.windowSize {
			// Calculate the initial SMA
			var sum float64
			for _, val := range ema.input {
				sum += val
			}
			ema.lastValue = sum / float64(ema.windowSize)
			ema.AddOutput(ema.lastValue)
		}
	} else {
		// Use EMA formula: EMA = (Close - previousEMA) * multiplier + previousEMA
		ema.lastValue = (value-ema.lastValue)*ema.alpha + ema.lastValue
		ema.AddOutput(ema.lastValue)
	}
}

// GetWindowSize returns the window size of the EMA
func (ema *EMA) GetWindowSize() int {
	return ema.windowSize
}
