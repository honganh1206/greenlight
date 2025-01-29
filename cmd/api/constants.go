package main

import "time"

const (
	minProcessingTime = 500 * time.Millisecond
	maxVariance       = 50 * time.Millisecond
	testIterations    = 5
)
