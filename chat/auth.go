package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
)

type authHandler struct {
	next http.Handler
}

func (a *authHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, err := request.Cookie("auth")
	if err == http.ErrNoCookie {
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

func MustAuth(handler http.Handler) *authHandler {
	return &authHandler{next: handler}
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	segs := strings.Split(request.URL.Path, "/")
	var action, provider string

	if len(segs) < 4 {
		action = "notfound"
	} else {
		action = segs[2]
		provider = segs[3]
	}

	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(writer, fmt.Sprintf("error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		// 第一个state参数是编码后的map数据，发送给授权服务商，服务商会把这些数据再发回给回调url。第二个参数也是发送给服务商的map数据，可能会改变授权行为，例如scope参数
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(writer, fmt.Sprintf("error when trying to GetBeginAuthURL for %s: %s", provider, err), http.StatusInternalServerError)
		}
		writer.Header().Set("Location", loginUrl)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	default:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "auth action %s not supported", action)
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
