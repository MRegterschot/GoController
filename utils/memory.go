package utils

import (
	"fmt"
	"runtime"
	"time"

	"go.uber.org/zap"
)

// Function to get the current memory stats
func getMemoryStats() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Alloc // Returns the number of bytes allocated for heap objects
}

// Function to convert bytes to megabytes
func bytesToMB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024)
}

// Round the value to 2 decimal places
func roundTo2Decimals(val float64) float64 {
	return float64(int(val*100)) / 100
}

// Create log message
func logMemoryUsage(currentMemory uint64, diff int64, runtime int) {
	currentMemoryMB := roundTo2Decimals(bytesToMB(currentMemory))
	diffMB := roundTo2Decimals(bytesToMB(uint64(abs(diff))))

	sign := "+"
	if diff < 0 {
		sign = "-"
	}

	zap.L().Info(fmt.Sprintf("Memory Usage: Current %.3f MB, Diff %s%.3f MB. Running for %d min.", currentMemoryMB, sign, diffMB, runtime))
}

// Absolute value function for int64
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Function to start checking memory every minute
func MemoryChecker() {
	runtime := 0
	initialMemory := getMemoryStats()
	logMemoryUsage(initialMemory, 0, runtime)

	for {
		time.Sleep(1 * time.Minute)
		runtime++

		currentMemory := getMemoryStats()
		diff := int64(currentMemory) - int64(initialMemory) // Ensure signed calculation

		logMemoryUsage(currentMemory, diff, runtime)
	}
}
