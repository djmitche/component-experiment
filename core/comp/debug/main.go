package debug

import (
	"comps/core"
	"context"
	"fmt"
	"net/http"
	"strings"
)

func componentPath(suffix string) core.ComponentPath {
	return core.ComponentPath("core/comp/debug." + suffix)
}

// Main is the component implementation for this package (`comp/debug.Main`).
//
// This component manages an `http.Handler` containing component debug
// information.  It is up to the caller to configure a server for this handler,
// using the `HandlerRequest` and `HandlerResponse` messages, or to start a
// server from this component with the `Serve` message.
//
// Other components can register handlers on this component with the
// `RegisterHandler` message.  Other components in this package do exactly
// that.
var Main = core.ComponentImpl{
	Path:         componentPath("Main"),
	Dependencies: []core.ComponentPath{},
	Start: func(map[core.ComponentPath]core.ComponentReference) core.Component {
		m := &main{
			handler:    http.NewServeMux(),
			registered: make(map[string]string),
		}
		m.register("", "/", http.HandlerFunc(m.root))
		return m
	},
}

type main struct {
	handler    *http.ServeMux
	registered map[string]string
}

var _ core.Component = &main{}
var _ core.ComponentReference = &main{}

// NewReference implements core.Component#NewReference.
func (m *main) NewReference() core.ComponentReference {
	return m
}

// Request implements core.ComponentReference#Request.
func (m *main) Request(ctx context.Context, msg core.Message) (core.Message, error) {
	switch v := msg.(type) {
	case HandlerRequest:
		return HandlerResponse{m.handler}, nil
	case Serve:
		m.serve(v.Port)
		return nil, nil
	case RegisterHandler:
		m.register(v.Name, v.Pattern, v.Handler)
		return nil, nil
	default:
		return nil, fmt.Errorf("Unrecognized message type %T", msg)
	}
}

// RequestAsync implements core.ComponentReference#RequestAsync.
func (m *main) RequestAsync(ctx context.Context, msg core.Message) {
	m.Request(ctx, msg)
}

func (m *main) serve(port int) {
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: m.handler,
	}
	go s.ListenAndServe()
}

func (m *main) register(name, pattern string, handler http.Handler) {
	m.handler.Handle(pattern, handler)
	if name != "" {
		m.registered[name] = pattern
	}
}

func (m *main) root(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	toc := []string{
		"<html>",
		"<head><title>Comps Debug</title></head>",
		"<body>",
		"<h1>Comps Debug</h1>",
		"<ul>",
	}
	for name, path := range m.registered {
		toc = append(toc, fmt.Sprintf("  <li><a href=\"%s\">%s</a></li>", path[1:], name))
	}
	toc = append(toc,
		"</ul>",
		"</body>",
		"</html>",
	)
	fmt.Fprintf(w, strings.Join(toc, "\n"))
}

// RegisterHandler requests the singleton http.Handler from this component
type RegisterHandler struct {
	// Name is the human-readable name for this handler.  It will appear in the
	// table of contents.  If this is empty, the handler will not be included
	// in the table of contents.
	Name string

	// Pattern is the pattern to which this handler should be attached.  See
	// https://pkg.go.dev/net/http#ServeMux.
	Pattern string

	// Handler is the handler to be registered.
	Handler http.Handler
}

// Serve requests the singleton http.Handler from this component
type Serve struct {
	// Port is the port on which to run the server
	Port int
}

// HandlerRequest requests the singleton http.Handler from this component
type HandlerRequest struct{}

// HandlerResponse returns the singleton http.Handler from this component
type HandlerResponse struct {
	Handler http.Handler
}
