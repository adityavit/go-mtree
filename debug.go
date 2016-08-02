package mtree

import (
	"fmt"
	"os"
	"time"
)

// DebugOutput is the where DEBUG output is written
var DebugOutput = os.Stderr

// Debugln writes output to DebugOutput, only if DEBUG environment variable is set
func Debugln(a ...interface{}) (n int, err error) {
	if os.Getenv("DEBUG") != "" {
		return fmt.Fprintf(DebugOutput, "[%d] [DEBUG] %s\n", time.Now().UnixNano(), a)
	}
	return 0, nil
}

// Debugf writes formatted output to DebugOutput, only if DEBUG environment variable is set
func Debugf(format string, a ...interface{}) (n int, err error) {
	if os.Getenv("DEBUG") != "" {
		return fmt.Fprintf(DebugOutput, "[%d] [DEBUG] %s\n", time.Now().UnixNano(), fmt.Sprintf(format, a...))
	}
	return 0, nil
}
