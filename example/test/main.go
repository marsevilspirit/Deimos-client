package main

import (
	"context"
	"log/slog"

	deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)

	ctx := context.Background()

	dir := "/phobos"

	_, _ = client.Delete(ctx, dir, deimos.WithDir(), deimos.WithRecursive())
	resp, err := client.Set(ctx, dir, "", deimos.WithDir())
	if err != nil {
		slog.Error("set dir", "err", err, "resp", resp)
		return
	}
	slog.Info("set dir", "resp", resp)

	resp, err = client.Set(ctx, dir+"/foo", "bar")
	if err != nil {
		slog.Error("set foo", "err", err, "resp", resp)
		return
	}
	slog.Info("set foo", "resp", resp)

	resp, err = client.Get(ctx, dir+"/foo")
	if err != nil {
		slog.Error("get dir", "err", err, "resp", resp)
		return
	}
	slog.Info("get", "resp", resp)
}
