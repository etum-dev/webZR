package scan

import(
	//"github.com/gorilla/websocket"
)

type CheckResult struct {
	Domain	string `json:"domain"`
	StatusCode	int	`json:"status_code"`
	Successful	bool	`json:"successful"`
}

func CheckHeaders(url string) error {
	// TODO: most likely has to tweak these 4 speed
	// ... Do I even need this? Can I connect to WS directly?
	/*
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			MaxIdleConns:	100,
			MaxIdleConnsPerHost: 10,
		},

	} */
	
}

func CheckJS(url string){
	// Parse sites js to find if it has ws estab.
	// Should be delegated to tools like katana and idk rg
}

func ScanEndpoint(endpointsFile string, url string){

}

func ScanSubdomain(url string){

}


func SendConnRequest(url string){
	// Attempts to connect via ws and wss. Should flag for: 
	// - If connection possible
	// - If connection possible and allows for ws://



}

