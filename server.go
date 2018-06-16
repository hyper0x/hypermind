package main

import (
	"flag"
	"fmt"
	"hypermind/core/base"
	"hypermind/core/controller"
	"net/http"
	_ "net/http/pprof"
)

var serverPort int

func init() {
	flag.IntVar(&serverPort, "port", 9091, "the server (http listen) port")
}

func main() {
	flag.Parse()
	fileServer := http.FileServer(http.Dir("web"))
	http.Handle("/css/", fileServer)
	http.Handle("/js/", fileServer)
	http.Handle("/img/", fileServer)
	http.HandleFunc("/", controller.RequestDispatcher)
	http.HandleFunc("/register", controller.Register)
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/logout", controller.Logout)
	http.HandleFunc("/upload", controller.Upload)
	http.HandleFunc("/get-cv", controller.GetCv)
	http.HandleFunc("/auth_code", controller.GetAuthCodeForAdmin)
	http.HandleFunc("/user_list", controller.GetUserListForAdmin)
	base.Logger().Infof("Starting hypermind http server (port=%d)...\n", serverPort)
	err := http.ListenAndServe(":"+fmt.Sprintf("%d", serverPort), nil)
	if err != nil {
		base.Logger().Fatalln("ListenAndServeError: ", err)
	} else {
		base.Logger().Infoln("Hypermind http server is started.")
	}
}
