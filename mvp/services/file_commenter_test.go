package services

import (
	"context"
	"os"
	"testing"

	"github.com/dcrystalj/blog-code/mvp/data"
	"github.com/stretchr/testify/assert"
)

const filePath = "./test_file"

func TestOpenFileForWriting(t *testing.T) {
	defer os.Remove(filePath)
	fc := NewFileCommenter(filePath)
	fc.Init()

	fc.file.WriteString("Can write")
	fc.file.Close()
}

func TestWriteComment(t *testing.T) {
	defer os.Remove(filePath)
	comment := data.Comment{
		Uid:     "123",
		Comment: "Test comment",
	}
	fc := NewFileCommenter(filePath)
	fc.Init()

	fc.writeComment(comment)

	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "{\"uid\":\"123\",\"comment\":\"Test comment\"}\n", string(content))
}

func TestServiceWrite3Comments(t *testing.T) {
	defer os.Remove(filePath)
	comment := data.Comment{
		Uid:     "123",
		Comment: "Test comment",
	}
	ctx, cancel := context.WithCancel(context.Background())
	fc := NewFileCommenter(filePath)
	fc.Init()
	go fc.Run(ctx)

	for i := 0; i < 3; i++ {
		fc.Input <- comment
	}

	cancel()

	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "{\"uid\":\"123\",\"comment\":\"Test comment\"}\n{\"uid\":\"123\",\"comment\":\"Test comment\"}\n{\"uid\":\"123\",\"comment\":\"Test comment\"}\n", string(content))
}

type FileCommenterError struct {
	FileCommenter
}

func (fc *FileCommenterError) writeComment(message data.Comment) error {
	return writeErr
}

func TestRunServiceReturnsOnWriteError(t *testing.T) {
	ctx := context.Background()
	fc := &FileCommenterError{}
	fc.Input = make(chan data.Comment)
	returnedError := make(chan error, 1)
	go func() {
		err := fc.Run(ctx)
		returnedError <- err
	}()
	fc.Input <- data.Comment{}

	assert.Equal(t, writeErr, <-returnedError)
}
