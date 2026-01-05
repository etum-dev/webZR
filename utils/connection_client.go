package utils

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

func NewWebSocketDialer(timeout time.Duration) *websocket.Dialer {
	return &websocket.Dialer{HandshakeTimeout: timeout}
}
