package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/reusing-code/dochan/eml"

	"gopkg.in/abiosoft/ishell.v2"
)

func main() {

	shell := ishell.New()

	shell.Println("Search Shell")

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
