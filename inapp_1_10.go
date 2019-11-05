// +build go1.10

package raven

import (
	"go/build"
	"os"
	"strings"
)

// As of go 1.10, the value of build.Default.GOROOT in a running program is
//  either $GOROOT of the executing machine, if set, or $GOROOT/$GOROOT_FINAL of
// the building machine (set by the linker).  We want the latter, so we can
// reliably identify stack frames from go stdlib packages.
//
// See https://go-review.googlesource.com/c/go/+/86835/ for more details.
//
var gorootFromLinker bool

func init() {
	gorootFromLinker = os.Getenv("GOROOT") == ""
}

// Determines whether frame should be marked as InApp
func isInAppFrame(frame StacktraceFrame, appPackagePrefixes []string) bool {
	if gorootFromLinker {
		if strings.HasPrefix(frame.AbsolutePath, build.Default.GOROOT) ||
			strings.Contains(frame.AbsolutePath, "/go/pkg/mod/") ||
			strings.Contains(frame.Module, "vendor") ||
			strings.Contains(frame.Module, "third_party") {
			return false
		}
		return true
	} else {
		if frame.Module == "main" {
			return true
		}
		for _, prefix := range appPackagePrefixes {
			if strings.HasPrefix(frame.Module, prefix) &&
				!strings.Contains(frame.AbsolutePath, "/go/pkg/mod/") &&
				!strings.Contains(frame.Module, "vendor") &&
				!strings.Contains(frame.Module, "third_party") {
				return true
			}
		}
		return false
	}
}
