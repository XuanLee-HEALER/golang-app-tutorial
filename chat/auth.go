package main

import (
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// ChatUser 代表登录聊天室的用户
type ChatUser interface {
	UniqueID() string
	Avatar() string
}

type chatUser struct {
	// 类型嵌入，User接口包含了AvatarURL方法，所以不需要专门实现
	goth.User
	uniqueID string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

func (u chatUser) Avatar() string {
	return u.AvatarURL
}

type authHandler struct {
	next http.Handler
}

func (a *authHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	cookies, err := request.Cookie("_gothic_session")
	if err == http.ErrNoCookie || cookies.Value == "" {
		// not authenticated
		writer.Header().Set("Location", "/login")
		writer.WriteHeader(http.StatusTemporaryRedirect) // send response
		return
	}
	if err != nil {
		// other error
		http.Error(writer, err.Error(), http.StatusInternalServerError) // send error response
		return
	}
	// success call next handler
	a.next.ServeHTTP(writer, request)
}

// MustAuth 返回一个AuthHandler
func MustAuth(handler http.Handler) *authHandler {
	return &authHandler{next: handler}
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	if _, err := gothic.CompleteUserAuth(writer, request); err == nil {
		writer.Header().Set("Location", "/")
		writer.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		gothic.BeginAuthHandler(writer, request)
	}

}

/*
	OAuth2的工作流（用户视角）
	1. 用户选择要使用的其它服务提供商登录第三方app
	2. 用户跳转到服务提供商的网站（询问客户是否允许第三方app访问他们的个人信息）
	3. 用户从OAuth2服务提供商登录并且同意第三方app访问客户个人信息的请求
	4. 用户跳转回第三方app（携带请求码 request code）
	5. 在后台，第三方app发送授权码到服务提供商，服务提供商再发回权限token
	6. app使用权限token来请求服务商，抓取tweet或post

	goauth2是完全用go实现的OAuth2协议，作者是核心go开发组成员
	goauth2激励了gomniauth项目（ruby的omniauth的go版本），这个项目提供了统一的访问OAuth2服务的解决方案，当OAuth3（或下一代授权协议）到来时，gomniauth
	可以在客户代码不变的情况下修改自身的代码

	google client id: 889254910425-iupbf91rpnb4e7ub90mq43jgk4r9kuqe.apps.googleusercontent.com
	google client secret: 384dIEsPznMKJmKweLdTaYgm
*/
