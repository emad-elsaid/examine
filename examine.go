package examine

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-delve/delve/service/debugger"
)

//go:embed public
var public embed.FS

func init() {
	if os.Getenv("DEBUGGER") == "true" {
		main()

		os.Exit(0)
	}

	fork()
}

func fork() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	env := os.Environ()
	env = append(env, "DEBUGGER=true")

	_, err = syscall.ForkExec(os.Args[0], nil, &syscall.ProcAttr{
		Dir: wd,
		Env: env,
		Files: []uintptr{
			os.Stdin.Fd(),
			os.Stdout.Fd(),
			os.Stderr.Fd(),
		},
		Sys: nil,
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	dbg := attachDebugger()
	if dbg == nil {
		os.Exit(1)
	}

	go handleExit(dbg)

	// functions(dbg)
	timeline := NewTimeline()

	funcs, err := dbg.Functions("main.")
	if err != nil {
		slog.Error(err.Error())
	}
	tracers := []Tracer{}
	for _, f := range funcs {
		if strings.Contains(f, "github.com/emad-elsaid/examine") {
			continue
		}
		tracers = append(tracers, &GenericTracer{functionName: f})
	}

	tracers = append(tracers, httpTracers...)

	go trace(dbg, timeline, tracers...)

	http.HandleFunc("/timeline.json", func(w http.ResponseWriter, r *http.Request) {
		points, err := timeline.JSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")
		fmt.Fprintf(w, string(points))
	})

	http.Handle("/public/", http.FileServerFS(public))
	http.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/public/index.html", http.StatusFound)
	})

	fmt.Println("Examine => http://127.0.0.1:9000")
	http.ListenAndServe(":9000", nil)
}

func attachDebugger() *debugger.Debugger {
	dbg, err := debugger.New(&debugger.Config{
		AttachPid: os.Getppid(),
		Backend:   "default",
	}, nil)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	return dbg
}

func handleExit(dbg *debugger.Debugger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		dbg.Detach(true)
	}
}

func functions(dbg *debugger.Debugger) {
	funcs, _ := dbg.Functions("")
	for _, f := range funcs {
		slog.Info(f)
	}
}
