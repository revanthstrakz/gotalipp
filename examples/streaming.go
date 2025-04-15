package examples

import (
	"fmt"
	"time"

	"github.com/revanthstrakz/gotalipp/talipp/indicators"
	"github.com/revanthstrakz/gotalipp/talipp/ohlcv"
)

func Streaming_example() {
	// Create an OHLCV stream
	stream := ohlcv.NewStream()

	// Create indicators
	sma := indicators.NewSMA(3)
	rsi := indicators.NewRSI(5)
	atr := indicators.NewATR(5)

	// Set up stream update callback
	stream.OnUpdate(func(candle *ohlcv.OHLCV) {
		// Update SMA with close price
		sma.AddValue(candle.Close)

		// Update RSI with close price
		rsi.AddValue(candle.Close)

		// Update ATR with high, low, close
		atr.AddOHLCValue(candle.High, candle.Low, candle.Close)

		// Print current indicator values
		fmt.Printf("New candle at %s: Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f\n",
			candle.Timestamp.Format("15:04:05"), candle.Open, candle.High, candle.Low, candle.Close)

		// Print indicator outputs if available
		smaOutput := sma.GetOutput()
		if len(smaOutput) > 0 {
			fmt.Printf("  SMA(3): %.4f\n", smaOutput[len(smaOutput)-1])
		}

		rsiOutput := rsi.GetOutput()
		if len(rsiOutput) > 0 {
			fmt.Printf("  RSI(5): %.4f\n", rsiOutput[len(rsiOutput)-1])
		}

		atrOutput := atr.GetOutput()
		if len(atrOutput) > 0 {
			fmt.Printf("  ATR(5): %.4f\n", atrOutput[len(atrOutput)-1])
		}
	})

	// Simulate streaming data
	fmt.Println("Simulating streaming data...")

	// Create some sample data
	now := time.Now()
	candles := []*ohlcv.OHLCV{
		ohlcv.NewOHLCV(now.Add(-7*time.Minute), 100.0, 105.0, 98.0, 103.0, 1000.0),
		ohlcv.NewOHLCV(now.Add(-6*time.Minute), 103.0, 107.0, 102.0, 105.0, 1200.0),
		ohlcv.NewOHLCV(now.Add(-5*time.Minute), 105.0, 108.0, 103.0, 107.0, 1100.0),
		ohlcv.NewOHLCV(now.Add(-4*time.Minute), 107.0, 110.0, 106.0, 108.0, 1300.0),
		ohlcv.NewOHLCV(now.Add(-3*time.Minute), 108.0, 112.0, 107.0, 110.0, 1400.0),
		ohlcv.NewOHLCV(now.Add(-2*time.Minute), 110.0, 115.0, 109.0, 112.0, 1500.0),
		ohlcv.NewOHLCV(now.Add(-1*time.Minute), 112.0, 116.0, 111.0, 115.0, 1600.0),
		ohlcv.NewOHLCV(now, 115.0, 120.0, 114.0, 118.0, 1700.0),
	}

	// Add candles one by one with a delay to simulate streaming
	for _, candle := range candles {
		stream.Add(candle)
		time.Sleep(500 * time.Millisecond)
	}

	// Show final stats
	fmt.Println("\nFinal indicator values:")
	fmt.Printf("SMA(3): %v\n", sma.GetOutput())
	fmt.Printf("RSI(5): %v\n", rsi.GetOutput())
	fmt.Printf("ATR(5): %v\n", atr.GetOutput())

	// Extract just the close prices from the stream
	fmt.Println("\nClose prices from stream:")
	fmt.Printf("%v\n", stream.Close())
}
