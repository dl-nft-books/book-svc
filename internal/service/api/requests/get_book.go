package requests

import (
	"gitlab.com/distributed_lab/urlval"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type GetBookByIDRequest struct {
	ID      int64   `json:"id"`
	ChainId []int64 `filter:"chain_id"`
}

func NewGetBookByIDRequest(r *http.Request) (GetBookByIDRequest, error) {
	var request GetBookByIDRequest

	id := chi.URLParam(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		return request, errors.New("id is not an integer")
	}
	request.ID = cast.ToInt64(id)
	return request, urlval.Decode(r.URL.Query(), &request)
}
