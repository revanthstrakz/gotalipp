// Package indicators provides technical analysis indicators
package indicators

import (
	"fmt"
)

// Indicator is the interface that all indicators must implement
type Indicator interface {
	// AddValue adds a new value to the indicator
	AddValue(float64)
	// GetOutput returns the current output values of the indicator
	GetOutput() []float64
	// GetName returns the name of the indicator
	GetName() string
	// Reset clears all values in the indicator
	Reset()
}

// IndicatorWithWindow is an indicator with a window size
type IndicatorWithWindow interface {
	Indicator
	// GetWindowSize returns the window size of the indicator
	GetWindowSize() int
}

// BaseIndicator provides common functionality for all indicators
type BaseIndicator struct {
	name        string
	input       []float64
	output      []float64
	initialized bool
}

// NewBaseIndicator creates a new BaseIndicator
func NewBaseIndicator(name string) *BaseIndicator {
	return &BaseIndicator{
		name:        name,
		input:       make([]float64, 0),
		output:      make([]float64, 0),
		initialized: false,
	}
}

// GetName returns the name of the indicator
func (bi *BaseIndicator) GetName() string {
	return bi.name
}

// Reset clears all values in the indicator
func (bi *BaseIndicator) Reset() {
	bi.input = make([]float64, 0)
	bi.output = make([]float64, 0)
	bi.initialized = false
}

// IsInitialized returns whether the indicator has been initialized
func (bi *BaseIndicator) IsInitialized() bool {
	return bi.initialized
}

// GetOutput returns the current output values of the indicator
func (bi *BaseIndicator) GetOutput() []float64 {
	// Return a copy of the output slice to prevent modification
	result := make([]float64, len(bi.output))
	copy(result, bi.output)
	return result
}

// AddOutput adds a value to the output
func (bi *BaseIndicator) AddOutput(value float64) {
	bi.output = append(bi.output, value)
	bi.initialized = true
}

// GetValue returns the value at the specified index
func (bi *BaseIndicator) GetValue(index int) (float64, error) {
	if index < 0 || index >= len(bi.output) {
		return 0, fmt.Errorf("index out of range: %d", index)
	}
	return bi.output[index], nil
}

// GetLastValue returns the last value in the output
func (bi *BaseIndicator) GetLastValue() (float64, error) {
	if len(bi.output) == 0 {
		return 0, fmt.Errorf("no values in indicator %s", bi.name)
	}
	return bi.output[len(bi.output)-1], nil
}
