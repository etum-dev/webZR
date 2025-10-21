package scan

import(
	//"github.com/gorilla/websocket"
	//"bytes"
	"bufio"
	"os"
	"fmt"
	"net/http"
	"github.com/etum-dev/WebZR/basicfuzz"
)

// This file will look for WS and then maybe send to the fuzz things
// Also checks for simple vuln potential, such as misconfigured cors

type Target struct {
	domain string
	baselineBody []byte
	abnormalBody []byte
}

type CheckResult struct {
	Domain	string `json:"domain"`
	StatusCode	int	`json:"status_code"`
	Successful	bool	`json:"successful"`
}

func ConnToWS(url string){
	basicfuzz.NewWSFuzzer(url)
}


func CheckJS(url string){
	// Parse sites js to find if it has ws estab.
	// Should be delegated to tools like katana and idk rg
}

func ScanEndpoint(endpointsFile string, url string){
	// just example one for now. add as opt later. 
	wordlist, err := os.Open("")
	if err != nil {
		fmt.Println("buh")
	}
	defer wordlist.Close()
	scanner := bufio.NewScanner(wordlist)
	for scanner.Scan(){
		fmt.Println("asdasdasda %s\n", scanner.Text())
	}
}

func ScanSubdomain(url string){

}


func SendConnRequest(url string){
	// Attempts to connect via ws and wss. Should flag for: 
	// - If connection possible
	// - If connection possible and allows for ws://
	// GET (endpoint) -> 101 resp
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} 
	fmt.Println(res)

}


func CorsITaket(url string, ownserver string){
	// Check if it validates origin'
	/* if (cors not validated) {
		send request with own server value
	}*/

	// ideally, i make this pipeable from eg interactsh. 
	// How do tools like httpx take such inputs?
	
}

