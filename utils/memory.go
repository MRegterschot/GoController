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
func logMemoryUsage(currentMemoryMB uint64, diffMB uint64) {
	zap.L().Info(fmt.Sprintf("Memory Usage: Current %.3f MB, Diff %.3f MB", bytesToMB(currentMemoryMB), bytesToMB(diffMB)))
}

// Function to start checking memory every minute
func MemoryChecker() {
	initialMemory := getMemoryStats()
	logMemoryUsage(initialMemory, 0)

	for {
		time.Sleep(1 * time.Minute)

		currentMemory := getMemoryStats()
		diff := currentMemory - initialMemory

		logMemoryUsage(currentMemory, diff)

		initialMemory = currentMemory
	}
}
