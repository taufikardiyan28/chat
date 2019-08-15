package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	client "github.com/taufikardiyan28/chat/client"
	room "github.com/taufikardiyan28/chat/room"
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

	// initiate room list
	r := make(map[string]*room.ChatRoom)
	room.Rooms = &r

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
	id := strings.TrimSpace(username)

	if a.isClientIDExists(id) {
		log.Printf("Key %s is exists", id)
		ws.Close()
		return
	}

	c_info := client.UserInfo{
		ID:       id,
		UserName: username,
	}
	c := &client.Connection{
		Conn:     ws,
		UserInfo: c_info,
		Public:   &a.Clients,
	}
	a.Clients[id] = c
	go c.Listen()
}

func (a *Server) isClientIDExists(id string) bool {
	for client_id := range a.Clients {
		if client_id == id {
			return true
		}
	}
	return false
}
