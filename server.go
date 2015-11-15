package main

import (
    "fmt"
    "net"
    "net/http"
    "net/http/fcgi"
    "os"
    "os/signal"
    "syscall"
)

type Server struct {
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request){ 
    str := "Hello World!"
    w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
    w.Header().Set("Content-Length", fmt.Sprint(len(str)))
    fmt.Fprint(w, str)
}

func main(){
    const SOCK = "/home/vagrant/dev/src/github.com/r-fujiwara/goroutine-fcgi/go-home.sock"
    sig := make(chan os.Signal)
    signal.Notify(sig, os.Interrupt)
    signal.Notify(sig, syscall.SIGTERM)

    server := Server{}
    l, _ := net.Listen("unix", SOCK)

    go func() {
        fcgi.Serve(l, server)
    }()

    <-sig

    if err := os.Remove(SOCK); err != nil {
        panic("socket file remove error.")
    }
}

