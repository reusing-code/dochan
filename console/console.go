package main

import (
	"flag"

	"github.com/reusing-code/dochan/search"
	"gopkg.in/abiosoft/ishell.v2"
)

func main() {
	dirPtr := flag.String("d", "", "input dir")

	flag.Parse()

	s, _ := search.NewDirectorySearch(*dirPtr)

	shell := ishell.New()

	shell.Println("Search Shell")

	shell.AddCmd(&ishell.Cmd{
		Name: "search",
		Help: "search pdfs",
		Func: func(c *ishell.Context) {
			// disable the '>>>' for cleaner same line input.
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// get username
			c.Print("SearchQuery: ")
			query := c.ReadLine()

			result := s.Search(query, true)
			for _, r := range result {
				c.Println(r)
			}

		},
	})

	// run shell
	shell.Run()
}
