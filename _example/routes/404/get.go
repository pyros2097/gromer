package not_found_404

import (
	. "github.com/pyros2097/gromer/gsx"
)

func GET(c Context) (*Node, int, error) {
	c.Meta("title", "Page Not Found")
	return c.Render(`
		<main class="box center">
			<h1>Page Not Found</h1>
			<h2 class="mt-6">
				<a class="is-underlined" href="/">Go Back</a>
			</h2>
		</main>
	`), 404, nil
}
