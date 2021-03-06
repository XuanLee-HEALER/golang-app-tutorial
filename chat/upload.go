package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func uploadHandler(rw http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// ioutil.ReadAll会读取Reader中的字节直至全部读取完
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := path.Join("avatars", userId+path.Ext(header.Filename))
	err = os.WriteFile(filename, data, 0600)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(rw, "Successful")
}
