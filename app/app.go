package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"risk-ext/config"
	"strings"

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

func App() *iris.Application {

	for k, m := range paths {
		v := reflect.ValueOf(m)
		for _, md := range method {
			fn := v.MethodByName(md)
			if fn.IsValid() {
				args_ := []reflect.Value{reflect.ValueOf(k), reflect.ValueOf(func(ctx iris.Context) {
					authFunc := v.MethodByName("Auth")
					authResult := true
					if authFunc.IsValid() {
						result := authFunc.Call([]reflect.Value{reflect.ValueOf(ctx)})
						authResult = result[0].Bool()
					}

					if !authResult {
						ctx.StatusCode(403)
						ctx.JSON("没有权限")
						return
					}

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
	return app
}

func Run() {
	app := App()
	host := config.GetString("host")
	port := config.GetString("port")
	if port == "" {
		port = "80"
	}
	app.Run(iris.Addr(host+":"+port), iris.WithConfiguration(conf))
}
func HttpClient(url string, args interface{}, method string, result interface{}, token ...string) error {
	client := &http.Client{}
	var params = ""
	var contentType = "application/x-www-form-urlencoded"

	if reflect.TypeOf(args).String() != "string" {
		jsonData, err := json.Marshal(args)
		if err != nil {
			return err
		}
		params = string(jsonData)
		contentType = "application/json"
	} else {
		params = args.(string)
	}

	var req *http.Request
	var err error
	if method == "GET" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, strings.NewReader(params))
	}

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", contentType)

	if len(token) != 0 {
		req.Header.Add("X-Token", token[0])
	}
	//fmt.Println(contentType, url, params)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	jsonStr := string(body)
	err = json.Unmarshal([]byte(jsonStr), result)
	return err
}
