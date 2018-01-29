package utils
import(
	"github.com/garyburd/redigo/redis"
	"time"
	"reflect"
	"encoding/json"
	"os"
	"strconv"
)
//redis连接池实例
var redisPool *redis.Pool = nil

//获取连接池实例
func getPool() *redis.Pool{
	if redisPool != nil {
		return redisPool
	}
	redisPool = &redis.Pool{
        MaxIdle:     20,
        MaxActive:   200,
        IdleTimeout: 240 * time.Second,
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", "127.0.0.1:6379")
            if err != nil {
                return nil, err
            }
            return c, err
        },
	}
	return redisPool
}

//可以直接设置struct对象和普通字符串
func SetValue(key string,value interface{}, second int) bool{
	c := getPool().Get()
	defer c.Close()
	t := reflect.TypeOf(value).Kind().String()	
	//如果要设置的value是struct类型，则通过将其序列化成json字符串
	if t == "struct" {
		b,_ := json.Marshal(value)
		os.Stdout.Write(b)
		_,err := c.Do("SET",key,b,"EX",strconv.Itoa(second))
		return err == nil
	}else{
		_,err := c.Do("SET",key,value,"EX",strconv.Itoa(second))
		return err == nil
	}
}


func GetValue(key string) (string,error){
	c := getPool().Get()
	defer c.Close()
	rep,err := c.Do("GET",key)
	if err != nil {
		return "",err
	}
	by,err := redis.Bytes(rep,nil)
	if err != nil {
		return "",err
	}
	return string(by[:]),nil
}

// is_key_exit, err := redis.Bool(c.Do("EXISTS", "mykey1"))
// if err != nil {
// 	fmt.Println("error:", err)
// } else {
// 	fmt.Printf("exists or not: %v \n", is_key_exit)
// }

// _, err = c.Do("DEL", "mykey")