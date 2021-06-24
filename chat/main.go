package main

import (
	"golang-app-tutorial/trace"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/namsral/flag"
	"github.com/stretchr/objx"
)

var (
	addr    string
	verbose bool
	avatars Avatar = TryAvatars{
		UseFileSystemAvatar,
		UseAuthAvatar,
		// 总能获取到一个默认的头像地址
		UseGravatarAvatar,
	}
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
	data := map[string]interface{}{
		"Host": request.Host,
	}
	if authCookies, err := request.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookies.Value)
	}

	t.templ.Execute(writer, data)
}

func init() {
	flag.StringVar(&addr, "addr", ":8080", "the address of the application")
	flag.BoolVar(&verbose, "verbose", false, "open verbose mode")
	flag.Parse()

	// 创建头像目录，忽略任何意外情况
	_ = os.Mkdir("avatars", os.ModeDir)

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:8080/"),
	)

	// 初始化goth的coockie store
	key := "funnyboy" // Replace with your SESSION_SECRET or similar
	maxAge := 60 * 30 // 30 days
	isProd := false   // Set to true when serving over https
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store
}

func main() {
	var verbose bool

	r := newRoom()
	if verbose {
		r.tracer = trace.New(os.Stdout)
	}

	/*
		goweb, pat, routes, or mux 如果需要更细致的路由管理，可以使用这些第三方包
	*/
	// StripPrefix会移除前缀，FileServer用来处理静态资源，Dir函数决定哪些文件夹是可以被访问的
	rtr := mux.NewRouter()

	// htmlHander的方法是指针类型的接收参数，所以传入Handle函数的也应该是指针类型
	rtr.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	rtr.Handle("/login", &templateHandler{filename: "login.html"})
	rtr.HandleFunc("/auth/{provider}", loginHandler)
	rtr.Handle("/room", r)
	// 添加logout功能，删除coockie并跳转回首页
	rtr.HandleFunc("/logout", func(rw http.ResponseWriter, r *http.Request) {
		http.SetCookie(rw, &http.Cookie{
			Name:   "auth",
			Value:  "", // 不是所有浏览器都会强制删除cookie，所以需要显示设置值
			Path:   "/",
			MaxAge: -1, // coockie应该被浏览器立即删除
		})
		rw.Header().Set("Location", "/chat")
		rw.WriteHeader(http.StatusTemporaryRedirect)
	})
	rtr.Handle("/upload", MustAuth(&templateHandler{filename: "upload.html"}))
	rtr.HandleFunc("/uploader", uploadHandler)

	http.Handle("/asset/", http.StripPrefix("/asset", http.FileServer(http.Dir(filepath.Join("templates", "asset")))))
	http.Handle("/avatars/", http.StripPrefix("/avatars", http.FileServer(http.Dir("./avatars"))))
	http.Handle("/", rtr)

	go r.run()
	// 监听localhost 8080，省略ip则监听localhost
	log.Println("starting web server on: ", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
