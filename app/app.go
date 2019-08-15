package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	client "github.com/taufikardiyan28/chat/client"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		/*token := r.Header.Get("Sec-WebSocket-Protocol")
		fmt.Println("token", token)
		if strings.TrimSpace(token) != "tes" {
			fmt.Println(token)
			return false
		}*/
		return true
	},
}

type (
	Server struct {
		Config  *Config
		Clients map[string]*client.Connection
	}

	Config struct {
		Port int `yaml:"port"`
	}
)

//Start the application
func (a *Server) Start() {
	a.Clients = make(map[string]*client.Connection)

	http.HandleFunc("/ws", a.handleWSConnections)
	log.Printf("http server started on :%d", a.Config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", a.Config.Port), nil)
}

//Handle websocket connection
func (a *Server) handleWSConnections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ERROR", err)
		//ws.WriteJSON(map[string]interface{}{"status": 1, "msg": "Unauthorized"})
		//ws.Close()
		return
	}
	//defer ws.Close()

	username := r.URL.Query().Get("username")
	key := strings.TrimSpace(username)

	if a.isKeyExists(key) {
		log.Printf("Key %s is exists", key)
		ws.Close()
		return
	}

	c_info := client.UserInfo{
		Key:      key,
		UserName: username,
	}
	c := &client.Connection{
		Conn:     ws,
		UserInfo: c_info,
		Public:   &a.Clients,
	}
	a.Clients[key] = c
	go c.Listen()
}

func (a *Server) isKeyExists(key string) bool {
	for clientKey := range a.Clients {
		if clientKey == key {
			return true
		}
	}
	return false
}
