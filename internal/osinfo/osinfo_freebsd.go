

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
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "unknown"
	}
	return strings.Split(string(out), "-")[0]
}
