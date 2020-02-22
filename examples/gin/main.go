package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/lafriks/go-spaproxy"
)

func main() {
	// Create new Svelte development server proxy service
	proxy, err := spaproxy.NewSvelteDevProxy(&spaproxy.SvelteDevProxyOptions{
		Dir: "../webapps/svelte/",
	})
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.NoRoute(func(c *gin.Context) {
		proxy.HandleFunc(c.Writer, c.Request)
	})

	// Catch interrupts for gracefully stopping background node proecess
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start React development server
	if err = proxy.Start(context.Background()); err != nil {
		panic(err)
	}

	// Start web server on port 8080
	go func() {
		if err = router.Run(":8080"); err != nil {
			panic(err)
		}
	}()

	<-done

	// Gracefully kill proxy background node process
	proxy.Stop()
}
