package server

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"time"

	"github.com/go-zoo/bone"
	"github.com/msales/pkg/log"
)

// Application represents the main application.
type Application interface {
	// Get the next message from the queue.
	GetNextMessage(ctx context.Context, topic string) ([]byte, error)
	// GetNextBatch gets the next count messages from the queue.
	GetNextBatch(ctx context.Context, topic string, count int) ([][]byte, error)
	// IsHealthy checks the health of the Application.
	IsHealthy() error
}

// Server represents a http server handler.
type Server struct {
	app Application
	mux *bone.Mux
}

// New creates a new Server instance.
func New(app Application) *Server {
	s := &Server{
		app: app,
		mux: bone.New(),
	}

	s.mux.GetFunc("/health", s.HealthHandler)
	s.mux.GetFunc("/:topic", s.TopicHandler)
	s.mux.GetFunc("/batch/:topic", s.BatchTopicHandler)

	s.mux.NotFound(NotFoundHandler())

	return s
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// TopicHandler handles requests to for messages.
func (s *Server) TopicHandler(w http.ResponseWriter, r *http.Request) {
	topic := bone.GetValue(r, "topic")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancel()

	b, err := s.app.GetNextMessage(ctx, topic)
	if err != nil {
		log.Error(r.Context(), "server: could not get message", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if b == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if len(b) == 0 {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(b)
}

// TopicHandler handles requests to for messages.
func (s *Server) BatchTopicHandler(w http.ResponseWriter, r *http.Request) {
	topic := bone.GetValue(r, "topic")
	size := 100
	ssize := r.URL.Query().Get("size")
	if i64, err := strconv.ParseInt(ssize, 10, 64); err == nil {
		size = int(i64)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancel()

	msgs, err := s.app.GetNextBatch(ctx, topic, size)
	if err != nil {
		log.Error(r.Context(), "server: could not get message", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if msgs == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if len(msgs) == 0 {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}

	mw := multipart.NewWriter(w)
	w.Header().Set("Content-Type", "multipart/mixed; boundary="+mw.Boundary())
	for _, b := range msgs {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Type", "application/octet-stream")
		pw, err := mw.CreatePart(h)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		pw.Write(b)
	}
	mw.Close()
}

// HealthHandler handles health requests.
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.app.IsHealthy(); err != nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// NotFoundHandler returns a 404.
func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
}
