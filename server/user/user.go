package user
import(
	"../utils"
	"time"
	
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
	session := getDbSession()
	defer session.Close()
	db := session.DB("wechat").C("users")
	
	var user User
	er := c.BindJSON(&user)
	if er != nil{
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "服务器错误",
		})
		return
	}

	if len(user.UserName) == 0 || len(user.Password) == 0{
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "用户名或密码非法",
		})
		return
	}
	err := db.Find(bson.M{"username": user.UserName}).One(&user)
	if err == nil{
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "用户名已存在",
		})
		return
	}

	insertUser := &User{
		UserName : user.UserName,
		Password : utils.Md5(user.Password),
		NickName : user.NickName,
		Gender : user.Gender,
		CreateAt : time.Now(), 
	}
	err = db.Insert(*insertUser)
	if err != nil {
		c.JSON(200, gin.H{
			"success" : false,
			"msg": "注册失败",
		})
		return
	}else{
		//写入redis
		utils.SetValue(insertUser.UserName,*insertUser,3600)
		c.JSON(200, gin.H{
			"success" : true,
			"msg": "注册成功",
		})
		return
	}
}

func login(user *User){
	utils.SetValue(user.UserName,*user,3600)
}