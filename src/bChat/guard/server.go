package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ListenAndServe() {

	// TODO: 读取配置文件
	host := "127.0.0.1"
	port := 8876
	router := NewRouter()

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("guard listen address is: %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
