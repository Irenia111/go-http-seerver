package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "net/http/pprof"
)

func main() {
	// 设置环境变量
	os.Setenv("VERSION", "1")
	// fmt.Println("VERSION:", os.Getenv("VERSION"))

	// http server 基础写法
	// http.HandleFunc("/", rootHandler)
	// err := http.ListenAndServe(":80", nil)

	// http server 中间件
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/", rootHandler)

	// err := http.ListenAndServe(":6060", nil)
	// mux.HandleFunc("/debug/pprof/", pprof.Index)
	// mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	
	// server := &http.Server{
	// 	Addr:         ":6000",
	// 	ReadTimeout:  60 * time.Second,
	// 	WriteTimeout: 60 * time.Second,
	// 	Handler:      mux,
	// }
	// server.ListenAndServe()

	http.ListenAndServe(":80", mux)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
		fmt.Printf("Server failed: %s\n", err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// io.WriteString(w, "ok\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	log.Print("Request host is: ", r.Host)
	log.Print("Response status code: ", http.StatusAccepted)
	// 从访问 url 中拿出 user
	user := r.URL.Query().Get("user")
	userInfo := "stranger user"
	if user != "" {
		userInfo = fmt.Sprintf("User[%s] login", user)
		log.Print(userInfo)
		// io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	}
	log.Print(userInfo)
	// io.WriteString(w, "VERSION:"+os.Getenv("VERSION")+"\n")
	// io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		// 设置 response header
		w.Header().Set(k, strings.Join(v, ","))
		// io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("VERSION", os.Getenv("VERSION"))
	resp := make(map[string]string)
	resp["message"] = "Status Accepted"
	resp["user"] = user

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	// 状态码
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResp)
}
