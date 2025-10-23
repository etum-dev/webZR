package scan

import (
	"net/url"

	"github.com/gorilla/websocket"
	//"bytes"
	"bufio"
	"fmt"
	"os"
)

func CheckJS(url string) {
	// Parse sites js to find if it has ws estab.
	// Should be delegated to tools like katana and idk rg
}

func ScanEndpoint( /*endpointsFile string*/ url string) {
	// just example one for now. add as opt later.
	fmt.Println("now scanning: ", url)
	wordlist, err := os.Open("ws-endpoints.txt")

	if err != nil {
		fmt.Println("buh", err)
	}
	defer wordlist.Close()
	scanner := bufio.NewScanner(wordlist)
	for scanner.Scan() {
		de := url + scanner.Text()
		fmt.Println(de)
		SendConnRequest(de)

	}
}

func ScanSubdomain(url string) {

}

func SendConnRequest(domain string) bool {
	/* Attempts to connect via ws and wss. Should flag for:
	// - If connection possible
	// - If connection possible and allows for ws://*/
	schemes := []string{"wss", "ws"}
	for _, scheme := range schemes {
		// Parse to separate host and path properly
		parsedURL, err := url.Parse("https://" + domain)
		if err != nil {
			fmt.Printf("Failed to parse domain %s: %v\n", domain, err)
			continue
		}

		u := url.URL{Scheme: scheme, Host: parsedURL.Host, Path: parsedURL.Path}
		if u.Path == "" {
			u.Path = "/"
		}
		wsUrl := u.String()

		dialer := websocket.Dialer{
			// add timeout
		}
		conn, resp, err := dialer.Dial(wsUrl, nil)

		if resp.StatusCode != 101 {
			continue
		}
		if err != nil {
			fmt.Printf("%s:  %v \n", wsUrl, err)
			//return false
		}
		defer conn.Close()

		fmt.Printf("yay conn %s (Status %d)\n", wsUrl, resp.StatusCode)

		/* flagging if we can conn with insecure ws:
		if scheme == "ws" {
			fmt.Printf("(`L_` )!! Insecure WS GET しました")
		} else {
			fmt.Printf("wss only, good boy\n")
			continue
		} */

	}
	return true
}

func CorsITaket(url string, ownserver string) {
	// Check if it validates origin'
	/* if (cors not validated) {
		send request with own server value
	}*/

	// ideally, i make this pipeable from eg interactsh.
	// How do tools like httpx take such inputs?

}
