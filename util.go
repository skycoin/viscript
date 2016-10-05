package main

import (
//"fmt"
)

// params: float value, negativemost, & positivemost bounds
func Clamp(f, negBoundary, posBoundary float32) float32 {
	if f < negBoundary {
		f = negBoundary
	}
	if f > posBoundary {
		f = posBoundary
	}

	return f
}
