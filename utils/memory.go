package utils

import (
	"fmt"
	"runtime"
	"time"

	"go.uber.org/zap"
)

// Get the current memory stats
func getMemoryStats() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Alloc
}

// Convert bytes to megabytes
func bytesToMB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024)
}

// Round the value to 2 decimal places
func roundTo2Decimals(val float64) float64 {
	return float64(int(val*100)) / 100
}

// Create log message
func logMemoryUsage(currentMemory uint64, diff int64, intervals int) {
	currentMemoryMB := roundTo2Decimals(bytesToMB(currentMemory))
	diffMB := roundTo2Decimals(bytesToMB(uint64(abs(diff))))

	sign := "+"
	if diff < 0 {
		sign = "-"
	}

	zap.L().Info(fmt.Sprintf("Memory Usage: Current %.3f MB, Diff %s%.3f MB. Running for %d intervals.", currentMemoryMB, sign, diffMB, intervals))
}

// Absolute value function for int64
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Start checking memory every minute
func MemoryChecker(interval time.Duration) {
	intervals := 0
	initialMemory := getMemoryStats()
	logMemoryUsage(initialMemory, 0, intervals)

	for {
		time.Sleep(interval)
		intervals++

		currentMemory := getMemoryStats()
		diff := int64(currentMemory) - int64(initialMemory)

		logMemoryUsage(currentMemory, diff, intervals)
	}
}
