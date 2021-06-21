package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	_, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	testURL := "http://url-to-avatar/"
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err := authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}
	if url != testURL {
		t.Error("AuthAvatar.GetAvatarURL should return correct url")
	}

	testURL2 := 1
	client.userData["avatar_url"] = testURL2
	_, err = authAvatar.GetAvatarURL(client)
	if err != ErrCastAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrCastAvatarURL when url type error")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{
		"userid": "d72026f779abdc571edd5555edad0266",
	}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("gravatarAvatar.GetAvatarURL should not return an err")
	}
	if url != "//www.gravatar.com/avatar/d72026f779abdc571edd5555edad0266" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0600)
	// 即使测试会失败，defer的函数也会执行
	defer os.Remove(filename)
	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userid": "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return the error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned: %s", url)
	}
}
