

//go:build !appsec
// +build !appsec

package appsec

import "git.proto.group/protoobp/pobp-trace-go/internal/log"

// Enabled returns true when AppSec is up and running. Meaning that the appsec build tag is enabled, the env var
// POBP_APPSEC_ENABLED is set to true, and the tracer is started.
func Enabled() bool {
	return false
}

// Start AppSec when enabled is enabled by both using the appsec build tag and
// setting the environment variable POBP_APPSEC_ENABLED to true.
func Start() {
	if enabled, err := isEnabled(); err != nil {
		// Something went wrong while checking the POBP_APPSEC_ENABLED configuration
		log.Error("appsec: error while checking if appsec is enabled: %v", err)
	} else if enabled {
		// The user is willing to enabled appsec but didn't have the build tag
		log.Info("appsec: enabled by the configuration but has not been activated during the compilation: please add the go build tag `appsec` to your build options to enable it")
	} else {
		// The user is not willing to start appsec, a simple debug log is enough
		log.Debug("appsec: not been not enabled during the compilation: please add the go build tag `appsec` to your build options to enable it")
	}
}

// Stop AppSec.
func Stop() {}

// Static rule stubs when disabled.
const staticRecommendedRule = ""
