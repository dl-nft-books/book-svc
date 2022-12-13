package data

import (
	"time"

	"gitlab.com/tokend/nft-books/book-svc/resources"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type BookQ interface {
	New() BookQ

	Get() (*Book, error)
	Select() ([]Book, error)
	Count() (*uint64, error)

	Insert(data Book) (int64, error)
	DeleteByID(id int64) error
	Update(updater BookUpdateParams, id int64) error

	UpdateContractParams(name, symbol, price string, id int64) error

	UpdatePrice(price string, id int64) error
	UpdateContractName(name string, id int64) error
	UpdateContractAddress(newAddress string, id int64) error
	UpdateLastBlock(newLastBlock uint64, id int64) error
	UpdateSymbol(newSymbol string, id int64) error
	UpdateDeployStatus(newStatus resources.DeployStatus, id int64) error

	// FilterActual does not include deleted books
	FilterActual() BookQ
	FilterByID(id ...int64) BookQ
	FilterByTokenId(tokenId ...int64) BookQ
	FilterByDeployStatus(status ...resources.DeployStatus) BookQ
	FilterByContractAddress(address ...string) BookQ

	Page(params pgdb.OffsetPageParams) BookQ
}

// BookUpdateParams is a structure for applicable update params on bookQ `Update`
type BookUpdateParams struct {
	Banner       *string
	File         *string
	Title        *string
	Description  *string
	Contract     *string
	DeployStatus *resources.DeployStatus
	Symbol       *string
	Price        *string
}

type Book struct {
	ID              int64                  `db:"id" structs:"-"`
	Title           string                 `db:"title" structs:"title"`
	Description     string                 `db:"description" structs:"description"`
	CreatedAt       time.Time              `db:"created_at" structs:"created_at"`
	Price           string                 `db:"price" structs:"price"`
	ContractAddress string                 `db:"contract_address" structs:"contract_address"`
	ContractName    string                 `db:"contract_name" structs:"contract_name"`
	ContractSymbol  string                 `db:"contract_symbol" structs:"contract_symbol"`
	ContractVersion string                 `db:"contract_version" structs:"contract_version"`
	Banner          string                 `db:"banner" structs:"banner"`
	File            string                 `db:"file" structs:"file"`
	Deleted         bool                   `db:"deleted" structs:"-"`
	TokenId         int64                  `db:"token_id" structs:"token_id"`
	DeployStatus    resources.DeployStatus `db:"deploy_status" structs:"deploy_status"`
	LastBlock       uint64                 `db:"last_block" structs:"last_block"`
}
