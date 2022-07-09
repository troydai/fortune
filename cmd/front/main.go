// Front HTTP server
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

const (
	_envDatastoreUrl = "DATASTORE_NAME"
)

type (
	server struct {
		client *http.Client
		logger *zap.Logger
	}
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("fail to provision zap logger: %s", err.Error())
	}

	s := &server{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logger,
	}

	http.Handle("/fortune", http.HandlerFunc(s.fortune))
	http.ListenAndServe(":8080", nil)
}

func (s *server) fortune(w http.ResponseWriter, req *http.Request) {
	resp, err := s.client.Get(os.Getenv(_envDatastoreUrl))
	if err != nil {
		s.error(w, "fail to access datastore", err)
		return
	}

	defer resp.Body.Close()
	fc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.error(w, "fail to read datastore response", err)
		return
	}

	if _, err := w.Write(fc); err != nil {
		s.error(w, "fail to write response", err)
		return
	}
}

func (s *server) error(w http.ResponseWriter, msg string, err error) {
	s.logger.Error(msg, zap.Error(err))
	http.Error(w, msg, http.StatusInternalServerError)
}
