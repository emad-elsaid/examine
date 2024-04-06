package examine

import (
	"log/slog"

	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/debugger"
)

var httpTracers = []Tracer{
	ServeHTTP{},
}

type ServeHTTP struct{}

func (s ServeHTTP) function() string { return "net/http.(*ServeMux).ServeHTTP" }
func (s ServeHTTP) trace(d *debugger.Debugger, state *api.DebuggerState, thread *api.Thread, bp *api.Breakpoint) {
	slog.Info(s.function())

	info := thread.BreakpointInfo
	if bp == nil {
		slog.Error("BreakpointInfo is nil", "tracer", s.function())
		return
	}

	var args Variables = info.Arguments
	if args == nil {
		slog.Error("Arguments is nil", "tracer", s.function())
		return
	}

	r := args.Name("r")
	if r == nil {
		slog.Error("Can't get argument: r")
		return
	}

	r = r.Dereference()

	method := r.Field("Method")
	if method.Error() != "" {
		slog.Error(method.Error())
	}

	url := r.Field("URL").Dereference()

	path := url.Field("Path")
	if path.Error() != "" {
		slog.Error(path.Error())
	}

	slog.Info(s.function(), "method", method.String(), "path", path.String())
}
