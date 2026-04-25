package observability

import (
	"context"
	"errors"
	"log"

	"github.com/ak-repo/wim/internal/errs"
)

func Report(ctx context.Context, err error, stack []string) {
	if err == nil {
		return
	}

	var e *errs.Error
	if !errors.As(err, &e) {
		return
	}

	log.Printf("[ERROR] kind=%s code=%s op_stack=%v error=%v",
		e.Kind.String(),
		e.Code,
		stack,
		err,
	)

	// TODO: Sentry integration
	// if sentryHub := sentry.GetHubFromContext(ctx); sentryHub != nil {
	//     sentryHub.CaptureException(err)
	//     sentryHub.ConfigureScope(func(scope *sentry.Scope) {
	//         scope.SetTag("kind", toKindString(e.Kind))
	//         scope.SetTag("code", e.Code)
	//         scope.SetExtra("op_stack", stack)
	//     })
	// }

	// TODO: OpenTelemetry integration
	// span := trace.SpanFromContext(ctx)
	// if span.SpanContext().IsValid() {
	//     span.RecordError(err)
	//     span.SetAttributes(
	//         attribute.String("error.op_stack", strings.Join(stack, " > ")),
	//         attribute.String("error.kind", toKindString(e.Kind)),
	//         attribute.String("error.code", e.Code),
	//     )
	// }
}
