// https://github.com/PalindromeLabs/STEWS/blob/main/discovery/README.md
package scan

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type VulnTestResult struct {
	URL	string
	Origin string
	Success bool
	Err string
}

// This function checks if target validates origin
func CorsITaket(wsurl string) {

	// if WebsocketAddr defined, use this server.
	origins := []string{
		"https://example.com", // invalid url
		wsurl,                 // TODO: parse url, if ws(s), replace with http(s). utils.AppendProto
	}

	for _, or := range origins {
		header := http.Header{}
		header.Set("Origin", or)
		r1 := websocket.DefaultDialer
		r1.HandshakeTimeout = 5 * time.Second

		conn, resp, err := r1.Dial(wsurl, header)
		if err != nil {
			if resp != nil {
				fmt.Printf("Origin=%q -> REJECTED (http %d)\n", or, resp.StatusCode)
			} else {
				fmt.Printf("Origin=%q -> ERROR (no response): %v\n", or, err)
			}
			continue
		}
		fmt.Printf("ヽ(●´ｗ｀○)ﾉ Connection successful with Origin: %q", or)
		defer conn.Close()
		//TODO check if origin is not the same as wsurl, if so, flag
	}
}
