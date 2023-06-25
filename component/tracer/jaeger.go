package tracer

import (
	"fmt"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/spf13/pflag"
	jg "go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

const defaultTracingRate = 1.0

type jaeger struct {
	id          string
	tracingRate float64
	agentURI    string
	port        int
	logger      appctx.Logger
}

func NewJaeger(id string) *jaeger {
	return &jaeger{
		id: id,
	}
}

func (j *jaeger) ID() string {
	return j.id
}

func (j *jaeger) InitFlags() {
	pflag.Float64Var(
		&j.tracingRate,
		"jaeger-tracing-rate",
		defaultTracingRate,
		"Sample tracing rate from OpenSensus: 0.0 -> 1.0 - Default: 1.0",
	)

	pflag.StringVar(
		&j.agentURI,
		"jaeger-agent-uri",
		"",
		"Jaeger agent URI to receive tracing data directly",
	)

	pflag.IntVar(
		&j.port,
		"jaeger-agent-port",
		6831,
		"Jaeger agent port",
	)
}

func (j *jaeger) isDisabled() bool {
	return j.agentURI == ""
}

func (j *jaeger) traceConfig() trace.Config {
	if j.tracingRate >= 1 {
		return trace.Config{
			DefaultSampler: trace.AlwaysSample(),
		}
	}

	return trace.Config{
		DefaultSampler: trace.ProbabilitySampler(j.tracingRate),
	}
}

func (j *jaeger) Run(ac appctx.AppContext) error {
	if j.isDisabled() {
		return nil
	}

	url := fmt.Sprintf("%s:%d", j.agentURI, j.port)

	je, err := jg.NewExporter(jg.Options{
		AgentEndpoint: url,
		Process: jg.Process{
			ServiceName: j.id,
		},
	})

	if err != nil {
		return err
	}

	trace.RegisterExporter(je)
	trace.ApplyConfig(j.traceConfig())

	j.logger = ac.Logger(j.ID())
	j.logger.Infof("Connect tracer (%s) on %s", j.ID(), url)

	return nil
}

func (j *jaeger) Stop() error {
	return nil
}
