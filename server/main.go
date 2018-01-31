
package main

import (
	"./user"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gorilla/websocket"
	"encoding/json"
	//"github.com/garyburd/redigo/redis"
	"./utils"
	"./kafka"
)
var (
	//websockt对象
	upgrader = websocket.Upgrader{}
	// 保存所有客户端连接,key是用户名	
	clients = make(map[string]*websocket.Conn) 
	// 用户上线/下线的消息通道
	onlineChan = make(chan Message,10)
	// 用户添加好友请求的消息通道
	addFriendChan  = make(chan Message,10)
	// 用户删除好友请求的消息通道
	delFriendChan  = make(chan Message,10)
	// 用户聊天消息通道
	chatChan = make(chan Message,10)
)
const (
	//MessageType
	ErrorMsg		= 0 //错误消息
	OnlineRemind	= 1	//上线提醒
	OfflineRemind   = 2 //下线提醒 
	AddFriendReq	= 3 //添加好友请求
	AgreeAdd		= 4 //同意好友请求
	DisAgreeAdd 	= 5 //拒绝好友请求
	
	NormalMsg		= 10 //普通通知消息
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
	// kafka消费者
	go kafkaConsumer()

	router := gin.Default()

	//用户操作相关路由
	u := router.Group("/user")
	{
		//用户注册
		u.POST("/register", user.Register)
		//用户登陆
		u.POST("/login",user.Login)
		//获取所有好友
		u.GET("/friends/:name",func(c *gin.Context){
			name := c.Param("name")
			ur,_ := user.GetUserByName(name)
			urs := ur.GetAllFriends()
			strUser,_ := json.Marshal(urs)
			c.JSON(200, gin.H{
				"success" : true,
				"msg": string(strUser),
			})
			return
		})
	}
	//websocket入口
	router.GET("/ws",wsConnHandler)

	router.Run(":8000") // listen and serve on 0.0.0.0:8000
}
func wsConnHandler(c *gin.Context){
	//获取客户端传来的用户名cookie，用于标示是哪个客户端
	userName := c.Query("a")
	//建立ws连接
	if _,ok := clients[userName]; !ok{
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
		clients[userName] = ws
	}
	//ws连接之后，向onlineChan通道传递消息，通知所有在线好友
	onlineChan <- Message{
		Src: userName,
		MessageType: OnlineRemind,
		Message: userName+"上线啦",
	}
	//ws断开连接
	defer func(){
		onlineChan <- Message{
			Src: userName,
			MessageType: OfflineRemind,
			Message: userName+"下线啦",
		}
		if _,ok := clients[userName]; ok{
			clients[userName].Close()
			delete(clients, userName)
		}
	}()

	for {
		//当有客户端发送消息过来，判断消息类型，分配给相应的消息处理管道
		var msg Message            
		err := clients[userName].ReadJSON(&msg)
		if err != nil{
			clients[userName].Close()
			delete(clients, userName)
			break
		}
		switch msg.MessageType{
		//添加、同意、拒绝好友请求	
		case AddFriendReq,AgreeAdd,DisAgreeAdd:
			addFriendChan <- msg
		}     
	}
}

// 处理好友相关请求
func handleFriendMessages() {
	for {
		msg := <- addFriendChan
		//首先要判断消息类型
		switch msg.MessageType{
		//如果是请求加好友，还要判断用户是否在线，先只做成只有在线才能加把	
		case AddFriendReq:
			dstUser,err := user.GetUserByName(msg.Dst)
			if err != nil {
				msg.Message = "用户不存在"
				msg.MessageType = ErrorMsg
				clients[msg.Src].WriteJSON(msg)
				break
			}
			if dstUser.IsMyFriend(msg.Src) {
				msg.Message = "不要重复添加好友"
				msg.MessageType = ErrorMsg
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
				msg.MessageType = ErrorMsg
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
				dstUser.AddFriendByName(msg.Src)
				srcUser.AddFriendByName(msg.Dst)
				msg.Message = "添加成功"
				msg.MessageType = NormalMsg
				clients[msg.Src].WriteJSON(msg)
				if ws,ok := clients[msg.Dst]; ok {
					ws.WriteJSON(msg)
				}
				break
			}else{
				msg.Message = "重复添加"
				msg.MessageType = ErrorMsg
				clients[msg.Src].WriteJSON(msg)
				break
			}
		}
	}
}
//处理好友上线下线的消息
func handleConnMessages(){
	for{
		msg := <- onlineChan
		srcUser,_ := user.GetUserByName(msg.Src)
		kafka.SendToKafka(msg)	
		for username,clientWs := range clients {
			if !srcUser.IsMyFriend(username) {
				continue
			}
			if err := clientWs.WriteJSON(msg);err != nil {
				clientWs.Close()
				delete(clients, username)
			}
			//修改redis在线状态
			if msg.MessageType == OnlineRemind {
				srcUser.Status = 1
			}else{
				srcUser.Status = 0
			}
			utils.SetValue(msg.Src,srcUser,3600)
			// 推送到kafka，进行数据库写入
		}
	}
}

// kafka消费者　处理消息
func kafkaConsumer(){
	consumer := kafka.GetConsumer()
	for {
		msg := <-consumer.Messages()
		m := &Message{}
		json.Unmarshal(msg.Value, &m)
		switch m.MessageType{
		// 用户在线状态更新到数据库	
		case OnlineRemind,OfflineRemind:
			fmt.Println(m.Message)
			srcUser,_ := user.GetUserByName(m.Src)
			srcUser.UserOnlineStatus(m.MessageType%2)
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