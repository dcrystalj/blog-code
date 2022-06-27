package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/dcrystalj/blog-code/mvp/services"
	"golang.org/x/sync/errgroup"
)

func main() {
	url := flag.String("url", ":8080", "hostname:port")
	path := flag.String("path", "./db/comments", "folder to save comments")
	flag.Parse()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	fc := services.NewFileCommenter(*path)
	srv := services.NewHttpServer(*url, fc.Input)

	fc.Init()
	srv.Init()

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return fc.Run(gctx)
	})
	g.Go(func() error {
		return srv.Run(gctx)
	})
	if err := g.Wait(); err != nil && err != context.Canceled {
		panic(err)
	}
}
