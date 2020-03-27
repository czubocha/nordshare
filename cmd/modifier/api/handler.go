package api

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github/czubocha/nordshare"
	"github/czubocha/nordshare/internal/api"
	"github/czubocha/nordshare/pkg/hash"
	"github/czubocha/nordshare/pkg/storage"
	"log"
	"net/http"
)

type (
	input struct {
		Content string `json:"content"`
		TTL     int64  `json:"ttl"`
	}
	Handler struct {
		modifier
		encrypter
	}
	modifier interface {
		ReadNote(context.Context, string) (nordshare.Note, error)
		UpdateNote(context.Context, []byte, int64, string) error
	}
	encrypter interface {
		Encrypt(*[]byte) error
	}
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest, h *Handler) (events.APIGatewayProxyResponse, error) {
	var in input
	id := request.PathParameters[nordshare.IDPathParamName]
	password := api.GetHeaderIncasesensible(request.Headers, nordshare.PasswordHeaderName)
	if err := json.Unmarshal([]byte(request.Body), &in); err != nil {
		log.Printf("modifier: %v", err)
		return api.NewResponse(http.StatusBadRequest)
	}
	n, err := h.ReadNote(ctx, id)
	if err != nil {
		log.Printf("reader: %v", err)
		switch err {
		case storage.ErrNoteExpired:
			// return status 401 to hide IDs of existing notes
			return api.NewResponse(http.StatusUnauthorized)
		default:
			return api.NewResponse(http.StatusInternalServerError)
		}
	}
	if !hash.HasWriteAccess(n, []byte(password)) {
		log.Print("modifier: incorrect password")
		return api.NewResponse(http.StatusUnauthorized)
	}
	content := []byte(in.Content)
	if err := h.Encrypt(&content); err != nil {
		log.Printf("modifier: %v", err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	if err := h.UpdateNote(ctx, content, in.TTL, id); err != nil {
		log.Printf("modifier: %v", err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	return api.NewResponse(http.StatusOK)
}

func NewHandler(saver modifier, encrypter encrypter) *Handler {
	return &Handler{modifier: saver, encrypter: encrypter}
}
