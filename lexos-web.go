package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/Jibble330/lexos"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func socket(c *gin.Context) {
    log.Println("WebSocket connection")
    ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer ws.Close()

    _, isbn, err := ws.ReadMessage()
    if err != nil {
        log.Println(err)
        return
    }

    lexile, atos, ar, err := lexos.Get(string(isbn))
    if err != nil {
        var message []byte
        if err == lexos.InvalidISBN {
            message = []byte("error:0")
        } else {
            log.Println(err.Error())
            message = []byte("error:" + err.Error())
        }
        ws.WriteMessage(websocket.TextMessage, message)
        return
    }
    msg := make(map[string]string)
    msg["lexile"] = fmt.Sprint(lexile)
    msg["atos"] = fmt.Sprint(atos)
    msg["ar"] = fmt.Sprint(ar)

    ws.WriteJSON(msg)
}

func home(c *gin.Context) {
    c.File("front.html")
}

func main() {
    lexos.Install()
    gin.SetMode(gin.ReleaseMode)

    r := gin.Default()
    r.GET("/", home)
    r.GET("/ws", socket)
    log.Println("Starting server on http://localhost:80")
    err := r.Run(":80")
    if err != nil {
        log.Fatal(err)
    }
}
