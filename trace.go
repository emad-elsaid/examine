package examine

import (
	"log/slog"
	"os"

	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/debugger"
)

type Tracer interface {
	function() string
	trace(*debugger.Debugger, *api.DebuggerState)
}

var continueCmd = api.DebuggerCommand{
	Name:                 api.Continue,
	ReturnInfoLoadConfig: &loadArgs,
}

func trace(dbg *debugger.Debugger, tracers []Tracer) {
	for _, t := range tracers {
		_, err := dbg.CreateBreakpoint(&api.Breakpoint{
			Name:         t.function(),
			FunctionName: t.function(),
			LoadArgs:     &loadArgs,
			LoadLocals:   &loadArgs,
			Stacktrace:   1,
			Goroutine:    true,
		}, "", nil, false)
		if err != nil {
			slog.Error("Error creating breakpoint", "error", err)
		}
	}

	for {
		resumed := make(chan struct{})
		state, err := dbg.Command(&continueCmd, resumed)
		if err != nil {
			slog.Error(err.Error())
		}

		// Wait until it's resumed
		<-resumed

		if state.Exited {
			os.Exit(0)
		}

		thread := state.CurrentThread
		bp := thread.Breakpoint
		if bp == nil {
			slog.Error("Current thread doesn't have a breakpoint")
			continue
		}

		for _, t := range tracers {
			if t.function() == bp.FunctionName {
				t.trace(dbg, state)
			}
		}

	}
}
