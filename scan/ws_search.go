package scan

import (
	"time"

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
	- If connection possible
	- allows for ws://
	*/

	schemes := []string{"wss", "ws"}
	for _, scheme := range schemes {
		// Parse to separate host and path properly
		wsUrl := scheme + "://" + domain

		dialer := websocket.Dialer{
			HandshakeTimeout: 5 * time.Second,
		}
		conn, resp, err := dialer.Dial(wsUrl, nil)

		if err != nil {
			// Enable this line if verbose mode
			//fmt.Printf("wtf %s:  %v \n", wsUrl, err)
			continue 
		}
		if resp == nil {
			fmt.Printf("rip %s, no response\n", wsUrl)
			continue // Skip to next scheme
		}

		// only defer close if conn is not nil
		if conn != nil {
			defer conn.Close()
		}
		if resp.StatusCode == 101 {
			fmt.Printf("yay conn %s (Status %d)\n", wsUrl, resp.StatusCode)
			if scheme == "ws" {
				fmt.Printf("(`L_` )!! Insecure WS GET しました\n")
			}
		} else {
			fmt.Printf("No 101, no bitches %s : %d\n", wsUrl, resp.StatusCode)
		}

	}
	return true
}

//

func CorsITaket(url string, ownserver string) {
	// Check if it validates origin'
	/* if (cors not validated) {
		send request with own server value
	}*/

}
