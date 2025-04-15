// Package talipp provides technical analysis indicators for Go
package talipp

import (
        _ "talipp/indicators" // Import indicators
        _ "talipp/ohlcv"      // Import OHLCV
)

// Version returns the version of the talipp package
func Version() string {
        return "0.1.0"
}
