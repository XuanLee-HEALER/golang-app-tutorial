package main

import (
	"errors"
)

// 自定义错误类型，获取头像url失败
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")
var ErrCastAvatarURL = errors.New("chat: Unable to cast avatar URL to string")

type Avatar interface {
	// 不同方式都可以通过同样的接口来获取头像url
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

// GetAvatarURL 忽略接收参数的名字可以让go抛弃对于实例自身的引用，且可以提醒开发者不使用这个引用
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	var urlStr string
	if url, ok := c.userData["avatar_url"]; !ok {
		return "", ErrNoAvatarURL
	} else if urlStr, ok = url.(string); !ok {
		return "", ErrCastAvatarURL
	}
	return urlStr, nil
}
