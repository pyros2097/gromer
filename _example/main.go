// Code generated by gromer. DO NOT EDIT.
package main

import (
	"os"

	"github.com/gorilla/mux"
	"github.com/pyros2097/gromer"
	"github.com/rs/zerolog/log"
	"gocloud.dev/server"

	"github.com/pyros2097/gromer/_example/assets"
	"github.com/pyros2097/gromer/_example/components"
	"github.com/pyros2097/gromer/_example/pages/404"
	"github.com/pyros2097/gromer/_example/pages"
	"github.com/pyros2097/gromer/_example/pages/about"
	"github.com/pyros2097/gromer/_example/pages/api/recover"
	"github.com/pyros2097/gromer/_example/pages/api/todos"
	"github.com/pyros2097/gromer/_example/pages/api/todos/_todoId_"
	
)

func init() {
	gromer.RegisterComponent(components.Header)
	gromer.RegisterComponent(components.Page)
	gromer.RegisterComponent(components.Todo)
	
}

func main() {
	port := os.Getenv("PORT")
	baseRouter := mux.NewRouter()
	baseRouter.Use(gromer.LogMiddleware)
	
	baseRouter.NotFoundHandler = gromer.StatusHandler(not_found_404.GET)
	
	staticRouter := baseRouter.NewRoute().Subrouter()
	staticRouter.Use(gromer.CacheMiddleware)
	gromer.StaticRoute(staticRouter, "/assets/", assets.FS)
	gromer.StylesRoute(staticRouter, "/styles.css")

	pageRouter := baseRouter.NewRoute().Subrouter()
	gromer.ApiExplorerRoute(pageRouter, "/explorer")
	gromer.Handle(pageRouter, "GET", "/", pages.GET)
	gromer.Handle(pageRouter, "GET", "/about", about.GET)
	

	apiRouter := baseRouter.NewRoute().Subrouter()
	apiRouter.Use(gromer.CorsMiddleware)
	gromer.Handle(apiRouter, "GET", "/api/recover", recover.GET)
	gromer.Handle(apiRouter, "GET", "/api/todos", todos.GET)
	gromer.Handle(apiRouter, "POST", "/api/todos", todos.POST)
	gromer.Handle(apiRouter, "DELETE", "/api/todos/{todoId}", todos_todoId_.DELETE)
	gromer.Handle(apiRouter, "GET", "/api/todos/{todoId}", todos_todoId_.GET)
	gromer.Handle(apiRouter, "PUT", "/api/todos/{todoId}", todos_todoId_.PUT)
	
	
	
	log.Info().Msg("http server listening on http://localhost:"+port)
	srv := server.New(baseRouter, nil)
	if err := srv.ListenAndServe(":"+port); err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to listen")
	}
}
