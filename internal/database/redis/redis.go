package redis

//save cach, if have return, sycunetlock()mutex, doing try late,
import (
	"github.com/go-redis/redis"
	"time"
)

// for execution context
type Context struct {
	config  *Config       //points to the redis config
	rClient *redis.Client //rCLient can be shared by multiple go routines
}

func CreateNewRedisContext(config *Config) (rContext *Context, err error) {

}

func (rContext *Context) CheckHealth() (pong string, err error) {

}

func (rContext *Context) Cleanup() {

}
func (rContext *Context) SetKeyValue(key, val string, kDuration time.Duration) (err error) {

}

func (rContext *Context) DeleteKey(key string) {

}

func (rContext *Context) GetValue(key string) (val string, err error) {

}

func (rContext *Context) GetValidTime(key string) (exp time.Duration, err error) {

}
