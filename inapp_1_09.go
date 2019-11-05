// +build !go1.10

package raven

import (
	"strings"
)

// Determines whether frame should be marked as InApp
func isInAppFrame(frame StacktraceFrame, appPackagePrefixes []string) bool {
	if frame.Module == "main" {
		return true
	}
	for _, prefix := range appPackagePrefixes {
		if strings.HasPrefix(frame.Module, prefix) &&
			!strings.Contains(frame.Module, "vendor") &&
			!strings.Contains(frame.Module, "third_party") {
			return true
		}
	}
	return false
}
