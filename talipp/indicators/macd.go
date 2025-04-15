package indicators

// MACD represents Moving Average Convergence Divergence indicator
type MACD struct {
	*BaseIndicator
	fastEMA    *EMA
	slowEMA    *EMA
	signalEMA  *EMA
	macdValues []float64
	signalLine []float64
	histograms []float64
}

// MACDOutput represents the output of MACD calculations
type MACDOutput struct {
	MACD      float64
	Signal    float64
	Histogram float64
}

// NewMACD creates a new MACD indicator with specified parameters
// fastLength: the period for the fast EMA (default 12)
// slowLength: the period for the slow EMA (default 26)
// signalLength: the period for the signal line EMA (default 9)
func NewMACD(fastLength, slowLength, signalLength int) *MACD {
	if fastLength <= 0 || slowLength <= 0 || signalLength <= 0 {
		panic("All periods must be greater than 0")
	}

	if fastLength >= slowLength {
		panic("Fast length must be less than slow length")
	}

	return &MACD{
		BaseIndicator: NewBaseIndicator("MACD"),
		fastEMA:       NewEMA(fastLength),
		slowEMA:       NewEMA(slowLength),
		signalEMA:     NewEMA(signalLength),
		macdValues:    make([]float64, 0),
		signalLine:    make([]float64, 0),
		histograms:    make([]float64, 0),
	}
}

// AddValue adds a new value to the MACD calculation
func (macd *MACD) AddValue(value float64) {
	macd.input = append(macd.input, value)

	// Add value to both EMAs
	macd.fastEMA.AddValue(value)
	macd.slowEMA.AddValue(value)

	// If both EMAs have outputs, calculate MACD line
	if macd.fastEMA.IsInitialized() && macd.slowEMA.IsInitialized() {
		fastValue, _ := macd.fastEMA.GetLastValue()
		slowValue, _ := macd.slowEMA.GetLastValue()

		// MACD line = fast EMA - slow EMA
		macdValue := fastValue - slowValue
		macd.macdValues = append(macd.macdValues, macdValue)

		// Feed the MACD value into the signal EMA
		macd.signalEMA.AddValue(macdValue)

		// If the signal EMA has an output, calculate histogram
		if macd.signalEMA.IsInitialized() {
			signalValue, _ := macd.signalEMA.GetLastValue()
			macd.signalLine = append(macd.signalLine, signalValue)

			// Histogram = MACD line - signal line
			histogram := macdValue - signalValue
			macd.histograms = append(macd.histograms, histogram)

			// Store the output
			macd.AddOutput(macdValue)
			macd.AddOutput(signalValue)
			macd.AddOutput(histogram)

		}
	}
}

// GetMACDOutput returns the complete MACD output (MACD, Signal, Histogram)
func (macd *MACD) GetMACDOutput() []MACDOutput {
	results := make([]MACDOutput, 0)

	// Get the minimum length of the three slices
	minLength := len(macd.macdValues)
	if len(macd.signalLine) < minLength {
		minLength = len(macd.signalLine)
	}
	if len(macd.histograms) < minLength {
		minLength = len(macd.histograms)
	}

	// Build the output
	for i := 0; i < minLength; i++ {
		results = append(results, MACDOutput{
			MACD:      macd.macdValues[i],
			Signal:    macd.signalLine[i],
			Histogram: macd.histograms[i],
		})
	}

	return results
}

// GetMACDLine returns just the MACD line values
func (macd *MACD) GetMACDLine() []float64 {
	result := make([]float64, len(macd.macdValues))
	copy(result, macd.macdValues)
	return result
}

// GetSignalLine returns just the signal line values
func (macd *MACD) GetSignalLine() []float64 {
	result := make([]float64, len(macd.signalLine))
	copy(result, macd.signalLine)
	return result
}

// GetHistogram returns just the histogram values
func (macd *MACD) GetHistogram() []float64 {
	result := make([]float64, len(macd.histograms))
	copy(result, macd.histograms)
	return result
}
