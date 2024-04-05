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
func (s ServeHTTP) trace(d *debugger.Debugger, state *api.DebuggerState) {
	slog.Info(s.function())
}
