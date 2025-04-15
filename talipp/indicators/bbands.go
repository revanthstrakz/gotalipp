package indicators

import (
	"math"
)

// BBands represents Bollinger Bands indicator
type BBands struct {
	*BaseIndicator
	windowSize      int
	deviationFactor float64
	sma             *SMA
	prices          []float64
	upperBands      []float64
	middleBands     []float64
	lowerBands      []float64
}

// BBandsOutput represents the output of Bollinger Bands calculations
type BBandsOutput struct {
	Upper  float64
	Middle float64
	Lower  float64
}

// NewBBands creates a new Bollinger Bands indicator
// windowSize: the period for the SMA (default 20)
// deviationFactor: the standard deviation factor (default 2.0)
func NewBBands(windowSize int, deviationFactor float64) *BBands {
	if windowSize <= 0 {
		panic("Window size must be greater than 0")
	}

	return &BBands{
		BaseIndicator:   NewBaseIndicator("BBands"),
		windowSize:      windowSize,
		deviationFactor: deviationFactor,
		sma:             NewSMA(windowSize),
		prices:          make([]float64, 0),
		upperBands:      make([]float64, 0),
		middleBands:     make([]float64, 0),
		lowerBands:      make([]float64, 0),
	}
}

// AddValue adds a new value to the Bollinger Bands calculation
func (bb *BBands) AddValue(value float64) {
	bb.input = append(bb.input, value)
	bb.prices = append(bb.prices, value)
	
	// Keep only the last windowSize values
	if len(bb.prices) > bb.windowSize {
		bb.prices = bb.prices[1:]
	}
	
	// Add value to SMA
	bb.sma.AddValue(value)
	
	// If we have enough values, calculate bands
	if len(bb.prices) == bb.windowSize && bb.sma.IsInitialized() {
		// Get latest SMA value
		smaValue, _ := bb.sma.GetLastValue()
		
		// Calculate standard deviation
		var sum float64
		for _, price := range bb.prices {
			sum += math.Pow(price-smaValue, 2)
		}
		stdDev := math.Sqrt(sum / float64(bb.windowSize))
		
		// Calculate bands
		upperBand := smaValue + (bb.deviationFactor * stdDev)
		lowerBand := smaValue - (bb.deviationFactor * stdDev)
		
		// Store band values
		bb.upperBands = append(bb.upperBands, upperBand)
		bb.middleBands = append(bb.middleBands, smaValue)
		bb.lowerBands = append(bb.lowerBands, lowerBand)
		
		// Add output - using middle band as the output value for the base indicator
		bb.AddOutput(smaValue)
	}
}

// GetWindowSize returns the window size of the Bollinger Bands
func (bb *BBands) GetWindowSize() int {
	return bb.windowSize
}

// GetBBandsOutput returns the complete Bollinger Bands output (Upper, Middle, Lower)
func (bb *BBands) GetBBandsOutput() []BBandsOutput {
	results := make([]BBandsOutput, 0)
	
	// Get the minimum length of the three slices
	minLength := len(bb.upperBands)
	if len(bb.middleBands) < minLength {
		minLength = len(bb.middleBands)
	}
	if len(bb.lowerBands) < minLength {
		minLength = len(bb.lowerBands)
	}
	
	// Build the output
	for i := 0; i < minLength; i++ {
		results = append(results, BBandsOutput{
			Upper:  bb.upperBands[i],
			Middle: bb.middleBands[i],
			Lower:  bb.lowerBands[i],
		})
	}
	
	return results
}

// GetUpperBand returns just the upper band values
func (bb *BBands) GetUpperBand() []float64 {
	result := make([]float64, len(bb.upperBands))
	copy(result, bb.upperBands)
	return result
}

// GetMiddleBand returns just the middle band values
func (bb *BBands) GetMiddleBand() []float64 {
	result := make([]float64, len(bb.middleBands))
	copy(result, bb.middleBands)
	return result
}

// GetLowerBand returns just the lower band values
func (bb *BBands) GetLowerBand() []float64 {
	result := make([]float64, len(bb.lowerBands))
	copy(result, bb.lowerBands)
	return result
}
