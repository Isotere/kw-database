package tcp_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Isotere/kw-database/pkg/tcp"
)

const (
	testHost = "localhost"
	testPort = 50000
)

func Test_Server(t *testing.T) {
	t.Run("Start and shutdown success", func(t *testing.T) {
		server := tcp.NewServer(tcp.NewFakeLogger(), testHost, testPort)

		// ToDo: в этом месте кажет data-race с флагом. Без флага все работает, при этом само приложение, по сути
		// использующее такую же схему shutdown билдится с флагом без проблем
		doneCh := make(chan struct{})
		go func(srv *tcp.Server) {
			<-doneCh
			srv.Shutdown()
		}(server)

		errCh := make(chan error)
		go func(srv *tcp.Server) {
			defer close(errCh)
			err := srv.Listen(context.Background(), func(_ context.Context, _ tcp.QuerySrv) {})
			errCh <- err
		}(server)

		select {
		case err := <-errCh:
			assert.NoError(t, err)
			doneCh <- struct{}{}
		case <-time.After(time.Second * 1):
			doneCh <- struct{}{}
		}
	})

	t.Run("Read Write success", func(t *testing.T) {
		server := tcp.NewServer(tcp.NewFakeLogger(), testHost, testPort)

		clientMsg := "some client message"
		serverMsg := "some server message"

		doneCh := make(chan struct{})
		go func(srv *tcp.Server) {
			<-doneCh
			srv.Shutdown()
		}(server)

		errCh := make(chan error)
		go func(srv *tcp.Server) {
			err := srv.Listen(context.Background(), func(_ context.Context, q tcp.QuerySrv) {
				msg, err := q.ReadClientQueryStr()
				if err != nil {
					errCh <- err
					return
				}

				if msg != clientMsg {
					errCh <- fmt.Errorf("invalid client msg exp: %s got: %s", clientMsg, msg)
					return
				}

				err = q.WriteClientResponse(tcp.TCPCodeOK, []byte(serverMsg))
				if err != nil {
					errCh <- err
					return
				}
			})
			errCh <- err
		}(server)

		client := tcp.NewClient(tcp.NewFakeLogger(), testHost, testPort)
		err := client.ProcessRequest(context.Background(), func(_ context.Context, q tcp.QueryClient) {
			err := q.WriteServerQueryStr(clientMsg)
			if err != nil {
				errCh <- err
				return
			}

			code, msg, err := q.ReadServerResponseStr()
			if err != nil {
				errCh <- err
				return
			}

			if code != tcp.TCPCodeOK {
				errCh <- fmt.Errorf("invalid server code exp: %d got: %d", tcp.TCPCodeOK, code)
				return
			}

			if msg != serverMsg {
				errCh <- fmt.Errorf("invalid server msg exp: %s got: %s", serverMsg, msg)
				return
			}
		})
		assert.NoError(t, err)
		close(doneCh)

		err = <-errCh
		assert.NoError(t, err)
		close(errCh)
	})
}
