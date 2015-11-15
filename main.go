package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"syscall"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Production World!")
}

func main() {

	socket := "/home/vagrant/dev/src/github.com/r-fujiwara/goroutine-fcgi/go-home.sock"
	userid := 33
	groupid := 33

	_ = syscall.Umask(0177)

	l, err := net.Listen("unix", socket)

	_ = syscall.Chown(socket, userid, groupid)

	defer l.Close()

	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Captured: %v", sig)
			l.Close()
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	fcgi.Serve(l, mux)
}
