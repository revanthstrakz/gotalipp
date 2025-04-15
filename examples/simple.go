package main

import (
        "fmt"

        "talipp/talipp/indicators"  // Using local import path
)

func main() {
        // Example 1: Simple Moving Average
        fmt.Println("SMA Example:")
        sma := indicators.NewSMA(3)
        
        // Add some values
        values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
        for _, value := range values {
                sma.AddValue(value)
                fmt.Printf("Added value: %.1f\n", value)
                
                // Print output only if we have any
                output := sma.GetOutput()
                if len(output) > 0 {
                        fmt.Printf("Current SMA values: %v\n", output)
                }
        }
        fmt.Println()
        
        // Example 2: Relative Strength Index
        fmt.Println("RSI Example:")
        rsi := indicators.NewRSI(3)
        
        // Add some values
        values = []float64{10.0, 15.0, 12.0, 15.0, 17.0, 14.0}
        for _, value := range values {
                rsi.AddValue(value)
                fmt.Printf("Added value: %.1f\n", value)
                
                // Print output only if we have any
                output := rsi.GetOutput()
                if len(output) > 0 {
                        fmt.Printf("Current RSI values: %v\n", output)
                }
        }
        fmt.Println()
        
        // Example 3: MACD
        fmt.Println("MACD Example:")
        macd := indicators.NewMACD(3, 5, 2)
        
        // Add some values (more needed for MACD to initialize)
        values = []float64{10.0, 11.0, 12.0, 13.0, 14.0, 13.0, 12.0, 11.0, 10.0}
        for _, value := range values {
                macd.AddValue(value)
        }
        
        // Get the MACD output
        macdOutput := macd.GetMACDOutput()
        if len(macdOutput) > 0 {
                fmt.Println("MACD outputs:")
                for i, output := range macdOutput {
                        fmt.Printf("%d: MACD=%.4f, Signal=%.4f, Histogram=%.4f\n", 
                                i, output.MACD, output.Signal, output.Histogram)
                }
        }
        fmt.Println()
        
        // Example 4: Bollinger Bands
        fmt.Println("Bollinger Bands Example:")
        bb := indicators.NewBBands(5, 2.0)
        
        // Add some values
        values = []float64{10.0, 12.0, 15.0, 14.0, 13.0, 16.0, 17.0, 18.0}
        for _, value := range values {
                bb.AddValue(value)
        }
        
        // Get the Bollinger Bands output
        bbOutput := bb.GetBBandsOutput()
        if len(bbOutput) > 0 {
                fmt.Println("Bollinger Bands outputs:")
                for i, output := range bbOutput {
                        fmt.Printf("%d: Upper=%.4f, Middle=%.4f, Lower=%.4f\n", 
                                i, output.Upper, output.Middle, output.Lower)
                }
        }
}
