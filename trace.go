package examine

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/debugger"
)

type Timeline struct {
	mutex       sync.RWMutex
	Tracepoints []*TracePoint
}

func NewTimeline() *Timeline {
	return &Timeline{
		Tracepoints: []*TracePoint{},
	}
}

func (t *Timeline) Add(p *TracePoint) {
	t.mutex.Lock()
	t.Tracepoints = append(t.Tracepoints, p)
	t.mutex.Unlock()
}

func (t *Timeline) JSON() ([]byte, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return json.Marshal(t.Tracepoints)
}

// All string fields will be repeated, so we can try to make them more effecient
// in storing and transfer by replacing them by a number maybe
type TracePoint struct {
	Function *string
	File     *string
	Line     int
	Info     any
	Time     int64
}

type Tracer interface {
	function() string
	trace(*debugger.Debugger, *api.DebuggerState, *api.Thread, *api.Breakpoint) any
}

var continueCmd = api.DebuggerCommand{
	Name:                 api.Continue,
	ReturnInfoLoadConfig: &loadArgs,
}

func trace(dbg *debugger.Debugger, timeline *Timeline, tracers ...Tracer) {
	for _, t := range tracers {
		slog.Info("Breakpoint", "function", t.function())
		_, err := dbg.CreateBreakpoint(&api.Breakpoint{
			Name:         t.function(),
			FunctionName: t.function(),
			LoadArgs:     &loadArgs,
			LoadLocals:   &loadArgs,
			Stacktrace:   0,
			Goroutine:    false,
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

		for _, thread := range state.Threads {
			bp := thread.Breakpoint
			if bp == nil {
				if thread == state.CurrentThread {
					slog.Error("Current thread doesn't have a breakpoint")
				}
				continue
			}

			for _, t := range tracers {
				if t.function() == bp.FunctionName {
					timeline.Add(&TracePoint{
						Function: &bp.FunctionName,
						File:     &bp.File,
						Line:     bp.Line,
						Info:     t.trace(dbg, state, thread, bp),
						Time:     time.Now().UnixNano(),
					})
				}
			}
		}
	}
}
