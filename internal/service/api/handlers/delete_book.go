package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/nft-books/book-svc/internal/service/api/helpers"
	"gitlab.com/tokend/nft-books/book-svc/internal/service/api/requests"
)

func DeleteBookByID(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetBookByIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	book, err := helpers.GetBookByID(r, req.ID)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if book == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	err = helpers.DB(r).Books().DeleteByID(req.ID)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusNoContent)
}
