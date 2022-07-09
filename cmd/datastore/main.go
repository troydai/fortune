// Back data store server
package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	_fortuneFiles = []string{"fortunes", "literature", "riddles"}
)

type (
	server struct {
		fortuneCookies []string
	}
)

func main() {
	var s server
	if err := s.load(); err != nil {
		log.Fatalf("fail to load fortune cookies: %s", err.Error())
	}

	http.Handle("/fetch", http.HandlerFunc(s.fetch))
	http.ListenAndServe(":8090", nil)
}

func (s *server) load() error {
	cmd := exec.Command("fortune", "-f")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	if !scanner.Scan() {
		return fmt.Errorf("cannot read the first line from stderr")
	}

	line := scanner.Text()
	fdir := line[strings.Index(line, "/"):]

	var content []string
	var current string
	for _, file := range _fortuneFiles {
		path := filepath.Join(fdir, file)

		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fail to read file %s\n", path)
			continue
		}
		s := bufio.NewScanner(f)
		for s.Scan() {
			if s.Text() == "%" {
				if len(current) != 0 {
					content = append(content, current)
					current = ""
				}
				continue
			}
			current = fmt.Sprintf("%s%s\n", current, s.Text())
		}
	}

	s.fortuneCookies = content
	return nil
}

func (s *server) fetch(w http.ResponseWriter, _ *http.Request) {
	idx := rand.Intn(len(s.fortuneCookies))
	w.Write([]byte(s.fortuneCookies[idx]))
}
