package indicators

// SMA represents a Simple Moving Average indicator
type SMA struct {
	*BaseIndicator
	windowSize int
	valueSum   float64
}

// NewSMA creates a new Simple Moving Average indicator with specified window size
func NewSMA(windowSize int) *SMA {
	if windowSize <= 0 {
		panic("Window size must be greater than 0")
	}

	return &SMA{
		BaseIndicator: NewBaseIndicator("SMA"),
		windowSize:    windowSize,
		valueSum:      0.0,
	}
}

// AddValue adds a new value to the SMA calculation
func (sma *SMA) AddValue(value float64) {
	sma.input = append(sma.input, value)
	sma.valueSum += value

	// If we have enough values, calculate SMA
	if len(sma.input) >= sma.windowSize {
		// Calculate SMA
		average := sma.valueSum / float64(sma.windowSize)
		sma.AddOutput(average)

		// Remove oldest value from sum
		sma.valueSum -= sma.input[0]
		// Remove oldest value from input
		sma.input = sma.input[1:]
	}
}

// GetWindowSize returns the window size of the SMA
func (sma *SMA) GetWindowSize() int {
	return sma.windowSize
}
