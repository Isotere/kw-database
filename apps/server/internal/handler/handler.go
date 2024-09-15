package handler

import (
	"context"

	"github.com/Isotere/kw-database/apps/server/internal"
	"github.com/Isotere/kw-database/pkg/tcp"
)

type Handler struct {
	log internal.Log
}

func New(log internal.Log) *Handler {
	return &Handler{
		log: log,
	}
}

func (h *Handler) Handle(_ context.Context, q tcp.QuerySrv) {
	inMsg, err := q.ReadClientQueryStr()
	if err != nil {
		h.log.WithError("Error receiving message", err)
	}

	h.log.Info("Receive message: ", inMsg)

	h.log.Info("Send message start")

	err = q.WriteClientResponseStr(tcp.TCPCodeOK, "Message Received Success! Try more!")
	if err != nil {
		h.log.WithError("Error sending message", err)
	}
}
