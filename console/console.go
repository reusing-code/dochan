package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/reusing-code/dochan/eml"

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
				err := parser.ParseDir(dir, func(f parser.File, strings []string) {
					count++
					result[f.Filename] = strings
					percent := int(count * 100 / totalCount)
					c.ProgressBar().Suffix(fmt.Sprintf(" %d%% [%d/%d]", percent, count, totalCount))
					c.ProgressBar().Progress(percent)
				}, parser.NoSkip)
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

	shell.AddCmd(&ishell.Cmd{
		Name: "parsemails",
		Help: "extract attachments from mails",
		Func: func(c *ishell.Context) {
			// disable the '>>>' for cleaner same line input.
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			c.Print("Input dir: ")
			input := c.ReadLine()

			c.Print("output dir: ")
			output := c.ReadLine()
			os.MkdirAll(output, 0777)
			num := 0

			err := eml.ExtractAttachmentsFromDirRec(input, func(filename string, content []byte, messageID string) error {
				target := filepath.Join(output, messageID+"-"+filename)
				err := ioutil.WriteFile(target, content, 0666)
				num++
				return err
			})
			if err != nil {
				c.Println(err)
			}

			c.Printf("Extracted attachments: %d\n", num)
		},
	})

	// run shell
	shell.Run()
}
