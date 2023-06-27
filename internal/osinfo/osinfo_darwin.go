

package osinfo

import (
	"os/exec"
	"runtime"
	"strings"
)

func osName() string {
	return runtime.GOOS
}

func osVersion() string {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return "unknown"
	}
	return strings.Trim(string(out), "\n")
}
