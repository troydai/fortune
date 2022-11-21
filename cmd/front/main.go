// Front HTTP server
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"go.uber.org/zap"
)

const (
	_envDatastoreUrl = "DATASTORE_NAME"
	_socketPath      = "unix:///run/spire/sockets/agent.sock"
)

type (
	server struct {
		client *http.Client
		logger *zap.Logger
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		startTls(ctx, http.HandlerFunc(s.fortune))
		wg.Done()
	}()

	wg.Wait()
}

func startTls(ctx context.Context, h http.Handler) error {
	// Create a `workloadapi.X509Source`, it will connect to Workload API using provided socket.
	// If socket path is not defined using `workloadapi.SourceOption`, value from environment variable `SPIFFE_ENDPOINT_SOCKET` is used.
	source, err := workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(_socketPath)))
	if err != nil {
		return fmt.Errorf("unable to create X509Source: %w", err)
	}
	defer source.Close()

	// Allowed SPIFFE ID
	// clientID := spiffeid.RequireFromString("spiffe://example.org/client")

	// Create a `tls.Config` to allow mTLS connections, and verify that presented certificate has SPIFFE ID `spiffe://example.org/client`
	tlsConfig := tlsconfig.MTLSServerConfig(source, source, tlsconfig.AuthorizeAny())
	server := &http.Server{
		Addr:      ":8081",
		TLSConfig: tlsConfig,
		Handler:   h,
	}

	return server.ListenAndServeTLS("", "")
}

func (s *server) fortune(w http.ResponseWriter, req *http.Request) {
	s.logger.Info(
		"hit fortune",
		zap.String("host", os.Getenv("HOSTNAME")),
		zap.Any("peer-uris", extractCaller(req)),
	)

	resp, err := s.client.Get(os.Getenv(_envDatastoreUrl))
	if err != nil {
		s.error(w, "fail to access datastore", err)
		return
	}

	defer resp.Body.Close()
	fc, err := io.ReadAll(resp.Body)
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

func extractCaller(req *http.Request) []string {
	if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
		return nil
	}

	var retval []string
	for _, c := range req.TLS.PeerCertificates {
		for _, u := range c.URIs {
			retval = append(retval, u.String())
		}
	}
	return retval
}
