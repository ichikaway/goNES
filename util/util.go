package util

import (
	"fmt"
	"os"
	"time"
)

func Bool2Uint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}


func GetFps(count int, startTime time.Time) int {
	end := time.Now()
	sec := int(end.Sub(startTime).Seconds())
	if sec == 0 {
		sec = 1
	}
	return count / sec
}

func Dec2Hex(i interface{}) interface{}{
	return fmt.Sprintf("%02X", i)
}

func DebugMode() bool {
	if len(os.Args) > 2 {
		debug := os.Args[2]
		if debug == "debug" {
			return true
		}
	}
	return false
}