package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/lafriks/go-spaproxy"
)

func main() {
	// Create new VueJS development server proxy service
	proxy, err := spaproxy.NewVueDevProxy(&spaproxy.VueDevProxyOptions{
		Dir: "../webapps/vue/",
	})
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Get("/*", proxy.HandleFunc)

	// Catch interrupts for gracefully stopping background node proecess
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start React development server
	if err = proxy.Start(context.Background()); err != nil {
		panic(err)
	}

	// Start web server on port 8080
	go func() {
		if err = http.ListenAndServe(":8080", r); err != nil {
			panic(err)
		}
	}()

	<-done

	// Gracefully kill proxy background node process
	proxy.Stop()
}
