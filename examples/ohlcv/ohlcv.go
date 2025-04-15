// Package ohlcv provides types and functions for working with OHLCV (Open, High, Low, Close, Volume) data
package ohlcv

import (
	"time"
)

// OHLCV represents Open, High, Low, Close, Volume data
type OHLCV struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// NewOHLCV creates a new OHLCV object
func NewOHLCV(timestamp time.Time, open, high, low, close, volume float64) *OHLCV {
	return &OHLCV{
		Timestamp: timestamp,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    volume,
	}
}

// Stream represents a stream of OHLCV data
type Stream struct {
	Data    []*OHLCV
	onUpdate func(*OHLCV)
}

// NewStream creates a new OHLCV stream
func NewStream() *Stream {
	return &Stream{
		Data:     make([]*OHLCV, 0),
		onUpdate: nil,
	}
}

// Add adds a new OHLCV to the stream
func (s *Stream) Add(candle *OHLCV) {
	s.Data = append(s.Data, candle)
	if s.onUpdate != nil {
		s.onUpdate(candle)
	}
}

// OnUpdate sets a callback to be called when new data is added
func (s *Stream) OnUpdate(callback func(*OHLCV)) {
	s.onUpdate = callback
}

// Size returns the number of candles in the stream
func (s *Stream) Size() int {
	return len(s.Data)
}

// Get returns the candle at the specified index
func (s *Stream) Get(index int) (*OHLCV, error) {
	if index < 0 || index >= len(s.Data) {
		return nil, nil
	}
	return s.Data[index], nil
}

// GetFromLast returns the candle that is 'offset' positions from the end
// offset 0 returns the last candle, offset 1 returns the second last, etc.
func (s *Stream) GetFromLast(offset int) (*OHLCV, error) {
	index := len(s.Data) - 1 - offset
	return s.Get(index)
}

// Clear removes all candles from the stream
func (s *Stream) Clear() {
	s.Data = make([]*OHLCV, 0)
}

// High returns the high prices from the stream
func (s *Stream) High() []float64 {
	result := make([]float64, len(s.Data))
	for i, candle := range s.Data {
		result[i] = candle.High
	}
	return result
}

// Low returns the low prices from the stream
func (s *Stream) Low() []float64 {
	result := make([]float64, len(s.Data))
	for i, candle := range s.Data {
		result[i] = candle.Low
	}
	return result
}

// Close returns the close prices from the stream
func (s *Stream) Close() []float64 {
	result := make([]float64, len(s.Data))
	for i, candle := range s.Data {
		result[i] = candle.Close
	}
	return result
}

// Open returns the open prices from the stream
func (s *Stream) Open() []float64 {
	result := make([]float64, len(s.Data))
	for i, candle := range s.Data {
		result[i] = candle.Open
	}
	return result
}

// Volume returns the volumes from the stream
func (s *Stream) Volume() []float64 {
	result := make([]float64, len(s.Data))
	for i, candle := range s.Data {
		result[i] = candle.Volume
	}
	return result
}
