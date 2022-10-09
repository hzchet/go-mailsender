package metrics

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func SetError(span trace.Span, err error) {
	span.SetStatus(codes.Error, err.Error())
}

func FollowSpan(ctx context.Context) (_ context.Context, span trace.Span) {
	ctx, span = tracer.Start(ctx, getCaller())
	return ctx, span
}

func getCaller() string {
	_, f, _, _ := runtime.Caller(2)
	return f
}

func FromCtx(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
