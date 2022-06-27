package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/dcrystalj/blog-code/mvp/data"
	"github.com/dcrystalj/blog-code/mvp/services"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

const path = "./db/test_comments"
const serverPort = 12345

func TestPostComments(t *testing.T) {
	os.Remove(path)
	defer os.Remove(path)

	ctx, cancel := context.WithCancel(context.Background())

	fc := services.NewFileCommenter(path)
	srv := services.NewHttpServer(fmt.Sprintf(":%d", serverPort), fc.Input)

	fc.Init()
	srv.Init()

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return fc.Run(gctx)
	})
	g.Go(func() error {
		return srv.Run(gctx)
	})

	comments := post100Comments(t)
	cancel()

	if err := g.Wait(); err != nil && err != context.Canceled {
		panic(err)
	}

	lines := readCommentsFromFile(t)
	assert.Equal(t, len(comments), len(lines))

	for _, line := range lines {
		_, found := comments[string(line)]
		assert.True(t, found)
	}
}

func post100Comments(t *testing.T) map[string]struct{} {
	contents := make(map[string]struct{})
	wgPost := sync.WaitGroup{}
	wgPost.Add(100)
	for i := 1; i <= 10; i++ {
		for j := 0; j < 10; j++ {
			go func(a, b int) {
				defer wgPost.Done()
				serializedComment, err := json.Marshal(data.Comment{
					Comment: fmt.Sprintf("test comment %d", a),
					Uid:     fmt.Sprint(b),
				})
				assert.NoError(t, err)
				c := http.Client{}
				response, err := c.Post(fmt.Sprintf("http://localhost:%d/comment/", serverPort), "application/json", bytes.NewBuffer(serializedComment))
				assert.NoError(t, err)
				contents[string(serializedComment)] = struct{}{}
				assert.Equal(t, 204, response.StatusCode)
			}(i, j)
		}
	}
	wgPost.Wait()
	return contents
}

func readCommentsFromFile(t *testing.T) [][]byte {
	content, err := os.ReadFile(path)
	assert.NoError(t, err)
	lines := bytes.Split(content, bytes.NewBufferString("\n").Bytes())
	return lines[:len(lines)-1]
}
