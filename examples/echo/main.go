package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lafriks/go-spaproxy"
)

func main() {
	// Create new Angular development server proxy service
	proxy, err := spaproxy.NewAngularDevProxy(&spaproxy.AngularDevProxyOptions{
		Dir: "../webapps/angular/",
	})
	if err != nil {
		panic(err)
	}

	// Catch interrupts for gracefully stopping background node proecess
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start React development server
	if err = proxy.Start(context.Background()); err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{
			URL: proxy.DevServerURL(),
		},
	})))

	// Start web server on port 8080
	go func() {
		e.Logger.Fatal(e.Start(":8080"))
	}()

	<-done

	// Gracefully kill proxy background node process
	proxy.Stop()
}
