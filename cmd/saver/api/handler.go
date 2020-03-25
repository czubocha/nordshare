package api

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github/czubocha/nordshare"
	"github/czubocha/nordshare/internal/api"
	"github/czubocha/nordshare/pkg/hash"
	"log"
	"net/http"
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
		SaveNote(context.Context, nordshare.Note, string) error
	}
	encrypter interface {
		Encrypt(*[]byte) error
	}
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest, h *Handler) (events.APIGatewayProxyResponse, error) {
	var in input
	if err := json.Unmarshal([]byte(request.Body), &in); err != nil {
		log.Printf("saver: %v", err)
		return api.NewResponse(http.StatusBadRequest)
	}
	n := convert(in)
	if err := h.Encrypt(&n.Content); err != nil {
		log.Printf("saver: %v", err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	if err := hash.Passwords(&n); err != nil {
		log.Printf("saver: %v", err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	id := request.RequestContext.RequestID
	if err := h.SaveNote(ctx, n, id); err != nil {
		log.Printf("saver: %v", err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	return api.NewResponse(http.StatusCreated, output{ID: id})
}

func convert(input input) nordshare.Note {
	return nordshare.Note{
		Content:       []byte(input.Content),
		ReadPassword:  []byte(input.ReadPassword),
		WritePassword: []byte(input.WritePassword),
		TTL:           input.TTL,
	}
}

func NewHandler(saver saver, encrypter encrypter) *Handler {
	return &Handler{saver: saver, encrypter: encrypter}
}
