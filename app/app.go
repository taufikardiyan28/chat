package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/taufikardiyan28/chat/db"
	"github.com/taufikardiyan28/chat/interfaces"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	client "github.com/taufikardiyan28/chat/client"
	"github.com/taufikardiyan28/chat/helper"
	room "github.com/taufikardiyan28/chat/room"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type (
	Server struct {
		Config  *helper.Configuration
		Clients map[string]*client.Connection
		DB      interfaces.Database
	}
)

//Start the application
func (a *Server) Start() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}\r\n    method=${method}, uri=${uri}, status=${status} remote_ip=${remote_ip}\n",
	}))
	e.Use(middleware.Recover())
	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	a.Clients = make(map[string]*client.Connection)

	// initiate room list
	r := make(map[string]*room.ChatRoom)
	room.Rooms = &r

	dbcon, err := db.NewConnection(a.Config)

	if err != nil {
		panic(err)
	}

	if err := dbcon.Ping(); err != nil {
		panic(err)
	}

	a.DB = dbcon

	e.HideBanner = true

	e.Static("/", "public")
	e.File("/demo", "web/index.html")

	e.GET("/ws", a.handleWSConnections)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", a.Config.Port)))
}

//Handle websocket connection
func (a *Server) handleWSConnections(c echo.Context) error {
	userId := c.QueryParam("id")
	id := strings.TrimSpace(userId)

	a.closePreviousConnection(id)

	c_info, err := a.DB.GetUserInfo(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(401, map[string]interface{}{"status": "error", "msg": err})
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println(err)
		return c.JSON(401, map[string]interface{}{"status": "error", "msg": err})
	}
	defer ws.Close()

	clientCon := &client.Connection{
		Conn:        ws,
		User:        c_info,
		OnlineUsers: &a.Clients,
		DB:          a.DB,
	}
	a.Clients[id] = clientCon
	clientCon.Start()
	return err
}

func (a *Server) closePreviousConnection(id string) bool {
	for client_id := range a.Clients {
		if client_id == id {
			a.Clients[client_id].Close()
			delete(a.Clients, client_id)
			return true
		}
	}
	return false
}
