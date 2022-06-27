package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/dcrystalj/blog-code/mvp/data"
)

var writeErr = errors.New("Writing comment failed")

type FileCommenter struct {
	Input chan data.Comment
	path  string
	file  *os.File
}

func NewFileCommenter(path string) *FileCommenter {
	return &FileCommenter{
		Input: make(chan data.Comment),
		path:  path,
	}
}

func (fc *FileCommenter) Init() {
	file, err := os.OpenFile(fc.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print("Could not open file")
		return
	}
	fc.file = file
}

func (fc *FileCommenter) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fc.file.Close()
		case m := <-fc.Input:
			err := fc.writeComment(m)
			switch {
			case err == writeErr:
				return err // something is wrong with file. To prevent furhter dataloss we stop the service.
			case err != nil:
				log.Println(err)
			}
		}
	}
}

func (fc *FileCommenter) writeComment(message data.Comment) error {
	data, err := fc.serializeComment(message)
	if err != nil {
		return err
	}
	data = append(data, bytes.NewBufferString("\n").Bytes()...)
	_, err = fc.file.Write(data)
	if err != nil {
		return writeErr
	}
	return nil
}

func (fc *FileCommenter) serializeComment(m data.Comment) ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, errors.New("Could not parse data")
	}
	return data, nil
}
