# go-spaproxy

Generic GoLang middleware to reverse proxy SPA development files to allow hot reload functionality while developing.

[![Build Status](https://cloud.drone.io/api/badges/lafriks/go-spaproxy/status.svg)](https://cloud.drone.io/lafriks/go-spaproxy)
[![codecov](https://codecov.io/gh/lafriks/go-spaproxy/branch/master/graph/badge.svg)](https://codecov.io/gh/lafriks/go-spaproxy)

## Supported JavaScript frameworks

Any framework can be supported as long as it has development web server to start but additional helpers are available for following frameworks:

* [Vue.js](https://vuejs.org/) - `NewVueDevProxy`
* [React](https://reactjs.org/) - `NewReactDevProxy`
* [Angular](https://angular.io/) - `NewAngularDevProxy`
* [Svelte](https://svelte.dev/) - `NewSvelteDevProxy`

## Usage

To use proxy instance needs to be created and then started using `SpaDevProxy.Start()` method.

Later You can use `SpaDevProxy.HandleFunc(w http.ResponseWriter, r *http.Request)` method to add it to applications *catch all* route that would proxy all not application requests to background development server.

**NB!** Application need to gracefully shutdown and call `SpaDevProxy.Stop()` method otherwise started node background server will not be stopped.

For examples on how to integrate see [examples](examples) folder:

* Net/HTTP example - [examples/simple](examples/simple/main.go)
* [Chi](https://github.com/go-chi/chi) example - [examples/go-chi](examples/go-chi/main.go)
* [Echo](https://echo.labstack.com/) example - [examples/echo](examples/echo/main.go)
* [Gin Web Framework](https://gin-gonic.com/) example - [examples/gin](examples/gin/main.go)
