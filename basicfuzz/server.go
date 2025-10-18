package basicfuzz

import(
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)


func ServeWs(){

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