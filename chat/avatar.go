package main

import (
	"errors"
	"os"
	"path"
)

// 自定义错误类型，获取头像url失败
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")
var ErrCastAvatarURL = errors.New("chat: Unable to cast avatar URL to string")

type Avatar interface {
	// 不同方式都可以通过同样的接口来获取头像url
	GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar

// GetAvatarURL 如果从一个avatar实现拿不到url，则继续找下一个
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

// GetAvatarURL 忽略接收参数的名字可以让go抛弃对于实例自身的引用，且可以提醒开发者不使用这个引用
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.Avatar()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	// crypto库中有很多加密方法，md5实现了io.Writer接口，可以使用WriteString来向其中写入数据，Sum方法可以获得hash值
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := os.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}

	return "", ErrNoAvatarURL
}
