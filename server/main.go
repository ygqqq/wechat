
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
	onlineChan = make(chan Message,2)
	// 用户添加好友请求的消息通道
	addFriendChan  = make(chan Message,2)
	// 用户删除好友请求的消息通道
	delFriendChan  = make(chan Message,2)
	// 用户聊天消息通道
	chatChan = make(chan Message,2)
)
const (
	//MessageType
	Error		    = 0 //错误消息
	OnlineRemind	= 1	//上线提醒
	OfflineRemind   = 2 //下线提醒 
	AddFriendReq	= 3 //添加好友请求
	AgreeAdd		= 4 //同意好友请求
	DisAgreeAdd 	= 5 //拒绝好友请求
	
)

type Message struct {
	Src    string `json:"src"`
	Dst string `json:"dst"`
	Message  string `json:"message"`
	MessageType int `json:"messagetype"`
}


func main() {
	// 处理添加好友请求
	go handleFriendMessages()
	// 处理用户上下线提醒消息
	go handleConnMessages()
	// 
	go handlerWs()

	router := gin.Default()

	//用户操作相关路由
	u := router.Group("/user")
	{
		//用户注册
		u.POST("/register", user.Register)
		//用户登陆
		u.POST("/login",user.Login)
	}
	//websocket入口
	router.GET("/ws",wsConnHandler)

	router.Run(":8000") // listen and serve on 0.0.0.0:8000
}
func wsConnHandler(c *gin.Context){
	//建立ws连接
	
	//ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	//获取客户端传来的用户名cookie，用于标示是哪个客户端
	userName := c.Query("a")
	//clients[userName] = ws

	if _,ok := clients[userName]; !ok{
		fmt.Println("新建立的链接")
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		clients[userName] = ws
		if err != nil {
			log.Fatal(err)
		}
	}else{
		fmt.Println("老的链接")
	}

	//ws连接之后，向onlineChan通道传递消息，通知所有在线好友
	// onlineChan <- Message{
	// 	Src: userName,
	// 	MessageType: OnlineRemind,
	// 	Message: userName+"上线啦",
	// }
	//ws断开连接
	defer func(){
		// onlineChan <- Message{
		// 	Src: userName,
		// 	MessageType: OnlineRemind,
		// 	Message: userName+"下线啦",
		// }
		fmt.Println(userName+"断开...............")
		if _,ok := clients[userName]; ok{
			clients[userName].Close()
			delete(clients, userName)

		}
	}()

	//for {
		//当有客户端发送消息过来，判断消息类型，分配给相应的消息处理管道
		var msg Message            
		err := clients[userName].ReadJSON(&msg)
		if err != nil{
			fmt.Println("错了")
			clients[userName].Close()
			delete(clients, userName)
		}
		// switch msg.MessageType{
		// //添加、同意、拒绝好友请求	
		// case AddFriendReq,AgreeAdd,DisAgreeAdd:
		// 	fmt.Println("有人要添加了")
		// 	addFriendChan <- msg
		// case 8:
		// 	fmt.Println("断开了")
		// }     
	//}
}
func handlerWs(){

}
// 处理好友相关请求
func handleFriendMessages() {
	for {
		msg := <- addFriendChan
		fmt.Println("这里处理添加请求")
		//首先要判断消息类型
		switch msg.MessageType{
		//如果是请求加好友，还要判断用户是否在线，先只做成只有在线才能加把	
		case AddFriendReq:
			dstUser,err := user.GetUserByName(msg.Dst)
			if err != nil {
				msg.Message = "用户不存在"
				clients[msg.Src].WriteJSON(msg)
				break
			}
			if dstUser.IsMyFriend(msg.Src) {
				msg.Message = "不要重复添加好友"
				fmt.Println("不要重复添加好友")
				clients[msg.Src].WriteJSON(msg)
				break
			}
			//如果目标用户存在　并且在线，则给目标用户推送加好友请求
			if ws,ok := clients[msg.Dst]; ok {
				ws.WriteJSON(msg)
				break
			}else{
				//目标用户不在线，给发起请求的用户推送消息
				msg.Message = "用户不在线"
				clients[msg.Src].WriteJSON(msg)
				break
			}
		case AgreeAdd:
			//如果是同意好友请求
			dstUser,_ := user.GetUserByName(msg.Dst)
			srcUser,_ := user.GetUserByName(msg.Src)
			//判断两者之前是否已经是好友
			if !dstUser.IsMyFriend(msg.Src) && !srcUser.IsMyFriend(msg.Dst){
				//将两者的好友列表append对方的用户名
				dstUser.AddOrDelFriendByName(msg.Src)
				srcUser.AddOrDelFriendByName(msg.Dst)
				msg.Message = "添加成功"
				clients[msg.Src].WriteJSON(msg)
				break
			}else{
				msg.Message = "重复添加"
				clients[msg.Src].WriteJSON(msg)
				break
			}
	
			//如果对方也在线，给对方页发送推送
			// if ws,ok := clients[msg.Dst]; ok {
			// 	ws.WriteJSON(msg)
			// 	return
			// }
		}
	}
}
func handleConnMessages(){
	for{
		msg := <- onlineChan
		for username,clientWs := range clients {
			if err := clientWs.WriteJSON(msg);err != nil {
				clientWs.Close()
				delete(clients, username)
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