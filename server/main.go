
package main

import (
	"./user"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	//"github.com/garyburd/redigo/redis"
	//"./utils"
)
var clients = make(map[*websocket.Conn]string) // connected clients
var broadcast = make(chan Message)           // broadcast channel
var upgrader = websocket.Upgrader{}
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
	}
type MyConn struct{
	userName string
	ws *websocket.Conn
}
type Y struct{
	Name string	`json:"name"`
	Age int		`json:"age"`
}
func main() {
	go handleMessages()
	router := gin.Default()

	//用户操作相关路由
	u := router.Group("/user")
	{
		
		//用户注册
		u.POST("/register", user.Register)
		u.GET("/yy",func(c *gin.Context){
		
			c.JSON(200, gin.H{
				"success" : true,
				"msg": "ok",
				"data": "hahaha",
			})
			return
		})
		u.GET("/ws",func(c *gin.Context){
				
			    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
				ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
				myConn := MyConn{c.Query("a"),ws}
				if err != nil {
						log.Fatal(err)
				}
				// Make sure we close the connection when the function returns
				defer func(){
					fmt.Println(myConn.userName+".....................离开了")
					ws.Close()
				}()
				fmt.Println(c.Query("a")+".....................来了")
				clients[ws] = c.Query("a")
				for {
					var msg Message            // Read in a new message as JSON and map it to a Message object
					err := ws.ReadJSON(&msg)
					if err != nil {
							log.Printf("error: %v", err)
							delete(clients, ws)
							break
					}
					// msg = Message{
					// 	Email:"haha SB",
					// }
					// Send the newly received message to the broadcast channel
					broadcast <- msg        
				}
			})

	}
	fmt.Println("----------------------------------")
	router.Run(":8000") // listen and serve on 0.0.0.0:8000
}
func handleMessages() {
	for {
			// Grab the next message from the broadcast channel
			msg := <-broadcast
			// Send it out to every client that is currently connected
			for client := range clients {
				log.Printf("client: %v", client)
				err := client.WriteJSON(msg)

				if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
				}
			}
	}
}
func Middleware(c *gin.Context) {
	id := c.Query("id")
	if id == "1" {
		fmt.Println("this is a middleware!")
		
	}else{
		c.JSON(200, gin.H{
			"message": "SB !=1",
		})
		c.Abort()
	}
    
}