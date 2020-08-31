package hotstart

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	LISTENER_FD           = 3
	DEFAULT_READ_TIMEOUT  = 60 * time.Second
	DEFAULT_WRITE_TIMEOUT = DEFAULT_READ_TIMEOUT
)

var (
	runMutex = sync.RWMutex{}
)

// HTTP server that supported hotstart shutdown or restart
func ListenAndServe(addr string, handler http.Handler) error {

	return NewServer(addr, handler, DEFAULT_READ_TIMEOUT, DEFAULT_WRITE_TIMEOUT).Server()
}

/*
Listen http server
*/
func (srv *HotServer) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := srv.getNetListener(addr)
	if err != nil {
		return err
	}

	srv.listener = ln

	srv.BeforeBegin(srv.Addr)

	return srv.Serve()
}

/*
监听 https server
*/
func (srv *HotServer) ListenAndServeTLS(certFile, keyFile string) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	ln, err := srv.getNetListener(addr)
	if err != nil {
		return err
	}

	srv.listener = tls.NewListener(ln, config)

	srv.BeforeBegin(srv.Addr)

	return srv.Serve()
}

/*
服务启动
*/
func (srv *HotServer) Serve() error {
	//监听信号
	go srv.handleSignals()
	err := srv.Server.Serve(srv.listener)

	srv.logf("waiting for connections closed.")
	//阻塞等待关闭
	<-srv.shutdownChan
	srv.logf("all connections closed.")

	return err
}

/*
get lister
*/
func (srv *HotServer) getNetListener(addr string) (ln net.Listener, err error) {
	if srv.isChild {
		file := os.NewFile(LISTENER_FD, "")
		ln, err = net.FileListener(file)
		if err != nil {
			err = fmt.Errorf("net.FileListener error: %v", err)
			return nil, err
		}
	} else {
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			err = fmt.Errorf("net.Listen error: %v", err)
			return nil, err
		}
	}
	return ln, nil
}

/*
监听信号
*/

func (srv *HotServer) handleSignals() {
	var sig os.Signal

	signal.Notify(
		srv.signalChan,
		syscall.SIGINT,
	)

	for {
		sig = <-srv.signalChan
		switch sig {
		case syscall.SIGINT:
			srv.logf("received SIGINT, hotstart shutting down HTTP server.")
			srv.shutdown()
		default:
		}
	}
}

/*
优雅关闭后台
*/
func (srv *HotServer) shutdown() {
	if err := srv.Shutdown(context.Background()); err != nil {
		srv.logf("HTTP server shutdown error: %v", err)
	} else {
		srv.logf("HTTP server shutdown success.")
		srv.shutdownChan <- true
	}
}

/*
格式化输出Log
*/

func (srv *HotServer) logf(format string, args ...interface{}) {
	pids := strconv.Itoa(os.Getpid())
	format = "[pid " + pids + "] " + format
	log.Printf(format, args...)
}
