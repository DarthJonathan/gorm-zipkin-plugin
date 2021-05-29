package gormopentracing

import (
	"github.com/openzipkin/zipkin-go"
)

type options struct {
	// logResult means log SQL operation result into span log which causes span size grows up.
	// This is advised to only open in developing environment.
	logResult bool

	// tracer allows users to use customized and different tracer to makes tracing clearly.
	tracer zipkin.Tracer

	// Whether to log statement parameters or leave placeholders in the queries.
	logSqlParameters bool
}

func defaultOption(tracer zipkin.Tracer) *options {
	return &options{
		logResult:        false,
		tracer:           tracer,
		logSqlParameters: true,
	}
}

type applyOption func(o *options)

// WithLogResult enable opentracingPlugin to log the result of each executed sql.
func WithLogResult(logResult bool) applyOption {
	return func(o *options) {
		o.logResult = logResult
	}
}

// WithTracer allows to use customized tracer rather than the global one only.
func WithTracer(tracer *zipkin.Tracer) applyOption {
	return func(o *options) {
		if tracer == nil {
			return
		}

		o.tracer = *tracer
	}
}

func WithSqlParameters(logSqlParameters bool) applyOption {
	return func(o *options) {
		o.logSqlParameters = logSqlParameters
	}
}
