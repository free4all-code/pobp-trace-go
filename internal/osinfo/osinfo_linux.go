

package osinfo

import (
	"bufio"
	"os"
	"strings"
)

func osName() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return "Linux (Unknown Distribution)"
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	name := "Linux (Unknown Distribution)"
	for s.Scan() {
		parts := strings.SplitN(s.Text(), "=", 2)
		switch parts[0] {
		case "Name":
			name = strings.Trim(parts[1], "\"")
		}
	}
	return name
}

func osVersion() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return "unknown"
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	version := "unknown"
	for s.Scan() {
		parts := strings.SplitN(s.Text(), "=", 2)
		switch parts[0] {
		case "VERSION":
			version = strings.Trim(parts[1], "\"")
		case "VERSION_ID":
			if version == "" {
				version = strings.Trim(parts[1], "\"")
			}
		}
	}
	return version
}
