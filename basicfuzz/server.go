package basicfuzz

import(
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	
)


func ServeWs(){
	//TODO: add wss https://tillitsdone.com/blogs/secure-websockets-in-golang/
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		var upgrader = websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
		}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		} 
		log.Println("ok",ws) 
	} )


	log.Fatal(http.ListenAndServe(":8000", nil)) 
}

func ServeWss(){
	certFile := "server.crt"
	keyFile := "server.key"
	http.HandleFunc("/secwss", func(w http.ResponseWriter, r *http.Request) {
		// fuck DRY all my homies hate DRY
		var upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		} 
		log.Println("ok ok",ws) 


	})
	log.Println("secure wss started")
	log.Fatal(http.ListenAndServeTLS(":8443", certFile, keyFile, nil));

}