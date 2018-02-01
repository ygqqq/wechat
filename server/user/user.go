package user
import(
	"../utils"
	"time"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/gin"
)
const URL = "127.0.0.1:27017"	//mongodb数据库地址
var dbSession *mgo.Session 
type User struct{
	Id_   bson.ObjectId `bson:"_id"`
	UserName string		`bson:"username"`
	Password string		`bson:"password"`
	NickName string		`bson:"nickname"`
	Gender string		`bson:"gender"`  //0：男　1：女
	CreateAt time.Time  `bson:"create_at"`
	Friends []string	`bson:"friends"`
	Status	int			`bson:"status"`	 //0:下线　1:在线
}
// 判断某人是否为自己的好友
func (u *User)IsMyFriend(username string) bool{
	for _,un := range u.Friends{
		if un == username {
			return true
		}
	}
	return false
}

// 添加好友
func (u *User)AddFriendByName(username string) {
	//将好友信息同步到redis
	u.Friends = append(u.Friends,username)
	utils.SetValue(u.UserName,*u,3600)
	//获取mongodb数据库连接
	session := getDbSession()
	defer session.Close()
	db := session.DB("wechat").C("users")
	db.Update(bson.M{"username": u.UserName},bson.M{"$push": bson.M{"friends": username}})
}
// 获取某用户的所有好友
func (u *User)GetAllFriends() []User{
	//获取mongodb数据库连接
	session := getDbSession()
	defer session.Close()
	db := session.DB("wechat").C("users")
	var users []User
	db.Find(bson.M{"username": bson.M{"$in": u.Friends}}).All(&users)
	return users
}
// 用户下线修改 status 0:离线　1:在线
func (u *User)SetUserOnlineStatus(status int){
	//获取mongodb数据库连接
	session := getDbSession()
	defer session.Close()
	db := session.DB("wechat").C("users")
	db.Update(bson.M{"_id": u.Id_},bson.M{"$set": bson.M{"status": status}})

}

// 防止每次操作mongo都重新建立连接
func getDbSession() *mgo.Session{
	if dbSession == nil{
		var err error
		dbSession, err = mgo.Dial(URL)
		if err != nil {
			panic(err) 
		}
	}
	return dbSession.Clone()
}

// 注册用户
func Register(c *gin.Context){
	//获取mongodb数据库连接
	session := getDbSession()
	defer session.Close()
	db := session.DB("wechat").C("users")
	//根据用户提交的注册信息转成User对象
	var user User
	err := c.BindJSON(&user)
	if err != nil{
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "服务器错误",
		})
		return
	}
	//一些用户信息的基本验证
	if len(user.UserName) == 0 || len(user.Password) == 0{
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "用户名或密码非法",
		})
		return
	}
	err = db.Find(bson.M{"username": user.UserName}).One(&user)
	if err == nil{
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "用户名已存在",
		})
		return
	}
	//通过验证后，将用户提交的信息插入数据库
	insertUser := &User{
		Id_ : bson.NewObjectId(),
		UserName : user.UserName,
		Password : utils.Md5(user.Password),
		NickName : user.NickName,
		Gender : user.Gender,
		CreateAt : time.Now(), 
		Status : 1,
	}
	err = db.Insert(*insertUser)
	if err != nil {
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "注册失败",
		})
		return
	}else{
		//注册成功后，将用户信息写入redis、cookie
		utils.SetValue(insertUser.UserName,*insertUser,3600)
		cookie := http.Cookie{Name: "username", Value: (*insertUser).UserName, Path: "/", MaxAge: 86400}
		http.SetCookie(c.Writer, &cookie)
		insertUser.Password = ""
		strUser,_ := json.Marshal(insertUser)
		c.JSON(200, gin.H{
			"success" : true,
			"msg": "注册成功",
			"user": string(strUser),
		})
		return
	}
}

// 用户登陆
func Login(c *gin.Context){
	//获取mongodb数据库连接
	session := getDbSession()
	defer session.Close()
	db := session.DB("wechat").C("users")

	//根据用户提交的用户名和密码转成User对象
	var postUser,dbUser User
	er := c.BindJSON(&postUser)
	if er != nil{
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "服务器错误",
		})
		return
	}
	//一些用户信息的基本验证
	if len(postUser.UserName) == 0 || len(postUser.Password) == 0{
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "用户名或密码非法",
		})
		return
	}
	//根据用户提交的用户名取到数据库中的user，如果不存在直接返回错误
	if err := db.Find(bson.M{"username": postUser.UserName}).One(&dbUser); err != nil{
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "用户名不存在",
		})
		return
	}
	//判断密码是否正确
	if utils.Md5(postUser.Password) != dbUser.Password{
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "密码错误",
		})
		return
	}
	//密码正确则登陆成功，写入redis、cookie
	dbUser.Status = 1 //修改用户的在线状态
	utils.SetValue(dbUser.UserName,dbUser,3600)
	cookie := http.Cookie{Name: "username", Value: dbUser.UserName, Path: "/", MaxAge: 86400}
	http.SetCookie(c.Writer, &cookie)
	dbUser.Password = ""

	strUser,_ := json.Marshal(dbUser)
	c.JSON(200, gin.H{
		"success" : true,
		"msg" : "登陆成功",
		"user" : strUser,
	})
	return
}

// 根据用户名获得用户
func GetUserByName(username string) (User,error){
	//先从redis获取用户，redis不存在则去数据库取
	if str,err := utils.GetValue(username); len(str) >0 && err == nil{
		user := &User{}
		json.Unmarshal([]byte(str), user)
		return *user,nil
	}
	//获取mongodb数据库连接
	session := getDbSession()
	defer session.Close()
	db := session.DB("wechat").C("users")
	var dbUser User
	err := db.Find(bson.M{"username": username}).One(&dbUser)
	return dbUser,err
}


