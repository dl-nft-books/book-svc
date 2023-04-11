package helpers

import (
	"github.com/dl-nft-books/book-svc/solidity/generated/contractsregistry"
	"github.com/dl-nft-books/book-svc/solidity/generated/rolemanager"
	"github.com/dl-nft-books/network-svc/connector/models"
	"github.com/ethereum/go-ethereum/common"
)

func CheckMarketplacePerrmision(network models.NetworkDetailedResponse, address string) (bool, error) {
	contractRegistry, err := contractsregistry.NewContractsregistry(common.HexToAddress(network.FactoryAddress), network.RpcUrl)
	if err != nil {
		return false, err
	}
	roleManagerContract, err := contractRegistry.GetRoleManagerContract(nil)
	if err != nil {
		return false, err
	}
	roleManager, err := rolemanager.NewRolemanager(roleManagerContract, network.RpcUrl)
	if err != nil {
		return false, err
	}
	return roleManager.RolemanagerCaller.IsMarketplaceManager(nil, common.HexToAddress(address))
}