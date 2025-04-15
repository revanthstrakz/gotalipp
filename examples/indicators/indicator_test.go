package indicators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMA(t *testing.T) {
	t.Run("Basic SMA calculation", func(t *testing.T) {
		sma := NewSMA(3)
		
		// Add values
		sma.AddValue(1.0)
		sma.AddValue(2.0)
		sma.AddValue(3.0)
		
		// Check output
		output := sma.GetOutput()
		assert.Equal(t, 1, len(output))
		assert.InDelta(t, 2.0, output[0], 0.0001)
		
		// Add another value
		sma.AddValue(4.0)
		
		// Check updated output
		output = sma.GetOutput()
		assert.Equal(t, 2, len(output))
		assert.InDelta(t, 2.0, output[0], 0.0001)
		assert.InDelta(t, 3.0, output[1], 0.0001)
	})
}

func TestEMA(t *testing.T) {
	t.Run("Basic EMA calculation", func(t *testing.T) {
		ema := NewEMA(3)
		
		// Add values
		ema.AddValue(1.0)
		ema.AddValue(2.0)
		ema.AddValue(3.0) // SMA should be calculated here (2.0)
		
		// Check output
		output := ema.GetOutput()
		assert.Equal(t, 1, len(output))
		assert.InDelta(t, 2.0, output[0], 0.0001)
		
		// Add another value - should use EMA formula
		// EMA = (Close - previousEMA) * multiplier + previousEMA
		// multiplier = 2 / (windowSize + 1) = 2 / 4 = 0.5
		// EMA = (4.0 - 2.0) * 0.5 + 2.0 = 3.0
		ema.AddValue(4.0)
		
		// Check updated output
		output = ema.GetOutput()
		assert.Equal(t, 2, len(output))
		assert.InDelta(t, 2.0, output[0], 0.0001)
		assert.InDelta(t, 3.0, output[1], 0.0001)
	})
}

func TestRSI(t *testing.T) {
	t.Run("Basic RSI calculation", func(t *testing.T) {
		rsi := NewRSI(2)
		
		// Add values that create equal gains and losses
		rsi.AddValue(10.0) // First value, just stored
		rsi.AddValue(15.0) // Gain: 5
		rsi.AddValue(10.0) // Loss: 5
		
		// Check output - with equal gains and losses, RSI should be 50
		output := rsi.GetOutput()
		assert.Equal(t, 1, len(output))
		assert.InDelta(t, 50.0, output[0], 0.0001)
		
		// Add a gain
		rsi.AddValue(15.0)
		
		// With more gains than losses, RSI should be higher than 50
		output = rsi.GetOutput()
		assert.Equal(t, 2, len(output))
		assert.Greater(t, output[1], 50.0)
	})
}

func TestMACD(t *testing.T) {
	t.Run("Basic MACD calculation", func(t *testing.T) {
		// Create MACD with small periods for testing
		macd := NewMACD(3, 5, 2)
		
		// Add enough values to get complete calculations
		for i := 1.0; i <= 10.0; i++ {
			macd.AddValue(i)
		}
		
		// Check we have some output
		output := macd.GetMACDOutput()
		assert.Greater(t, len(output), 0)
		
		// Verify that the MACD and Signal lines are different
		lastOutput := output[len(output)-1]
		assert.NotEqual(t, lastOutput.MACD, lastOutput.Signal)
		
		// Verify histogram calculation
		assert.InDelta(t, lastOutput.MACD-lastOutput.Signal, lastOutput.Histogram, 0.0001)
	})
}

func TestBBands(t *testing.T) {
	t.Run("Basic Bollinger Bands calculation", func(t *testing.T) {
		// Create Bollinger Bands with window size 3 and 2 standard deviations
		bb := NewBBands(3, 2.0)
		
		// Add constant values - should result in zero standard deviation
		bb.AddValue(10.0)
		bb.AddValue(10.0)
		bb.AddValue(10.0)
		
		// Get bands
		output := bb.GetBBandsOutput()
		assert.Equal(t, 1, len(output))
		
		// With all identical values, all bands should be equal
		assert.InDelta(t, 10.0, output[0].Middle, 0.0001)
		assert.InDelta(t, 10.0, output[0].Upper, 0.0001)
		assert.InDelta(t, 10.0, output[0].Lower, 0.0001)
		
		// Add a different value
		bb.AddValue(13.0)
		
		// Now the bands should diverge
		output = bb.GetBBandsOutput()
		assert.Equal(t, 2, len(output))
		assert.Greater(t, output[1].Upper, output[1].Middle)
		assert.Less(t, output[1].Lower, output[1].Middle)
	})
}

func TestATR(t *testing.T) {
	t.Run("Basic ATR calculation", func(t *testing.T) {
		// Create ATR with window size 3
		atr := NewATR(3)
		
		// Add OHLC values
		atr.AddOHLCValue(10.0, 8.0, 9.0)   // TR = 2.0 (high-low)
		atr.AddOHLCValue(11.0, 9.0, 10.0)  // TR = 2.0 (high-low)
		atr.AddOHLCValue(10.0, 7.0, 8.0)   // TR = 3.0 (high-low)
		
		// First ATR should be average of TRs: (2+2+3)/3 = 2.33
		output := atr.GetOutput()
		assert.Equal(t, 1, len(output))
		assert.InDelta(t, 2.33, output[0], 0.1)
	})
}

func TestStoch(t *testing.T) {
	t.Run("Basic Stochastic calculation", func(t *testing.T) {
		// Create Stochastic with window size 3, no %K smoothing, %D period 2
		stoch := NewStoch(3, 1, 2)
		
		// Add HLC values in an uptrend
		stoch.AddHLCValue(10.0, 8.0, 9.0)
		stoch.AddHLCValue(11.0, 9.0, 10.0)
		stoch.AddHLCValue(12.0, 10.0, 11.0)
		
		// First K should be 100 (close at highest point in range)
		// No D yet (need one more K value)
		output := stoch.GetKValues()
		assert.Equal(t, 3, len(output))
		assert.InDelta(t, 100.0, output[2], 0.1)
		
		// Add one more value
		stoch.AddHLCValue(12.0, 9.0, 10.0)
		
		// Now we should have D values
		stochOutput := stoch.GetStochOutput()
		assert.Greater(t, len(stochOutput), 0)
		
		// K value should be less than 100 now
		lastOutput := stochOutput[len(stochOutput)-1]
		assert.Less(t, lastOutput.K, 100.0)
	})
}
