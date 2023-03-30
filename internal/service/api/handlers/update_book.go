package handlers

import (
	"errors"
	"fmt"
	"github.com/dl-nft-books/book-svc/solidity/generated/contractsregistry"
	"github.com/dl-nft-books/book-svc/solidity/generated/rolemanager"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"

	"github.com/dl-nft-books/book-svc/internal/data"
	"github.com/dl-nft-books/book-svc/internal/service/api/helpers"
	"github.com/dl-nft-books/book-svc/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func UpdateBookByID(w http.ResponseWriter, r *http.Request) {

	logger := helpers.Log(r)
	networker := helpers.Networker(r)
	request, err := requests.NewUpdateBookRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	address := r.Context().Value("address").(string)
	bookData, err := helpers.GetBookByID(r, requests.GetBookByIDRequest{
		ID: request.ID,
	})
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if bookData == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}
	book, err := helpers.NewBook(bookData)
	for _, networkData := range book.Attributes.Networks {
		network, err := networker.GetNetworkDetailedByChainID(networkData.Attributes.ChainId)
		if err != nil {
			logger.WithError(err).Error("default failed to check if network exists")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		contractRegistry, err := contractsregistry.NewContractsregistry(common.HexToAddress(network.FactoryAddress), network.RpcUrl)
		if err != nil {
			logger.WithError(err).Debug("failed to create contract registry")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		roleManagerContract, err := contractRegistry.GetRoleManagerContract(nil)
		if err != nil {
			logger.WithError(err).Debug("failed to get role manager contract")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		roleManager, err := rolemanager.NewRolemanager(roleManagerContract, network.RpcUrl)
		if err != nil {
			logger.WithError(err).Debug("failed to create role manager")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		isAdmin, err := roleManager.RolemanagerCaller.IsAdmin(nil, common.HexToAddress(address))
		if err != nil {
			logger.WithError(err).Debug("failed to check is admin")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if !isAdmin {
			isManager, err := roleManager.RolemanagerCaller.IsAdmin(nil, common.HexToAddress(address))
			if err != nil {
				logger.WithError(err).Debug("failed to check is admin")
				ape.RenderErr(w, problems.InternalError())
				return
			}
			if !isManager {
				logger.WithFields(logan.F{"account": address}).Debug("you don't have permission to create book")
				ape.RenderErr(w, problems.Forbidden())
				return
			}
		}
	}
	updateParams := data.BookUpdateParams{
		Banner:      nil,
		File:        nil,
		Description: nil,
	}

	// Collecting update params
	banner := request.Data.Attributes.Banner
	if banner != nil {
		if err = helpers.CheckBannerMimeType(banner.Attributes.MimeType, r); err != nil {
			helpers.Log(r).WithError(err).Error("failed to validate banner mime type")
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}

		if err = helpers.SetMediaLink(r, banner); err != nil {
			helpers.Log(r).WithError(err).Error("failed to set banner link")
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}

		bannerMediaRaw := helpers.MarshalMedia(banner)
		updateParams.Banner = &bannerMediaRaw[0]
	}

	file := request.Data.Attributes.File
	if file != nil {
		if err = helpers.CheckFileMimeType(file.Attributes.MimeType, r); err != nil {
			helpers.Log(r).WithError(err).Error("failed to validate file mime type")
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}

		if err = helpers.SetMediaLink(r, file); err != nil {
			helpers.Log(r).WithError(err).Error("failed to set file link")
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}

		fileMediaRaw := helpers.MarshalMedia(file)
		updateParams.File = &fileMediaRaw[0]
	}

	description := request.Data.Attributes.Description
	if description != nil {
		if len(*description) > requests.MaxDescriptionLength {
			err = errors.New(fmt.Sprintf("invalid description length (max len is %v)", requests.MaxDescriptionLength))
			helpers.Log(r).WithError(err).Error("failed to validate book description")
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}

		updateParams.Description = description
	}

	// Updating collected params
	if err = helpers.DB(r).Books().Update(updateParams, request.ID); err != nil {
		helpers.Log(r).WithError(err).Error("failed to update book params")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if request.Data.Attributes.Network != nil {
		if err = helpers.DB(r).Books().UpdateDeployStatus(
			request.Data.Attributes.Network.Attributes.DeployStatus,
			request.ID,
			request.Data.Attributes.Network.Attributes.ChainId); err != nil {
			helpers.Log(r).WithError(err).Error("failed to update book deploy status param")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}
	ape.Render(w, http.StatusNoContent)
}
