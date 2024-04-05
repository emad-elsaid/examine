package examine

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/debugger"
)

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})

	dbg := attachDebugger()
	if dbg == nil {
		os.Exit(1)
	}

	go handleExit(dbg)

	// functions(dbg)
	trace(dbg, httpTracers)

	// fmt.Println("Examine => http://127.0.0.1:9000")
	// http.ListenAndServe(":9000", nil)
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

func printvar(vars []api.Variable) {
	for _, v := range vars {
		slog.Info("var",
			"name", v.Name,
			"type", v.Type,
		)
		fmt.Println(v.MultilineString("\t", ""))
	}
}
