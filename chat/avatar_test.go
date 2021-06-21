package main

import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	testURL := "http://url-to-avatar/"
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}
	if url != testURL {
		t.Error("AuthAvatar.GetAvatarURL should return correct url")
	}

	testURL2 := 1
	client.userData["avatar_url"] = testURL2
	url, err = authAvatar.GetAvatarURL(client)
	if err != ErrCastAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrCastAvatarURL when url type error")
	}
}
