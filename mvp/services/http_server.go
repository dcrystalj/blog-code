package services

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/dcrystalj/blog-code/mvp/data"
)

type HttpServer struct {
	srv         *http.Server
	listener    net.Listener
	messageSink chan data.Comment
}

func NewHttpServer(url string, messageSink chan data.Comment) *HttpServer {
	return &HttpServer{
		srv: &http.Server{
			Addr: url,
		},
		messageSink: messageSink,
	}
}

func (h *HttpServer) Init() {
	http.HandleFunc("/comment/", h.addCommentHandler)

	listener, err := net.Listen("tcp", h.srv.Addr)
	if err != nil {
		panic(err)
	}
	h.listener = listener
}

func (h *HttpServer) Run(ctx context.Context) error {
	go h.srv.Serve(h.listener)
	for {
		select {
		case <-ctx.Done():
			return h.srv.Shutdown(ctx)
		}
	}
}

func (h *HttpServer) addCommentHandler(w http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(w, "Invalid method type", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	c := data.Comment{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&c)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	h.messageSink <- c
	w.WriteHeader(http.StatusNoContent)
}
