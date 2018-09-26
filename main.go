package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sjsafranek/ligneous"
)

const (
	DEFAULT_HTTP_PORT = 8000
)

var (
	port    int    = DEFAULT_HTTP_PORT
	logger         = ligneous.NewLogger()
	version string = "0.0.1"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	flag.IntVar(&port, "p", DEFAULT_HTTP_PORT, "http server port")
	flag.Parse()

	if 1 < len(os.Args) {
		if "createuser" == os.Args[1] {
			email := os.Args[2]
			password := os.Args[3]
			CreateUser(email, password)
			fmt.Println("User created")
			os.Exit(0)
		} else if "version" == os.Args[1] {
			fmt.Println(version)
			os.Exit(0)
		}
	}
}

func main() {

	router := mux.NewRouter()

	// Create a simple file server
	fs := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	router.HandleFunc("/login", LoginHandler).Methods("GET", "POST")
	router.HandleFunc("/logout", LogoutHandler).Methods("GET")
	router.HandleFunc("/chat", ChatHandler).Methods("GET")
	router.HandleFunc("/ws/{chatroom}", WebSocketHandler)
	http.Handle("/", router)

	router.Use(LoggingMiddleWare, SetHeadersMiddleWare)

	// source: http://patorjk.com/software/taag/#p=display&f=Slant&t=CryptoChat
	fmt.Println(`
   ______                 __        ________          __
  / ____/______  ______  / /_____  / ____/ /_  ____ _/ /_
 / /   / ___/ / / / __ \/ __/ __ \/ /   / __ \/ __ '/ __/
/ /___/ /  / /_/ / /_/ / /_/ /_/ / /___/ / / / /_/ / /_
\____/_/   \__, / .___/\__/\____/\____/_/ /_/\__,_/\__/
          /____/_/
	`)

	// Start the server on localhost port 8000 and log any errors
	logger.Infof("Magic happens on port :%v", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		logger.Critical("ListenAndServe: ", err)
	}

	// google-chrome --app=http://localhost:8000/login --window-size=600,400
}
