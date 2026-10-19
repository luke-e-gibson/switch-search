package main

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type SearchBang struct {
	Name       string // eg: "Google"
	Identifier string // eg: "g" -> used in the search query like "!g"
	URL        string // eg: "https://www.google.com/search?q="
}

var searchBangs = []SearchBang{
	{Name: "Google", Identifier: "g", URL: "https://www.google.com/search?q="},
	{Name: "DuckDuckGo", Identifier: "d", URL: "https://duckduckgo.com/?q="},
	{Name: "Bing", Identifier: "b", URL: "https://www.bing.com/search?q="},
	{Name: "Yahoo", Identifier: "y", URL: "https://search.yahoo.com/search?p="},
	{Name: "Wikipedia", Identifier: "w", URL: "https://en.wikipedia.org/wiki/Special:Search?search="},
}

func main() {
	app := fiber.New()

	app.Get("/", static.New("assets"))

	app.Get("/search", func(c fiber.Ctx) error {
		// Get search query from the request
		query := c.Query("q")
		// Allow a bang to appear anywhere in the query, not just at the start
		if idx := strings.Index(query, "!"); idx != -1 && idx+1 < len(query) {
			for _, bang := range searchBangs {
				if query[idx+1] == bang.Identifier[0] {
					// Redirect to the corresponding search engine using the remainder after the bang+identifier
					return c.Redirect().To(bang.URL + query[idx+2:])
				}
			}
		}
		return c.Redirect().To("https://www.google.com/search?q=" + query) // Default to Google if no bang is found
	})

	app.Listen(":3000")
}
