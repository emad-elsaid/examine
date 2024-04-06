# Examine

A drop-in Go package to trace your program automatically.

> [!CAUTION]
> Examine is still in very early development

# Problem
1. When working on a project I would like to trace execution especially for network requests (incoming and outgoing)
2. The way to do it is either change the code to log these requests or add breakpoints and attach a debugger
3. I would like to Do 1 without doing 2. I would like to just import a package that traces the program network requests

# How to use
1. Import Examine to your program with `import _ "github.com/emad-elsaid/examine"`
2. Build and run your program `go build ./cmd/path/to/cmd && ./cmd-name`
3. Examine prints traces of your network requests (this can change during development)
4. As Examine pauses programs it will affect performance. and shouldn't be used in production for any case


# How it works
1. Examine forks your program
2. Attaches Delve debugger to your program
3. Adds breakpoints to Go standard library
4. Whenever a breakpoint is hit. it'll store relevant information and print them and continue

# Known issues
## Can't print any variables if the program ran with `go run`
Turns out `go run` doesn't include DWARF debugging information. so you have to `go build` then run the program

## Couldn't attach to PID
If you get this message:
>ERROR Could not attach to pid 10042: this could be caused by a kernel security setting, try writing "0" to /proc/sys/kernel/yama/ptrace_scope

You need to turn on `ptrace_scope` in your kernel by
```shell
sudo su
echo "0" > /proc/sys/kernel/yama/ptrace_scope
```
