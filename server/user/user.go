package user
import(
	"../utils"
	"time"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/gin"
)
const URL = "127.0.0.1:27017"	//mongodb数据库地址
var dbSession *mgo.Session 
type User struct{
	UserName string		`bson:"username"`
	Password string		`bson:"password"`
	NickName string		`bson:"nickname"`
	Gender string		`bson:"gender"` //0：男　1：女
	CreateAt time.Time  `bson:"create_at"`
	Friends []string	`bson:"friends"`
	Status	int		//0:下线　1:在线
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
		c.JSON(200, gin.H{
			"success" : true,
			"msg": "注册成功",
		})
		return
	}
}

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
	utils.SetValue(dbUser.UserName,dbUser,3600)
	cookie := http.Cookie{Name: "username", Value: dbUser.UserName, Path: "/", MaxAge: 86400}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(200, gin.H{
		"success" : true,
		"msg": "登陆成功",
	})
	return
}
