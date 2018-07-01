package eml

import (
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/karrick/godirwalk"

	"github.com/jhillyerd/enmime"
)

func ExtractAttachments(eml io.Reader, cb func(filename string, content []byte, messageID string) error) error {
	env, err := enmime.ReadEnvelope(eml)
	if err != nil {
		return err
	}
	msgID := env.GetHeader("Message-ID")
	if len(msgID) == 0 {
		msgID = uuid.New().String()
	}
	for _, att := range env.Attachments {
		err = cb(att.FileName, att.Content, msgID)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExtractAttachmentsFromDirRec(dir string, cb func(filename string, content []byte, messageID string) error) error {
	err := godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				fmt.Println(path)
				f, err := os.Open(path)
				defer f.Close()
				if err != nil {
					return err
				}
				err = ExtractAttachments(f, cb)
				if err != nil {
					fmt.Printf("Error in file %q: %q\n", path, err)
					return nil
				}
			}
			return nil
		},
	})

	if err != nil {
		return err
	}

	return nil
}
