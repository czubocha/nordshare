package api

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github/czubocha/nordshare"
	"github/czubocha/nordshare/internal/api"
	"github/czubocha/nordshare/pkg/hash"
	"github/czubocha/nordshare/pkg/storage"
	"log"
	"net/http"
)

const (
	idPathParamName    = "id"
	passwordHeaderName = "password"
)

type (
	Handler struct {
		remover
		decrypter
	}
	remover interface {
		ReadNote(context.Context, string) (nordshare.Note, error)
		DeleteNote(context.Context, string) error
	}
	decrypter interface {
		Decrypt(*nordshare.Note) error
	}
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest, h *Handler) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters[idPathParamName]
	password := request.Headers[passwordHeaderName]
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
		log.Print("remover: incorrect password")
		return api.NewResponse(http.StatusUnauthorized)
	}
	if err = h.DeleteNote(ctx, id); err != nil {
		log.Printf("remover: %v", err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	return api.NewResponse(http.StatusOK)
}

func NewHandler(remover remover) *Handler {
	return &Handler{remover: remover}
}
