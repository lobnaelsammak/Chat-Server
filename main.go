package main
import (
	"log" //log library 3ashn n-keep track of events happening 
	"net/http" //library responsible for http - tcp
	"github.com/gorilla/websocket" //socket api choosed 
)
//LINK EACH ATTRIBUTE TO JSON 
type Message struct{
	Username string `json:"username"` 
	Email	string 	`json:"email"`
	Message	string 	`json:"message"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main(){
	//ServeFile replies to the request with the contents of the named file or directory. (linking html,js files)
	fileserver := http.FileServer(http.Dir("../public")) 
	//Handle registers the handler for the given pattern in the DefaultServeMux.
	http.Handle("/",fileserver)

	//HandleFunc registers the handler function for the given pattern in the DefaultServeMux.
	http.HandleFunc("/ws",connectionHandler)

	go messageHandler() //start threading 


	log.Println("http server started on : 8081 ")
	//ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections
	err := http.ListenAndServe(":8081", nil) 
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}


func connectionHandler(w http.ResponseWriter, r *http.Request){
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close() //to make sure we close the connection after request 

	clients[socket] = true //register new clients 

	for {
		var msg Message
		err := socket.ReadJSON(&msg) //map new message to JSON object 
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, socket) //clear objects to use next time
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}

}

func messageHandler() {
	for {
		//get message from channel
		msg := <-broadcast
		//send message to all connected clients
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

