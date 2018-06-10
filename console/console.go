package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/reusing-code/dochan/persist"

	"github.com/reusing-code/dochan/parser"

	"gopkg.in/abiosoft/ishell.v2"
)

func main() {

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
			//query := c.ReadLine()

			//result := s.Search(query, true)
			//for _, r := range result {
			//	c.Println(r)
			//}

		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "parse",
		Help: "parse dir",
		Func: func(c *ishell.Context) {
			// disable the '>>>' for cleaner same line input.
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// get username
			c.Print("Dir: ")
			dir := c.ReadLine()
			c.ProgressBar().Start()
			result := make(map[interface{}][]string)
			totalCount, _ := parser.GetFileCount(dir)
			if totalCount <= 0 {
				totalCount = 1 // no division by 0...
			}
			count := 0
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := parser.ParseDir(dir, func(filename string, strings []string) {
					count++
					result[filename] = strings
					percent := int(count * 100 / totalCount)
					c.ProgressBar().Suffix(fmt.Sprintf(" %d%% [%d/%d]", percent, count, totalCount))
					c.ProgressBar().Progress(percent)
				})
				if err != nil {
					c.Println(err)
				}
			}()
			wg.Wait()
			c.ProgressBar().Stop()

			output, err := os.OpenFile("searchdb.gob.gz", os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				c.Println(err)
			}
			defer output.Close()

			err = persist.PersistData(result, output)
			if err != nil {
				c.Println(err)
			}
		},
	})

	// run shell
	shell.Run()
}
