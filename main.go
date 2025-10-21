package main 

import(
	"github.com/etum-dev/WebZR/scan"
	"github.com/etum-dev/WebZR/basicfuzz"
	"fmt"
	"flag"

)


func main(){

	isTest := flag.Bool("test",false,"Use the local test server")
	fuzzOpt := flag.Int("fuzz",1,"Fuzzing option. 1 for basic, 2 for custom header, 3 for mutation")
	proxyOpt := flag.String("proxy")
	//customheader := flag.String("BugBounty: xyz")
	flag.Parse()
	if *isTest == true {
		basicfuzz.ServeWs()
	}
	if *fuzzOpt == 1 {
		fmt.Println("Doing a simple fuzz.")
		basicfuzz.SimpleFuzz()
	} else {
		fmt.Println("Not fuzzing")
	}



	//scan.AuthShodan()
	scan.SendConnRequest("https://stream.binance.com/stream")
}
// need more methods to identify websockets

/*

https://github.com/projectdiscovery/nuclei-templates/issues/11243

SOCKET.IO:

<script src="/socket.io/socket.io.js"></script>
<script type="module">
  import { io } from "https://cdn.socket.io/4.8.1/socket.io.esm.min.js";
</script>
*/