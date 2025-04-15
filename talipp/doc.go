/*
Package talipp is a Go implementation of technical analysis indicators
that offers similar functionality to the Python TaLipp package.

The package provides:
- Implementation of common technical indicators
- Support for streaming data input
- Immutable outputs with automatic caching
- Simple API for easy integration

Example usage:

    import (
        "talipp/indicators"
    )

    // Create a new Simple Moving Average with window size 3
    sma := indicators.NewSMA(3)

    // Add values
    sma.AddValue(10.0)
    sma.AddValue(11.0)
    sma.AddValue(12.0)

    // Get output values
    fmt.Println(sma.GetOutput()) // prints [11.0]

    // Add another value
    sma.AddValue(13.0)
    
    // Get updated output values
    fmt.Println(sma.GetOutput()) // prints [11.0, 12.0]
*/
package talipp
