
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
var (
	//websockt对象
	upgrader = websocket.Upgrader{}
	// 保存所有客户端连接,key是用户名	
	clients = make(map[string]*websocket.Conn) 
	// 用户上线/下线的消息通道
	onlineChan = make(chan Message)
	// 用户添加好友请求的消息通道
	addFriendChan  = make(chan Message)
	// 用户删除好友请求的消息通道
	delFriendChan  = make(chan Message)
	// 用户聊天消息通道
	chatChan = make(chan Message)
)
const (

)


var chans = make(chan Message)           // broadcast channel
//var test=make(map[int]string)

type Message struct {
	Src    string `json:"src"`
	Dst string `json:"dst"`
	Message  string `json:"message"`
}
// type MyConn struct{
// 	userName string
// 	ws *websocket.Conn
// }
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
		u.POST("/login",user.Login)
		u.GET("/yy/:name",func(c *gin.Context){
	
			username, _ := c.Cookie("username")
			c.JSON(200, gin.H{
				"success" : true,
				"msg": "ok",
				"data": "hahaha"+username,
			})
			return
		})
		u.GET("/ws",func(c *gin.Context){
				
			    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
				ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
				
				//myConn := MyConn{c.Query("a"),ws}
				if err != nil {
					log.Fatal(err)
				}
				clients[c.Query("a")] = ws
				fmt.Println(c.Query("a")+"来了.......................")
				// Make sure we close the connection when the function returns
				defer func(){
					fmt.Println(c.Query("a")+".....................离开了")
					delete(clients, c.Query("a"))
					ws.Close()
				}()
	
				for {
					var msg Message            // Read in a new message as JSON and map it to a Message object
					err := ws.ReadJSON(&msg)
					//fmt.Println(msg.Src,msg.Dst,msg.Message)
					if err != nil {
						log.Printf("error: %v", err)
						delete(clients, msg.Src)
						break
					}
					
					// msg = Message{
					// 	Email:"haha ",
					// }
					// Send the newly received message to the broadcast channel
					chans <- msg        
				}
			})

	}
	fmt.Println("----------------------------------")
	router.Run(":8000") // listen and serve on 0.0.0.0:8000
}
func handleMessages() {
	for {
		// Grab the next message from the broadcast channel

		msg := <-chans
		fmt.Println(msg)
		//Send it out to every client that is currently connected
		clients[msg.Dst].WriteJSON(msg)
		// for username,clientWs := range clients {
		// 	fmt.Println(username,msg.Src,msg.Dst,msg.Message)
		// 	if(msg.Dst == username){
		// 		//err :=clientWs.WriteMessage(1,[]byte(msg.Message))
		// 		err := clientWs.WriteJSON(msg)
		// 		if err != nil {
		// 			log.Printf("error: %v", err)
		// 			clientWs.Close()
		// 			delete(clients, username)
		// 		}
		// 	}
		// }
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