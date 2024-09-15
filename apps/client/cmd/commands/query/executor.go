package query

import (
	"context"
	"fmt"
	"os"

	"github.com/Isotere/kw-database/pkg/tcp"
)

const (
	success = 0
	fail    = 1
)

func Handle(query string, host string, port int) {
	os.Exit(run(query, host, port))
}

func run(queryArg string, host string, port int) (exitCode int) {
	log := tcp.NewFakeLogger()

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	client := tcp.NewClient(log, host, port)

	q := query{q: queryArg}

	err := client.ProcessRequest(ctx, q.handlerFn)
	if err != nil {
		log.WithError("query error", err)
		return fail
	}

	return success
}

type query struct {
	q    string
	code uint32
	res  string
}

func (q *query) handlerFn(ctx context.Context, cl tcp.QueryClient) {
	err := cl.WriteServerQueryStr(q.q)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	code, msg, err := cl.ReadServerResponseStr()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(code)
	fmt.Println(msg)
}
