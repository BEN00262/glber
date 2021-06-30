package slb

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/ipsn/go-libtor"
	"github.com/pkg/errors"
)

// Server is an interface for representing server load balancer implementation.
type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type serverLoadBalancer struct {
	*Config
	*http.Server
	RequestDirector func(target *url.URL) func(*http.Request)
	HandlerDirector HandlerDirector
}

// CreateSLB returns Server implementation(*serverLoadBalancer) from the given Config.
func CreateSLB(cfg *Config, ops ...Option) (Server, error) {
	err := cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	sbl := &serverLoadBalancer{
		Config: cfg,
		RequestDirector: func(target *url.URL) func(*http.Request) {
			return func(req *http.Request) {
				req.URL.Scheme = target.Scheme
				req.URL.Host = target.Host
				req.URL.Path = target.Path

				if target.RawQuery == "" || req.URL.RawQuery == "" {
					req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
				} else {
					req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
				}
				if _, ok := req.Header["User-Agent"]; !ok {
					req.Header.Set("User-Agent", "")
				}
			}
		},
		HandlerDirector: cfg.Balancing.CreateHandler,
	}
	sbl.apply(ops...)

	sbl.Server = &http.Server{
		Handler: sbl.HandlerDirector(cfg.BackendServerConfigs.getURLs(), sbl),
	}

	return sbl, nil
}

func (s *serverLoadBalancer) apply(ops ...Option) {
	for _, op := range ops {
		op(s)
	}
}

func (s *serverLoadBalancer) Proxy(target *url.URL, w http.ResponseWriter, req *http.Request) {
	(&httputil.ReverseProxy{
		Director: s.RequestDirector(target),
	}).ServeHTTP(w, req)
}

func (s *serverLoadBalancer) ListenAndServe() error {
	// we really dont have any control over this :)
	// addr := s.Config.Host + ":" + s.Config.Port

	// var (
	// 	ls  net.Listener
	// 	err error
	// )

	// we wont handle the tls shit here ( its tor nigga )
	// if s.Config.TLSConfig.Enabled {
	// 	ls, err = createTLSListenter(addr, s.Config.TLSConfig.CertKey, s.Config.TLSConfig.KeyKey)
	// 	if err != nil {
	// 		return errors.Wrap(err, "faild to create tls lisner")
	// 	}
	// } else {
	// this listener created will be a tor listener
	ls, err := createListener()
	if err != nil {
		return errors.Wrap(err, "faild to create listener")
	}

	// }

	// inject the tor network here
	err = s.Server.Serve(ls)
	if err != nil {
		return errors.Wrap(err, "faild to serve")
	}
	return nil
}

// we have created a darkweb thing
func createListener() (net.Listener, error) {
	// create the tor listener
	// Start tor with some defaults + elevated verbosity
	fmt.Println("Starting and registering onion service, please wait a bit...")
	t, err := tor.Start(nil, &tor.StartConf{ProcessCreator: libtor.Creator, DebugWriter: os.Stderr})
	if err != nil {
		log.Panicf("Failed to start tor: %v", err)
	}
	defer t.Close()

	// Wait at most a few minutes to publish the service
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// Create an onion service to listen on any port but show as 80
	onion, err := t.Listen(ctx, &tor.ListenConf{RemotePorts: []int{80}})
	if err != nil {
		log.Panicf("Failed to create onion service: %v", err)
	}
	defer onion.Close()

	fmt.Printf("Please open a Tor capable browser and navigate to http://%v.onion\n", onion.ID)

	return onion, nil
}

// we dont care about this
func createTLSListenter(addr string, certFile, keyFile string) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, errors.Wrapf(err, "faild to load 509 key parir, certFile: %s, keyFile: %s", certFile, keyFile)
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ls, err := tls.Listen("tcp", addr, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "faild to create listener, network: tcp, addr: %s", addr)
	}

	return ls, nil
}

func (s *serverLoadBalancer) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err := s.Server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "faild to shutdown")
	}
	return nil
}
