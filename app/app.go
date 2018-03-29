package app

import (
	"reflect"
	"risk-ext/config"

	"github.com/kataras/iris/context"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

var (
	app    = iris.New()
	method = []string{"Get", "Post", "Put", "Delete"}

	paths = context.Map{}
)

var conf = iris.Configuration{ // default configuration:
	DisableStartupLog:                 false,
	DisableInterruptHandler:           false,
	DisablePathCorrection:             false,
	EnablePathEscape:                  false,
	FireMethodNotAllowed:              false,
	DisableBodyConsumptionOnUnmarshal: false,
	DisableAutoFireStatusCode:         false,
	TimeFormat:                        "2006-1-2 15:04:05",
	Charset:                           "UTF-8",
}

func init() {
	if config.GetBool("debug") {
		app.Logger().SetLevel("debug")
	}

	if config.GetBool("logs") {
		app.Use(recover.New())
		app.Use(logger.New())
	}
}
func AddPath(path string, obj interface{}) {
	paths[path] = obj
}

func Run() {
	host := config.GetString("host")
	port := config.GetString("port")

	for k, m := range paths {
		v := reflect.ValueOf(m)
		for _, md := range method {
			fn := v.MethodByName(md)
			if fn.IsValid() {
				args_ := []reflect.Value{reflect.ValueOf(k), reflect.ValueOf(func(ctx iris.Context) {
					args := []reflect.Value{reflect.ValueOf(ctx)}
					rs := fn.Call(args)
					ctx.StatusCode(int(rs[0].Int()))
					ctx.JSON(rs[1].Interface())
				})}
				a := reflect.ValueOf(app)
				afn := a.MethodByName(md)
				if afn.IsValid() {
					afn.Call(args_)
				}
			}
		}

	}
	if port == "" {
		port = "80"
	}
	app.Run(iris.Addr(host+":"+port), iris.WithConfiguration(conf))
}
