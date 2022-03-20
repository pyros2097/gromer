// Code generated by gromer. DO NOT EDIT.
package main

import (
	c "context"
	"embed"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pyros2097/gromer"
	"github.com/rs/zerolog/log"
	"gocloud.dev/server"

	"github.com/pyros2097/gromer/example/context"
	"github.com/pyros2097/gromer/example/pages/api/todos"
	"github.com/pyros2097/gromer/example/pages"
	"github.com/pyros2097/gromer/example/pages/about"
	"github.com/pyros2097/gromer/example/pages/api/recover"
	"github.com/pyros2097/gromer/example/pages/api/todos/_todoId_"
)

//go:embed assets/*
var assetsFS embed.FS

func main() {
	r := mux.NewRouter()
	r.Use(gromer.LogMiddleware)
	r.NotFoundHandler = gromer.NotFoundHandler
	r.PathPrefix("/assets/").Handler(wrapCache(http.FileServer(http.FS(assetsFS))))
	handle(r, "GET", "/api", gromer.ApiExplorer(apiDefinitions()))
	handle(r, "GET", "/about", about.GET)
	handle(r, "DELETE", "/api/todos/{todoId}", todos_todoId_.DELETE)
	handle(r, "GET", "/api/todos/{todoId}", todos_todoId_.GET)
	handle(r, "PUT", "/api/todos/{todoId}", todos_todoId_.PUT)
	handle(r, "GET", "/api/todos", todos.GET)
	handle(r, "POST", "/api/todos", todos.POST)
	handle(r, "GET", "/api/recover", recover.GET)
	handle(r, "GET", "/", pages.GET)
	println("http server listening on http://localhost:3000")
	srv := server.New(r, nil)
	if err := srv.ListenAndServe(":3000"); err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to listen")
	}
}

func wrapCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
		h.ServeHTTP(w, r)
	})
}

func handle(router *mux.Router, method, route string, h interface{}) {
	router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ctx, err := context.WithContext(c.WithValue(
			c.WithValue(
				c.WithValue(r.Context(), "assetsFS", assetsFS),
					"url", r.URL),
			"header", r.Header))
		if err != nil {
			gromer.RespondError(w, 500, err)
			return
		}
		gromer.PerformRequest(route, h, ctx, w, r)
	}).Methods(method)
}

func apiDefinitions() []gromer.ApiDefinition {
	return []gromer.ApiDefinition{
		
		{
			Method: "DELETE",
			Path: "/api/todos/{todoId}",
			PathParams: []string{ "todoId",  },
			Params: map[string]interface{}{
				
			},
		},
		{
			Method: "GET",
			Path: "/api/todos/{todoId}",
			PathParams: []string{ "todoId",  },
			Params: map[string]interface{}{
				"show": "string", 
			},
		},
		{
			Method: "PUT",
			Path: "/api/todos/{todoId}",
			PathParams: []string{ "todoId",  },
			Params: map[string]interface{}{
				"completed": "bool", 
			},
		},
		{
			Method: "GET",
			Path: "/api/todos",
			PathParams: []string{  },
			Params: map[string]interface{}{
				"limit": "int", "offset": "int", 
			},
		},
		{
			Method: "POST",
			Path: "/api/todos",
			PathParams: []string{  },
			Params: map[string]interface{}{
				"text": "string", 
			},
		},
	}
}
