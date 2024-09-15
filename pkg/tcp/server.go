package tcp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Server struct {
	log Logger

	srv    net.Listener
	wg     *sync.WaitGroup
	doneCh chan struct{}

	host string
	port int
}

type HandlerSrv func(context.Context, QuerySrv)

func NewServer(log Logger, host string, port int) *Server {
	return &Server{
		log:  log,
		host: host,
		port: port,

		doneCh: make(chan struct{}),
	}
}

func (s *Server) Listen(ctx context.Context, handler HandlerSrv) error {
	if s.srv != nil {
		s.log.Error("tcp server already started")
		return ErrServerAlreadyStarted
	}

	var err error
	s.srv, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return errors.Wrap(err, "Cannot start tcp server")
	}
	defer func() {
		_ = s.srv.Close()
	}()

	s.wg = &sync.WaitGroup{}

	s.log.Info("TCP Server started on: ", s.host, ":", s.port)

	for {
		cli, errA := s.srv.Accept()
		if errA != nil {
			select {
			case <-s.doneCh:
				s.log.Info("TCP Server shutting down by signal")
				return nil
			default:
				return errors.Wrap(errA, "Error accepting incoming connection")
			}
		}

		s.wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, c net.Conn, hd HandlerSrv) {
			defer wg.Done()
			defer func() {
				_ = c.Close()
			}()

			hd(ctx, NewQuery(c))
		}(ctx, s.wg, cli, handler)
	}
}

func (s *Server) Shutdown() {
	close(s.doneCh)
	_ = s.srv.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		s.log.Warning("Timed out waiting for connections to finish.")
		return
	}
}
