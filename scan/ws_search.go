package scan

import (
	"net/url"

	"github.com/gorilla/websocket"
	//"bytes"
	"bufio"
	"fmt"
	"os"

	"github.com/etum-dev/WebZR/basicfuzz"
)

// This file will look for WS and then maybe send to the fuzz things
// Also checks for simple vuln potential, such as misconfigured cors

type Target struct {
	domain       string
	baselineBody []byte
	abnormalBody []byte
}

type CheckResult struct {
	Domain     string `json:"domain"`
	StatusCode int    `json:"status_code"`
	Successful bool   `json:"successful"`
}

func ConnToWS(url string) {
	basicfuzz.NewWSFuzzer(url)
}

func CheckJS(url string) {
	// Parse sites js to find if it has ws estab.
	// Should be delegated to tools like katana and idk rg
}

func ScanEndpoint(endpointsFile string, url string) {
	// just example one for now. add as opt later.
	wordlist, err := os.Open("")
	if err != nil {
		fmt.Println("buh")
	}
	defer wordlist.Close()
	scanner := bufio.NewScanner(wordlist)
	for scanner.Scan() {
		fmt.Println("asdasdasda %s\n", scanner.Text())
	}
}

func ScanSubdomain(url string) {

}

func SendConnRequest(domain string) {
	/* Attempts to connect via ws and wss. Should flag for:
	// - If connection possible
	// - If connection possible and allows for ws://*/
	schemes := []string{"wss", "ws"}
	for _, scheme := range schemes {
		u := url.URL{Scheme: scheme, Host: domain, Path: "/"} //path maybe looped on ws-endpoints.txt idk
		wsUrl := u.String()

		dialer := websocket.Dialer{
			// add timeout
		}
		conn, resp, err := dialer.Dial(wsUrl, nil)

		if err != nil {
			fmt.Printf("%s:  %v \n", wsUrl, err)
			continue
		}
		defer conn.Close()

		fmt.Printf("yay conn %s (Status %d)\n", wsUrl, resp.StatusCode)

		// flagging if we can conn with insecure ws://
		if scheme == "ws" {
			fmt.Printf("(`L_` )!! Insecure WS GET しました")
		}
		return
	}

}

func CorsITaket(url string, ownserver string) {
	// Check if it validates origin'
	/* if (cors not validated) {
		send request with own server value
	}*/

	// ideally, i make this pipeable from eg interactsh.
	// How do tools like httpx take such inputs?

}
