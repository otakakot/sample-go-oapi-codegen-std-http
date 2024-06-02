package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/otakakot/sample-go-oapi-codegen-std-http/pkg/api"
)

func main() {
	port := cmp.Or(os.Getenv("PORT"), "8080")

	serv := &Server{}

	mux := http.NewServeMux()

	hdl := api.HandlerWithOptions(serv, api.StdHTTPServerOptions{
		BaseRouter: mux,
	})

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           hdl,
		ReadHeaderTimeout: 30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	go func() {
		slog.Info("start server listen")

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	slog.Info("start server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

	slog.Info("done server shutdown")
}

var _ api.ServerInterface = (*Server)(nil)

type Server struct{}

// Health implements api.ServerInterface.
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
