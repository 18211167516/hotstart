package hotStart

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
)

type Hot struct {
	srv *http.Server
	net.Listener
}

const (
	GRACEFUL_ENVIRON_KEY    = "IS_GRACEFUL"
	GRACEFUL_ENVIRON_STRING = GRACEFUL_ENVIRON_KEY + "=1"
	GRACEFUL_LISTENER_FD    = 3
)

func (this *Hot) getTCPListenerFd() (uintptr, error) {
	f, err := this.Listener.(*net.TCPListener).File()
	if err != nil {
		return 0, fmt.Errorf("failed to get socket file descriptor: %v", err)
	}
	return f.Fd(), nil
}

// 启动子进程执行新程序
func (this *Hot) startNewProcess() error {

	listenerFd, err := this.getTCPListenerFd()
	if err != nil {
		return fmt.Errorf("failed to get socket file descriptor: %v", err)
	}

	envs := []string{}
	for _, value := range os.Environ() {
		if value != GRACEFUL_ENVIRON_STRING {
			envs = append(envs, value)
		}
	}
	envs = append(envs, GRACEFUL_ENVIRON_STRING)

	execSpec := &syscall.ProcAttr{
		Env:   envs,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFd},
	}

	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		return 0, fmt.Errorf("failed to forkexec: %v", err)
	}

	return uintptr(fork), nil
}
