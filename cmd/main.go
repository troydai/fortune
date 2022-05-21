package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	fmt.Println("starting HTTP server ...")

	http.Handle("/echo", http.HandlerFunc(echo))

	http.ListenAndServe(":8080", nil)
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
