package main

import (
	"flag"
	"golang-app-tutorial/trace"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

const (
	googleClientId     = "889254910425-iupbf91rpnb4e7ub90mq43jgk4r9kuqe.apps.googleusercontent.com"
	googleClientSecret = "384dIEsPznMKJmKweLdTaYgm"
)

// 这个类型满足http.Handler接口，所以可以直接传入 http.Handle函数
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Do只会在第一次调用时执行，无论多少个goroutine执行这段代码，都只执行一次
	// 模板只需要编译一次，可以写一个函数加载模板并编译，返回这个模板（适合单goroutine执行）
	// 这是一种懒加载的方式，只有在请求（第一次）到达的时候才会编译模板
	t.once.Do(func() {
		// Must是对返回*Template类型值的函数的包装函数，遮蔽了error可以不为nil的情况
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(writer, request)
}

func main() {
	var verbose bool

	addr := flag.String("addr", ":8080", "the addr of the application.")
	flag.BoolVar(&verbose, "v", false, "open verbose mode")
	// parse the flag
	flag.Parse()
	r := newRoom()
	if verbose {
		r.tracer = trace.New(os.Stdout)
	}

	// setup gomniauth
	gomniauth.SetSecurityKey("Life goes on!")
	gomniauth.WithProviders(
		facebook.New("key", "secret", "http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret", "http://localhost:8080/auth/callback/github"),
		google.New(googleClientId, googleClientSecret, "http://localhost:8080/auth/callback/google"),
	)

	/*
		goweb, pat, routes, or mux 如果需要更细致的路由管理，可以使用这些第三方包
	*/
	http.Handle("/asset/", http.StripPrefix("/asset", http.FileServer(http.Dir(filepath.Join("templates", "asset")))))
	// htmlHander的方法是指针类型的接收参数，所以传入Handle函数的也应该是指针类型
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	// 监听localhost 8080，省略ip则监听localhost
	log.Println("starting web server on: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
