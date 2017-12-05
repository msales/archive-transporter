package server_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msales/transporter/server"
	"github.com/stretchr/testify/assert"
)

type testApp struct {
	nextMessage func() ([]byte, error)
	isHealthy   func() error
}

func (a testApp) GetNextMessage(ctx context.Context, topic string) ([]byte, error) {
	return a.nextMessage()
}

func (a testApp) IsHealthy() error {
	return a.isHealthy()
}

func TestServer_TopicHandler(t *testing.T) {
	tests := []struct {
		url  string
		data []byte
		err  error
		code int
	}{
		{"/foobar", []byte{'a'}, nil, 200},
		{"/foobar", []byte{}, nil, 204},
		{"/foobar", nil, nil, 404},
		{"/foobar", nil, errors.New("test error"), 500},
	}

	for _, tt := range tests {
		app := testApp{
			nextMessage: func() ([]byte, error) {
				return tt.data, tt.err
			},
		}
		srv := server.New(app)

		r := httptest.NewRequest("GET", tt.url, nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, r)

		assert.Equal(t, tt.code, w.Code)
	}
}

func TestServer_HealthHandler(t *testing.T) {
	tests := []struct {
		err  error
		code int
	}{
		{nil, http.StatusOK},
		{errors.New(""), http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		app := testApp{
			isHealthy: func() error {
				return tt.err
			},
		}
		srv := server.New(app)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/health", nil)
		srv.ServeHTTP(w, req)

		assert.Equal(t, tt.code, w.Code)
	}
}

func TestNotFoundHandler(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	server.NotFoundHandler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
