

package tracer

import (
	"math"
	"os"
	"testing"

	"git.proto.group/protoobp/pobp-trace-go/internal/globalconfig"

	"github.com/stretchr/testify/assert"
)

func TestStartupLog(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		assert := assert.New(t)
		tp := new(testLogger)
		tracer, _, _, stop := startTestTracer(t, WithLogger(tp))
		defer stop()

		tp.Reset()
		logStartup(tracer)
		lines := removeAppSec(tp.Lines())
		assert.Len(lines, 2)
		assert.Regexp(`POBP Tracer v[0-9]+\.[0-9]+\.[0-9]+(-rc\.[0-9]+)? INFO: POBP-TRACE-GO TRACER CONFIGURATION {"date":"[^"]*","os_name":"[^"]*","os_version":"[^"]*","version":"[^"]*","lang":"Go","lang_version":"[^"]*","env":"","service":"tracer\.test","agent_url":"http://localhost:9/v0.4/traces","agent_error":"Post .*","debug":false,"analytics_enabled":false,"sample_rate":"NaN","sample_rate_limit":"disabled","sampling_rules":null,"sampling_rules_error":"","service_mappings":null,"tags":{"runtime-id":"[^"]*"},"runtime_metrics_enabled":false,"health_metrics_enabled":false,"profiler_code_hotspots_enabled":((false)|(true)),"profiler_endpoints_enabled":((false)|(true)),"pobp_version":"","architecture":"[^"]*","global_service":"","lambda_mode":"false","appsec":((true)|(false)),"agent_features":{"DropP0s":false,"Stats":false,"StatsdPort":0}}`, lines[1])
	})

	t.Run("configured", func(t *testing.T) {
		assert := assert.New(t)
		tp := new(testLogger)
		os.Setenv("POBP_TRACE_SAMPLE_RATE", "0.123")
		defer os.Unsetenv("POBP_TRACE_SAMPLE_RATE")
		tracer, _, _, stop := startTestTracer(t,
			WithLogger(tp),
			WithService("configured.service"),
			WithAgentAddr("test.host:1234"),
			WithEnv("configuredEnv"),
			WithServiceMapping("initial_service", "new_service"),
			WithGlobalTag("tag", "value"),
			WithGlobalTag("tag2", math.NaN()),
			WithRuntimeMetrics(),
			WithAnalyticsRate(1.0),
			WithServiceVersion("2.3.4"),
			WithSamplingRules([]SamplingRule{ServiceRule("mysql", 0.75)}),
			WithDebugMode(true),
		)
		defer globalconfig.SetAnalyticsRate(math.NaN())
		defer globalconfig.SetServiceName("")
		defer stop()

		tp.Reset()
		logStartup(tracer)
		assert.Len(tp.Lines(), 2)
		assert.Regexp(`POBP Tracer v[0-9]+\.[0-9]+\.[0-9]+(-rc\.[0-9]+)? INFO: POBP-TRACE-GO TRACER CONFIGURATION {"date":"[^"]*","os_name":"[^"]*","os_version":"[^"]*","version":"[^"]*","lang":"Go","lang_version":"[^"]*","env":"configuredEnv","service":"configured.service","agent_url":"http://localhost:9/v0.4/traces","agent_error":"Post .*","debug":true,"analytics_enabled":true,"sample_rate":"0\.123000","sample_rate_limit":"100","sampling_rules":\[{"service":"mysql","name":"","sample_rate":0\.75}\],"sampling_rules_error":"","service_mappings":{"initial_service":"new_service"},"tags":{"runtime-id":"[^"]*","tag":"value","tag2":"NaN"},"runtime_metrics_enabled":true,"health_metrics_enabled":true,"profiler_code_hotspots_enabled":((false)|(true)),"profiler_endpoints_enabled":((false)|(true)),"pobp_version":"2.3.4","architecture":"[^"]*","global_service":"configured.service","lambda_mode":"false","appsec":((true)|(false)),"agent_features":{"DropP0s":false,"Stats":false,"StatsdPort":0}}`, tp.Lines()[1])
	})

	t.Run("limit", func(t *testing.T) {
		assert := assert.New(t)
		tp := new(testLogger)
		os.Setenv("POBP_TRACE_SAMPLE_RATE", "0.123")
		defer os.Unsetenv("POBP_TRACE_SAMPLE_RATE")
		os.Setenv("POBP_TRACE_RATE_LIMIT", "1000.001")
		defer os.Unsetenv("POBP_TRACE_RATE_LIMIT")
		tracer, _, _, stop := startTestTracer(t,
			WithLogger(tp),
			WithService("configured.service"),
			WithAgentAddr("test.host:1234"),
			WithEnv("configuredEnv"),
			WithServiceMapping("initial_service", "new_service"),
			WithGlobalTag("tag", "value"),
			WithGlobalTag("tag2", math.NaN()),
			WithRuntimeMetrics(),
			WithAnalyticsRate(1.0),
			WithServiceVersion("2.3.4"),
			WithSamplingRules([]SamplingRule{ServiceRule("mysql", 0.75)}),
			WithDebugMode(true),
		)
		defer globalconfig.SetAnalyticsRate(math.NaN())
		defer globalconfig.SetServiceName("")
		defer stop()

		tp.Reset()
		logStartup(tracer)
		assert.Len(tp.Lines(), 2)
		assert.Regexp(`POBP Tracer v[0-9]+\.[0-9]+\.[0-9]+(-rc\.[0-9]+)? INFO: POBP-TRACE-GO TRACER CONFIGURATION {"date":"[^"]*","os_name":"[^"]*","os_version":"[^"]*","version":"[^"]*","lang":"Go","lang_version":"[^"]*","env":"configuredEnv","service":"configured.service","agent_url":"http://localhost:9/v0.4/traces","agent_error":"Post .*","debug":true,"analytics_enabled":true,"sample_rate":"0\.123000","sample_rate_limit":"1000.001","sampling_rules":\[{"service":"mysql","name":"","sample_rate":0\.75}\],"sampling_rules_error":"","service_mappings":{"initial_service":"new_service"},"tags":{"runtime-id":"[^"]*","tag":"value","tag2":"NaN"},"runtime_metrics_enabled":true,"health_metrics_enabled":true,"profiler_code_hotspots_enabled":((false)|(true)),"profiler_endpoints_enabled":((false)|(true)),"pobp_version":"2.3.4","architecture":"[^"]*","global_service":"configured.service","lambda_mode":"false","appsec":((true)|(false)),"agent_features":{"DropP0s":false,"Stats":false,"StatsdPort":0}}`, tp.Lines()[1])
	})

	t.Run("errors", func(t *testing.T) {
		assert := assert.New(t)
		tp := new(testLogger)
		os.Setenv("POBP_TRACE_SAMPLING_RULES", `[{"service": "some.service", "sample_rate": 0.234}, {"service": "other.service"}]`)
		defer os.Unsetenv("POBP_TRACE_SAMPLING_RULES")
		tracer, _, _, stop := startTestTracer(t, WithLogger(tp))
		defer stop()

		tp.Reset()
		logStartup(tracer)
		assert.Len(tp.Lines(), 2)
		assert.Regexp(`POBP Tracer v[0-9]+\.[0-9]+\.[0-9]+(-rc\.[0-9]+)? INFO: POBP-TRACE-GO TRACER CONFIGURATION {"date":"[^"]*","os_name":"[^"]*","os_version":"[^"]*","version":"[^"]*","lang":"Go","lang_version":"[^"]*","env":"","service":"tracer\.test","agent_url":"http://localhost:9/v0.4/traces","agent_error":"Post .*","debug":false,"analytics_enabled":false,"sample_rate":"NaN","sample_rate_limit":"100","sampling_rules":\[{"service":"some.service","name":"","sample_rate":0\.234}\],"sampling_rules_error":"found errors:\\n\\tat index 1: rate not provided","service_mappings":null,"tags":{"runtime-id":"[^"]*"},"runtime_metrics_enabled":false,"health_metrics_enabled":false,"profiler_code_hotspots_enabled":((false)|(true)),"profiler_endpoints_enabled":((false)|(true)),"pobp_version":"","architecture":"[^"]*","global_service":"","lambda_mode":"false","appsec":((true)|(false)),"agent_features":{"DropP0s":false,"Stats":false,"StatsdPort":0}}`, tp.Lines()[1])
	})

	t.Run("lambda", func(t *testing.T) {
		assert := assert.New(t)
		tp := new(testLogger)
		tracer, _, _, stop := startTestTracer(t, WithLogger(tp), WithLambdaMode(true))
		defer stop()

		tp.Reset()
		logStartup(tracer)
		assert.Len(tp.Lines(), 1)
		assert.Regexp(`POBP Tracer v[0-9]+\.[0-9]+\.[0-9]+(-rc\.[0-9]+)? INFO: POBP-TRACE-GO TRACER CONFIGURATION {"date":"[^"]*","os_name":"[^"]*","os_version":"[^"]*","version":"[^"]*","lang":"Go","lang_version":"[^"]*","env":"","service":"tracer\.test","agent_url":"http://localhost:9/v0.4/traces","agent_error":"","debug":false,"analytics_enabled":false,"sample_rate":"NaN","sample_rate_limit":"disabled","sampling_rules":null,"sampling_rules_error":"","service_mappings":null,"tags":{"runtime-id":"[^"]*"},"runtime_metrics_enabled":false,"health_metrics_enabled":false,"profiler_code_hotspots_enabled":((false)|(true)),"profiler_endpoints_enabled":((false)|(true)),"pobp_version":"","architecture":"[^"]*","global_service":"","lambda_mode":"true","appsec":((true)|(false)),"agent_features":{"DropP0s":false,"Stats":false,"StatsdPort":0}}`, tp.Lines()[0])
	})
}

func TestLogSamplingRules(t *testing.T) {
	assert := assert.New(t)
	tp := new(testLogger)
	os.Setenv("POBP_TRACE_SAMPLING_RULES", `[{"service": "some.service", "sample_rate": 0.234}, {"service": "other.service"}, {"service": "last.service", "sample_rate": 0.56}, {"odd": "pairs"}, {"sample_rate": 9.10}]`)
	defer os.Unsetenv("POBP_TRACE_SAMPLING_RULES")
	_, _, _, stop := startTestTracer(t, WithLogger(tp))
	defer stop()

	lines := removeAppSec(tp.Lines())
	assert.Len(lines, 2)
	assert.Contains(lines[0], "WARN: at index 4: ignoring rule {Service: Name: Rate:9.10}: rate is out of [0.0, 1.0] range")
	assert.Regexp(`POBP Tracer v[0-9]+\.[0-9]+\.[0-9]+(-rc\.[0-9]+)? WARN: DIAGNOSTICS Error\(s\) parsing POBP_TRACE_SAMPLING_RULES: found errors:\n\tat index 1: rate not provided\n\tat index 3: rate not provided$`, lines[1])
}

func TestLogAgentReachable(t *testing.T) {
	assert := assert.New(t)
	tp := new(testLogger)
	tracer, _, _, stop := startTestTracer(t, WithLogger(tp))
	defer stop()
	tp.Reset()
	logStartup(tracer)
	assert.Len(tp.Lines(), 2)
	assert.Regexp(`POBP Tracer v[0-9]+\.[0-9]+\.[0-9]+(-rc\.[0-9]+)? WARN: DIAGNOSTICS Unable to reach agent intake: Post`, tp.Lines()[0])
}

func TestLogFormat(t *testing.T) {
	assert := assert.New(t)
	tp := new(testLogger)
	tracer := newTracer(WithLogger(tp), WithRuntimeMetrics(), WithDebugMode(true))
	defer tracer.Stop()
	tp.Reset()
	tracer.StartSpan("test", ServiceName("test-service"), ResourceName("/"), WithSpanID(12345))
	assert.Len(tp.Lines(), 1)
	assert.Regexp(`POBP Tracer v[0-9]+\.[0-9]+\.[0-9]+(-rc\.[0-9]+)? DEBUG: Started Span: pobp.trace_id="12345" pobp.span_id="12345", Operation: test, Resource: /, Tags: map.*, map.*`, tp.Lines()[0])
}