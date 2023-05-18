package starter2

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Route interface {
	http.Handler

	// Pattern reports the path at which this is registered.
	Pattern() string
}

type EchoHandler struct {
	log *zap.Logger
}

func (h *EchoHandler) Pattern() string {
	return "/echo"
}
func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		h.log.Warn("Failed to handle request", zap.Error(err))
	}
}

func NewEchoHandler(log *zap.Logger) *EchoHandler {
	return &EchoHandler{log: log}
}

type HelloHandler struct {
	log *zap.Logger
}

// NewHelloHandler builds a new HelloHandler.
func NewHelloHandler(log *zap.Logger) *HelloHandler {
	return &HelloHandler{log: log}
}

func (*HelloHandler) Pattern() string {
	return "/hello"
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error("Failed to read request", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "Hello, %s\n", body); err != nil {
		h.log.Error("Failed to write response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("Starting HTTP server", zap.String("addr", srv.Addr))
			//{"level":"info","msg":"Starting HTTP server","addr":":8080"}
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

// usage:
// curl -X POST -d 'hello' http://localhost:8080/echo
// output:hello
// curl -X POST -d 'gopher' http://localhost:8080/hello
// output: Hello, gopher
func Main() {
	/*
		fx.New(
			fx.Provide(
				NewHTTPServer,
				fx.Annotate(
					NewServeMux,
					fx.ParamTags(`group:"routes"`),
				),
				fx.Annotate(
					NewEchoHandler,
					fx.As(new(Route)), // convert struct EchoHandler to interface Route
					fx.ResultTags(`group:"routes"`),
				),
				fx.Annotate(
					NewHelloHandler,
					fx.As(new(Route)),
					fx.ResultTags(`group:"routes"`),
				),
				zap.NewExample,
			),
			fx.Invoke(func(*http.Server) {}),
		).Run()*/
	fx.New(
		fx.Provide(
			NewHTTPServer,
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			AsRoute(NewEchoHandler),
			AsRoute(NewHelloHandler),
			zap.NewExample,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()

}
