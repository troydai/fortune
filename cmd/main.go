package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var _welcome = []byte("welcome")

func main() {
	fmt.Println("starting HTTP server ...")
	defer fmt.Println("stopping HTTP server ...")

	http.Handle("/echo", http.HandlerFunc(echo))
	http.Handle("/", http.HandlerFunc(welcome))
	http.ListenAndServe(":8080", nil)
}

func welcome(w http.ResponseWriter, _ *http.Request) {
	w.Write(_welcome)
}

func echo(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("error: unexpect error during read: %s\n", err.Error())
		return
	}

	if _, err := w.Write(content); err != nil {
		w.WriteHeader(500)
		log.Printf("error: unexpect error during write: %s\n", err.Error())
		return
	}
}
