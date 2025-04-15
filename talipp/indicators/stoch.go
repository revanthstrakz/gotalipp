package indicators

import (
	"math"
)

// Stoch represents a Stochastic Oscillator indicator
type Stoch struct {
	*BaseIndicator
	windowSize      int
	smoothK         int
	smoothD         int
	highValues      []float64
	lowValues       []float64
	closeValues     []float64
	kValues         []float64
	dValues         []float64
}

// StochOutput represents the output of Stochastic Oscillator calculations
type StochOutput struct {
	K float64
	D float64
}

// NewStoch creates a new Stochastic Oscillator indicator
// windowSize: the period for the %K calculation (default 14)
// smoothK: the period for %K smoothing (default 1 - no smoothing)
// smoothD: the period for %D calculation (default 3)
func NewStoch(windowSize, smoothK, smoothD int) *Stoch {
	if windowSize <= 0 || smoothK <= 0 || smoothD <= 0 {
		panic("All periods must be greater than 0")
	}

	return &Stoch{
		BaseIndicator: NewBaseIndicator("Stoch"),
		windowSize:    windowSize,
		smoothK:       smoothK,
		smoothD:       smoothD,
		highValues:    make([]float64, 0),
		lowValues:     make([]float64, 0),
		closeValues:   make([]float64, 0),
		kValues:       make([]float64, 0),
		dValues:       make([]float64, 0),
	}
}

// AddValue is not the preferred method for Stochastic, but included for interface compatibility
func (stoch *Stoch) AddValue(value float64) {
	stoch.AddHLCValue(value, value, value)
}

// AddHLCValue adds a new high, low, close data to the Stochastic Oscillator calculation
func (stoch *Stoch) AddHLCValue(high, low, close float64) {
	stoch.highValues = append(stoch.highValues, high)
	stoch.lowValues = append(stoch.lowValues, low)
	stoch.closeValues = append(stoch.closeValues, close)
	
	// Keep only the windowSize values
	if len(stoch.highValues) > stoch.windowSize {
		stoch.highValues = stoch.highValues[1:]
		stoch.lowValues = stoch.lowValues[1:]
		stoch.closeValues = stoch.closeValues[1:]
	}
	
	// If we have enough values, calculate %K
	if len(stoch.closeValues) == stoch.windowSize {
		// Find highest high and lowest low in the window
		highestHigh := stoch.highValues[0]
		lowestLow := stoch.lowValues[0]
		
		for i := 1; i < stoch.windowSize; i++ {
			highestHigh = math.Max(highestHigh, stoch.highValues[i])
			lowestLow = math.Min(lowestLow, stoch.lowValues[i])
		}
		
		// Calculate raw %K
		var kValue float64
		if highestHigh == lowestLow {
			kValue = 50.0 // To avoid division by zero
		} else {
			kValue = 100.0 * ((close - lowestLow) / (highestHigh - lowestLow))
		}
		
		// Add K value for smoothing
		stoch.kValues = append(stoch.kValues, kValue)
		
		// Apply smoothing to %K if required
		if stoch.smoothK > 1 {
			if len(stoch.kValues) >= stoch.smoothK {
				var sum float64
				for i := len(stoch.kValues) - stoch.smoothK; i < len(stoch.kValues); i++ {
					sum += stoch.kValues[i]
				}
				kValue = sum / float64(stoch.smoothK)
				
				// Keep only the needed K values
				if len(stoch.kValues) > stoch.smoothK {
					stoch.kValues = stoch.kValues[1:]
				}
			} else {
				// Not enough data for K smoothing yet
				return
			}
		}
		
		// Calculate %D (SMA of %K)
		stoch.dValues = append(stoch.dValues, kValue)
		
		if len(stoch.dValues) >= stoch.smoothD {
			var sum float64
			for i := len(stoch.dValues) - stoch.smoothD; i < len(stoch.dValues); i++ {
				sum += stoch.dValues[i]
			}
			dValue := sum / float64(stoch.smoothD)
			
			// Keep only the needed D values
			if len(stoch.dValues) > stoch.smoothD {
				stoch.dValues = stoch.dValues[1:]
			}
			
			
			// Use K as the main indicator output
			stoch.AddOutput(kValue)
			stoch.AddOutput(dValue)
		}
	}
}

// GetWindowSize returns the window size of the Stochastic Oscillator
func (stoch *Stoch) GetWindowSize() int {
	return stoch.windowSize
}

// GetStochOutput returns the complete Stochastic Oscillator output (K, D)
func (stoch *Stoch) GetStochOutput() []StochOutput {
	// Find the minimum length between K and D values
	kLen := len(stoch.kValues)
	dLen := len(stoch.dValues)
	
	// Get the minimum of the two
	resultLen := kLen
	if dLen < resultLen {
		resultLen = dLen
	}
	
	// Start from the appropriate positions
	kStart := kLen - resultLen
	dStart := dLen - resultLen
	
	// Create result slice
	results := make([]StochOutput, resultLen)
	
	// Fill in the results
	for i := 0; i < resultLen; i++ {
		results[i] = StochOutput{
			K: stoch.kValues[kStart+i],
			D: stoch.dValues[dStart+i],
		}
	}
	
	return results
}

// GetKValues returns just the %K values
func (stoch *Stoch) GetKValues() []float64 {
	result := make([]float64, len(stoch.kValues))
	copy(result, stoch.kValues)
	return result
}

// GetDValues returns just the %D values
func (stoch *Stoch) GetDValues() []float64 {
	result := make([]float64, len(stoch.dValues))
	copy(result, stoch.dValues)
	return result
}
