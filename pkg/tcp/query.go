package tcp

import (
	"encoding/binary"
	"net"
	"net/http"

	"github.com/pkg/errors"
)

type QuerySrv interface {
	ReadClientQuery() ([]byte, error)
	ReadClientQueryStr() (string, error)
	WriteClientResponse(code ResponseCode, msg []byte) error
	WriteClientResponseStr(code ResponseCode, msg string) error
}

type QueryClient interface {
	ReadServerResponse() (ResponseCode, []byte, error)
	ReadServerResponseStr() (ResponseCode, string, error)
	WriteServerQuery(msg []byte) error
	WriteServerQueryStr(msg string) error
}

type ResponseCode uint32

const (
	TCPCodeOK ResponseCode = http.StatusOK

	TCPCodeBadRequest    ResponseCode = http.StatusBadRequest
	TCPCodeNotAuthorized ResponseCode = http.StatusUnauthorized
	TCPCodeForbidden     ResponseCode = http.StatusForbidden

	TCPCodeInternalError ResponseCode = http.StatusInternalServerError
)

type query struct {
	conn net.Conn

	code ResponseCode
}

func NewQuery(conn net.Conn) *query {
	return &query{
		conn: conn,

		code: TCPCodeOK,
	}
}

func (q *query) ReadClientQuery() ([]byte, error) {
	msgLenLen := make([]byte, 4) // takes the first four bytes: the Length Message.
	ind, err := q.conn.Read(msgLenLen)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading TCP Package")
	}

	msgLenVal := binary.BigEndian.Uint32(msgLenLen[:ind]) // convert the bytes into int32.

	bodyBuf := make([]byte, msgLenVal) // create the byte buffer for the body.
	ind, err = q.conn.Read(bodyBuf)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading TCP Package")
	}

	return bodyBuf[:ind], nil
}

func (q *query) ReadServerResponse() (ResponseCode, []byte, error) {
	codeBuf := make([]byte, 4) // takes the second four bytes: the Length Message.
	ind, err := q.conn.Read(codeBuf)
	if err != nil {
		return 0, nil, errors.Wrap(err, "Error reading TCP Package")
	}
	code := binary.BigEndian.Uint32(codeBuf[:ind])

	msgLenLen := make([]byte, 4) // takes the second four bytes: the Length Message.
	ind, err = q.conn.Read(msgLenLen)
	if err != nil {
		return ResponseCode(code), nil, errors.Wrap(err, "Error reading TCP Package")
	}

	msgLenVal := binary.BigEndian.Uint32(msgLenLen[:ind]) // convert the bytes into int32.

	bodyBuf := make([]byte, msgLenVal) // create the byte buffer for the body.
	ind, err = q.conn.Read(bodyBuf)
	if err != nil {
		return ResponseCode(code), nil, errors.Wrap(err, "Error reading TCP Package")
	}

	return ResponseCode(code), bodyBuf[:ind], nil
}

func (q *query) ReadClientQueryStr() (string, error) {
	msgBytes, err := q.ReadClientQuery()
	if err != nil {
		return "", err
	}

	if len(msgBytes) > 0 {
		return string(msgBytes), nil
	}

	return "", nil
}

func (q *query) ReadServerResponseStr() (ResponseCode, string, error) {
	code, msgBytes, err := q.ReadServerResponse()
	if err != nil {
		return code, "", err
	}

	if len(msgBytes) > 0 {
		return code, string(msgBytes), nil
	}

	return code, "", nil
}

func (q *query) WriteClientResponse(code ResponseCode, msg []byte) error {
	if q == nil {
		panic("not valid response object")
	}

	codeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(codeBuf, uint32(code))

	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(msg)))

	codeBuf = append(codeBuf, lenBuf...)

	_, err := q.conn.Write(append(codeBuf, msg...))

	return err
}

func (q *query) WriteServerQuery(msg []byte) error {
	if q == nil {
		panic("not valid response object")
	}

	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(msg)))

	_, err := q.conn.Write(append(lenBuf, msg...))

	return err
}

func (q *query) WriteClientResponseStr(code ResponseCode, msg string) error {
	if q == nil {
		panic("not valid response object")
	}

	msgBuf := []byte(msg)

	return q.WriteClientResponse(code, msgBuf)
}

func (q *query) WriteServerQueryStr(msg string) error {
	if q == nil {
		panic("not valid response object")
	}

	msgBuf := []byte(msg)

	return q.WriteServerQuery(msgBuf)
}
