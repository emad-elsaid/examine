package examine

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	track(dbg)
	functions(dbg)
	go trace(dbg)

	fmt.Println("Examine => http://127.0.0.1:9000")
	http.ListenAndServe(":9000", nil)
}

func attachDebugger() *debugger.Debugger {
	dbg, err := debugger.New(&debugger.Config{
		AttachPid: os.Getppid(),
		Backend:   "native",
	}, nil)
	if err != nil {
		fmt.Println("Err: ", err)
		return nil
	}

	return dbg
}

func functions(dbg *debugger.Debugger) {
	funcs, _ := dbg.Functions("")
	for _, f := range funcs {
		fmt.Println(f)
	}
}

func track(dbg *debugger.Debugger) {
	fs := []string{
		"net/http.(*ServeMux).ServeHTTP",
	}

	for _, f := range fs {
		_, err := dbg.CreateBreakpoint(&api.Breakpoint{
			Name:         f,
			FunctionName: f,
		}, "", nil, false)
		if err != nil {
			log.Printf("Can't create BP: %s, %s", f, err)
		}
	}
}

func trace(dbg *debugger.Debugger) {
	for {
		dbg.TargetGroup().Continue()

		state, err := dbg.State(false)
		if err != nil {
			fmt.Println(err)
			return
		}

		if state.Exited {
			os.Exit(0)
		}

		fmt.Printf("%v\n", state.CurrentThread.Function.Name())
	}
}
