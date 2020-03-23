package api

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"nordshare/internal/api"
	"nordshare/pkg/hash"
	"nordshare/pkg/note"
	"nordshare/pkg/storage"
)

const (
	idPathParamName    = "id"
	passwordHeaderName = "password"
)

type (
	input struct {
		Content       string `json:"content"`
		ReadPassword  string `json:"readPassword"`
		WritePassword string `json:"writePassword"`
		TTL           int64  `json:"ttl"`
	}
	output struct {
		Content string `json:"content"`
		TTL     int64  `json:"ttl"`
	}
	Handler struct {
		reader
		decrypter
	}
	reader interface {
		ReadNote(context.Context, string) (note.Note, error)
	}
	decrypter interface {
		DecryptContent(*note.Note) error
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
	if hash.HasReadAccess(n, []byte(password)) == false {
		log.Print("reader: incorrect password")
		return api.NewResponse(http.StatusUnauthorized)
	}
	if err := h.DecryptContent(&n); err != nil {
		log.Printf("reader: %v", err)
		return api.NewResponse(http.StatusInternalServerError)
	}
	return api.NewResponse(http.StatusOK, output{
		Content: string(n.Content),
		TTL:     n.TTL,
	})
}

func NewHandler(saver reader, decrypter decrypter) *Handler {
	return &Handler{reader: saver, decrypter: decrypter}
}
