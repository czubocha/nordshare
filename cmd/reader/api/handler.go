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

type (
	output struct {
		Content string `json:"content"`
		TTL     int64  `json:"ttl"`
	}
	Handler struct {
		reader
		decrypter
	}
	reader interface {
		ReadNote(context.Context, string) (nordshare.Note, error)
	}
	decrypter interface {
		Decrypt(*[]byte) error
	}
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest, h *Handler) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters[nordshare.IDPathParamName]
	password := api.GetHeaderIncasesensible(request.Headers, nordshare.PasswordHeaderName)
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
	if !hash.HasReadAccess(n, []byte(password)) {
		log.Print("reader: incorrect password")
		return api.NewResponse(http.StatusUnauthorized)
	}
	if err := h.Decrypt(&n.Content); err != nil {
		log.Printf("reader: %v", err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	return api.NewResponse(http.StatusOK, output{
		Content: string(n.Content),
		TTL:     n.TTL,
	})
}

func NewHandler(reader reader, decrypter decrypter) *Handler {
	return &Handler{reader: reader, decrypter: decrypter}
}
