package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dcrystalj/blog-code/mvp/data"
	"github.com/stretchr/testify/assert"
)

func TestAddCommentHandlerGETFails(t *testing.T) {
	h := &HttpServer{}
	req := httptest.NewRequest("GET", "/comment/", nil)
	w := httptest.NewRecorder()

	h.addCommentHandler(w, req)

	resp := w.Result()
	assert.Equal(t, 405, resp.StatusCode)
}

func TestPostComment(t *testing.T) {
	c := make(chan data.Comment, 1)
	h := &HttpServer{messageSink: c}
	postComment, _ := json.Marshal(data.Comment{Uid: "123", Comment: "Some comment"})
	req := httptest.NewRequest("POST", "/comment/", bytes.NewReader(postComment))
	w := httptest.NewRecorder()

	h.addCommentHandler(w, req)

	resp := w.Result()
	assert.Equal(t, 204, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	comment := <-h.messageSink
	assert.Equal(t, "123", comment.Uid)
	assert.Equal(t, "Some comment", comment.Comment)
}

func TestServerRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan data.Comment, 1)
	s := NewHttpServer(":8083", c)
	s.Init()
	go s.Run(ctx)

	time.Sleep(time.Duration(100 * time.Millisecond))
	cancel()
}
