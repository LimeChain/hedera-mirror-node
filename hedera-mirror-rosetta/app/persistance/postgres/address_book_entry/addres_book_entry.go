package address_book_entry

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/jinzhu/gorm"
	"strconv"
)

type addressBookEntry struct {
	Id                 int32  `gorm:"type:integer;primary_key"`
	ConsensusTimestamp int64  `gorm:"type:bigint"`
	Ip                 string `gorm:"size:128"`
	Port               int32  `gorm:"type:integer"`
	Memo               string `gorm:"size:128"`
	PublicKey          string `gorm:"size:1024"`
	NodeId             *int   `gorm:"type:bigint"`
	NodeAccountId      int64  `gorm:"type:bigint"`
	NodeCertHash       []byte `gorm:"type:bytea"`
}

func (addressBookEntry) TableName() string {
	return "address_book_entry"
}

func (abe *addressBookEntry) GetPeerId() string {
	if abe.NodeId == nil {
		return abe.Memo
	}

	return strconv.Itoa(*abe.NodeId)
}

// AddressBookEntryRepository struct that has connection to the Database
type AddressBookEntryRepository struct {
	dbClient *gorm.DB
}

// Entries return all found Address Book Entries
func (aber *AddressBookEntryRepository) Entries() (*types.AddressBookEntries, *rTypes.Error) {
	dbEntries := aber.retrieveEntries()

	entries := make([]*types.AddressBookEntry, len(dbEntries))
	for i, e := range dbEntries {
		entries[i] = &types.AddressBookEntry{
			PeerId: e.GetPeerId(),
			Metadata: map[string]interface{}{
				"ip":   e.Ip,
				"port": e.Port,
			},
		}
	}

	return &types.AddressBookEntries{
		Entries: entries}, nil
}

func (aber *AddressBookEntryRepository) retrieveEntries() []addressBookEntry {
	var entries []addressBookEntry
	aber.dbClient.Find(&entries)
	return entries
}

// NewAddressBookEntryRepository creates an instance of a AddressBookEntryRepository struct.
func NewAddressBookEntryRepository(dbClient *gorm.DB) *AddressBookEntryRepository {
	return &AddressBookEntryRepository{
		dbClient: dbClient,
	}
}
