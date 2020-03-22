package api

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"nordshare/internal/api"
	"nordshare/pkg/hash"
	"nordshare/pkg/note"
)

type (
	input struct {
		Content       string `json:"content"`
		ReadPassword  string `json:"readPassword"`
		WritePassword string `json:"writePassword"`
		TTL           int64  `json:"ttl"`
	}
	output struct {
		ID string `json:"id"`
	}
	Handler struct {
		saver
		encrypter
	}
	saver interface {
		SaveNote(context.Context, note.Note, string) error
	}
	encrypter interface {
		EncryptContent(*note.Note) error
	}
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest, h *Handler) (events.APIGatewayProxyResponse, error) {
	var in input
	if err := json.Unmarshal([]byte(request.Body), &in); err != nil {
		log.Print(err)
		return api.NewResponse(http.StatusBadRequest)
	}
	n := convert(in)
	if err := h.EncryptContent(&n); err != nil {
		log.Print(err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	if err := hash.HashNote(&n); err != nil {
		log.Print(err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	id := request.RequestContext.RequestID
	if err := h.SaveNote(ctx, n, id); err != nil {
		log.Print(err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	return api.NewResponse(http.StatusOK, output{ID: id})
}

func convert(input input) note.Note {
	return note.Note{
		Content:       []byte(input.Content),
		ReadPassword:  []byte(input.ReadPassword),
		WritePassword: []byte(input.WritePassword),
		TTL:           input.TTL,
	}
}

func NewHandler(saver saver, encrypter encrypter) *Handler {
	return &Handler{saver: saver, encrypter: encrypter}
}
