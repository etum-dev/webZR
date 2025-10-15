package main 

import(
	"github.com/etum-dev/WebZR/scan"

)

func main(){
	scan.AuthShodan()
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