# gromer

[![Version](https://badge.fury.io/gh/pyros2097%2Fgromer.svg)](https://github.com/pyros2097/gromer)

**gromer** is a framework and cli to build web apps in golang.
It uses a declarative syntax using inline templates for components and pages.
It also generates http handlers for your routes which follow a particular folder structure. Similar to other frameworks like nextjs, sveltekit.
These handlers are also normal functions and can be imported in other packages directly. ((inspired by [Encore](https://encore.dev/)).
More information on the templating syntax is given [here](https://github.com/pyrossh/gromer/blob/master/handlebars/README.md),

# Requirements

```sh
go >= v1.18
```

# Install

```sh
go get -u -v github.com/pyros2097/gromer/cmd/gromer
```

# Using

You need to follow this directory structure similar to nextjs for the api route handlers to be generated and run the gromer command.

[Example](https://github.com/pyros2097/gromer/tree/master/_example)


**These are normal page routes**
```go
// /pages/get.go
package todos_page

import (
	"context"

	. "github.com/pyros2097/gromer"
	_ "github.com/pyros2097/gromer/_example/components"
	"github.com/pyros2097/gromer/_example/pages/api/todos"
	. "github.com/pyros2097/gromer/handlebars"
)

type GetParams struct {
	Filter string `json:"filter"`
	Page   int    `json:"limit"`
}

func GET(ctx context.Context, params GetParams) (HtmlContent, int, error) {
	index := Default(params.Page, 1)
	todos, status, err := todos.GET(ctx, todos.GetParams{
		Filter: params.Filter,
		Limit:  index * 10,
	})
	if err != nil {
		return HtmlErr(status, err)
	}
	return Html(`
		<Page title="gromer example">
			<Header></Header>
			<section class="todoapp">
					<section class="main">
						<ul class="todo-list" id="todo-list">
							{{#each todos as |todo|}}
								{{#Todo todo=todo}}{{/Todo}}
							{{/each}}
						</ul>
					</section>
				{{/if}}
			</section>
		</Page>
		`).
		Prop("todos", todos).
		Render()
}
```


**These are API routes**
```go
// /pages/api/todos/get.go
package todos

import (
	"context"

	. "github.com/pyros2097/gromer"
	"github.com/pyros2097/gromer/_example/services"
)

type GetParams struct {
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
}

func GET(ctx context.Context, params GetParams) ([]*services.Todo, int, error) {
	limit := Default(params.Limit, 10)
	todos := services.GetAllTodo(ctx, services.GetAllTodoParams{
		Limit: limit,
	})
	if params.Filter == "completed" {
		newTodos := []*services.Todo{}
		for _, v := range todos {
			if v.Completed {
				newTodos = append(newTodos, v)
			}
		}
		return newTodos, 200, nil
	}
	if params.Filter == "active" {
		newTodos := []*services.Todo{}
		for _, v := range todos {
			if !v.Completed {
				newTodos = append(newTodos, v)
			}
		}
		return newTodos, 200, nil
	}
	return todos, 200, nil
}

```

```go
// /pages/api/todos/post.go
package todos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pyros2097/gromer/_example/services"
)

type PostParams struct {
	Text string `json:"text"`
}

func POST(ctx context.Context, b PostParams) (*services.Todo, int, error) {
	todo, err := services.CreateTodo(ctx, services.Todo{
		ID:        uuid.New().String(),
		Text:      b.Text,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, 500, err
	}
	return todo, 200, nil
}
```

And then run the gromer cli command annd it will generate the route handlers in a main.go file,
```go
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
```

# TODO:
Add inline css formatting
ADd inline html formatting