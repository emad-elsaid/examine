package examine

import (
	"log/slog"

	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/debugger"
)

type GenericTracer struct {
	functionName string
}

func (s GenericTracer) function() string { return s.functionName }
func (s GenericTracer) trace(d *debugger.Debugger, state *api.DebuggerState, thread *api.Thread, bp *api.Breakpoint) any {
	slog.Info(s.function())

	return nil
}
