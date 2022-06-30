package gsx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TodoData struct {
	ID        string
	Text      string
	Completed bool
}

func Todo(c *Context, todo *TodoData) []*Tag {
	return c.Render(`
		<li id="todo-{todo.ID}" class="{ completed: todo.Completed }">
			<div class="upper">
				<span>{todo.Text}</span>
				<span>{todo.Text}</span>
			</div>
			{children}
			<div class="bottom">
				<span>{todo.Completed}</span>
				<span>{todo.Completed}</span>
			</div>
		</li>
	`)
}

func TodoList(ctx Context, todos []*TodoData) []*Tag {
	return ctx.Render(`
		<ul id="todo-list" class="relative" x-for="todo in todos">
			<Todo />
		</ul>
	`)
}

func TodoCount(ctx Context, count int) []*Tag {
	return ctx.Render(`
		<span id="todo-count" class="todo-count" hx-swap-oob="true">
			<strong>{count}</strong> item left
		</span>
	`)
}

func WebsiteName() string {
	return "My Website"
}

func TestComponent(t *testing.T) {
	r := require.New(t)
	RegisterComponent(Todo, nil, "todo")
	RegisterFunc(WebsiteName)
	h := Context{
		data: map[string]interface{}{
			"todo": &TodoData{ID: "4", Text: "My fourth todo", Completed: false},
		},
	}
	actual := renderString(h.Render(`
		<div>
			<Todo>
				<div class="todo-panel">
					<span>{todo.Text}</span>
					<span>{todo.Completed}</span>
				</div>
			</Todo>
			<Todo />
		</div>
	`))
	expected := `
		<div>
			<todo>
				<li id="todo-4" class="">
					<div class="view"><span>My fourth todo</span><span>My fourth todo</span></div>
					<div class="todo-panel"><span>My fourth todo</span><span>false</span></div>
					<div class="count"><span>false</span><span>false</span></div>
				</li>
			</todo>
		</div>
	`
	r.Equal(expected, actual)
}

// func TestFor(t *testing.T) {
// 	r := require.New(t)
// 	RegisterComponent(Todo, nil, "todo")
// 	RegisterFunc(WebsiteName)
// 	h := Context{
// 		data: map[string]interface{}{
// 			"todos": []*TodoData{
// 				{ID: "1", Text: "My first todo", Completed: true},
// 				{ID: "2", Text: "My second todo", Completed: false},
// 				{ID: "3", Text: "My third todo", Completed: false},
// 			},
// 		},
// 	}
// 	actual := h.Render(`
// 		<div>
// 			<ul x-for="todo in todos" class="relative">
// 				<li>
// 					<span>{todo.Text}</span>
// 					<span>{todo.Completed}</span>
// 					<a>link to {todo.ID}</a>
// 				</li>
// 			</ul>
// 			<ol x-for="todo in todos">
// 				<Todo>
// 					<div class="todo-panel">
// 						<span>{todo.Text}</span>
// 						<span>{todo.Completed}</span>
// 					</div>
// 				</Todo>
// 			</ol>
// 		</div>
// 	`).String()
// 	expected := stripWhitespace(`
// 		<div>
// 			<ul x-for="todo in todos" class="relative">
// 				<li><span>My first todo</span><span>true</span><a>link to 1</a></li>
// 				<li><span>My second todo</span><span>false</span><a>link to 2</a></li>
// 				<li><span>My third todo</span><span>false</span><a>link to 3</a></li>
// 			</ul>
// 			<ol x-for="todo in todos">
// 				<li id="todo-1" class="completed">
// 					<div class="view"><span>My first todo</span><span>My first todo</span></div>
// 					<div class="todo-panel"><span>My first todo</span><span>true</span></div>
// 					<div class="count"><span>true</span><span>true</span></div>
// 				</li>
// 				<li id="todo-2" class="">
// 					<div class="view"><span>My second todo</span><span>My second todo</span></div>
// 					<div class="todo-panel"><span>My second todo</span><span>false</span></div>
// 					<div class="count"><span>false</span><span>false</span></div>
// 				</li>
// 				<li id="todo-3" class="">
// 					<div class="view"><span>My third todo</span><span>My third todo</span></div>
// 					<div class="todo-panel"><span>My third todo</span><span>false</span></div>
// 					<div class="count"><span>false</span><span>false</span></div>
// 				</li>
// 			</ol>
// 		</div>
// 	`)
// 	r.Equal(expected, actual)
// }

// func TestForComponent(t *testing.T) {
// 	r := require.New(t)
// 	RegisterComponent(Todo, nil, "todo")
// 	RegisterComponent(TodoList, nil, "todos")
// 	RegisterFunc(WebsiteName)
// 	h := Context{
// 		data: map[string]interface{}{
// 			"todos": []*TodoData{
// 				{ID: "1", Text: "My first todo", Completed: true},
// 				{ID: "2", Text: "My second todo", Completed: false},
// 				{ID: "3", Text: "My third todo", Completed: false},
// 			},
// 		},
// 	}
// 	actual := h.Render(`
// 		<div>
// 			<TodoList />
// 		</div>
// 	`).String()
// 	expected := stripWhitespace(`
// 		<div>
// 			<ul id="todo-list" class="relative" x-for="todo in todos">
// 				<li id="todo-1" class="completed">
// 					<div class="view"><span>My first todo</span><span>My first todo</span></div>
// 					<div class="todo-panel"><span>My first todo</span><span>true</span></div>
// 					<div class="count"><span>true</span><span>true</span></div>
// 				</li>
// 				<li id="todo-2" class="">
// 					<div class="view"><span>My second todo</span><span>My second todo</span></div>
// 					<div class="todo-panel"><span>My second todo</span><span>false</span></div>
// 					<div class="count"><span>false</span><span>false</span></div>
// 				</li>
// 				<li id="todo-3" class="">
// 					<div class="view"><span>My third todo</span><span>My third todo</span></div>
// 					<div class="todo-panel"><span>My third todo</span><span>false</span></div>
// 					<div class="count"><span>false</span><span>false</span></div>
// 				</li>
// 			</ul>
// 		</div>
// 	`)
// 	r.Equal(expected, actual)
// }

// func TestMultipleComonent(t *testing.T) {
// 	r := require.New(t)
// 	RegisterComponent(Todo, nil, "todo")
// 	RegisterComponent(TodoCount, nil, "count")
// 	h := Context{
// 		data: map[string]interface{}{
// 			"todo": &TodoData{
// 				ID:        "3",
// 				Text:      "My third todo",
// 				Completed: false,
// 			},
// 		},
// 	}
// 	actual := h.Render(`
// 			<Todo />
// 			<TodoCount />
// 	`).String()
// 	expected := stripWhitespace(`
// 	`)
// 	r.Equal(expected, actual)
// }
