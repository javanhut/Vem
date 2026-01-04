package appcore

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	debugOnce sync.Once
	debugOn   bool
)

func debugEnabled() bool {
	debugOnce.Do(func() {
		value := strings.TrimSpace(strings.ToLower(os.Getenv("VEM_DEBUG")))
		debugOn = value == "1" || value == "true" || value == "yes" || value == "on"
	})
	return debugOn
}

func debugf(format string, args ...any) {
	if !debugEnabled() {
		return
	}
	if len(format) == 0 || format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Printf(format, args...)
}
