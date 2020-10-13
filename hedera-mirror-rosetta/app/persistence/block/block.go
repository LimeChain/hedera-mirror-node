/*-
 * ‌
 * Hedera Mirror Node
 *
 * Copyright (C) 2019 - 2020 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ‍
 */

package block

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/jinzhu/gorm"
	"log"
)

const genesisPreviousHash = "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

const (
	// selectLatestWithIndex - Selects row with the latest consensus_end and adds additional info about the position of that row using count.
	// The information about the position is used as Block Index
	selectLatestWithIndex string = `SELECT rd.file_hash,
                                           rd.consensus_start,
                                           rd.consensus_end,
                                           rd.prev_hash,
                                           rcd_index.block_index
                                    FROM   (SELECT *
                                            FROM   record_file
                                            WHERE  consensus_end = (SELECT MAX(consensus_end)
                                                                    FROM   record_file)) AS rd,
                                           (SELECT COUNT(*) - 1 - ? AS block_index
                                            FROM   record_file) AS rcd_index`

	// selectByHashWithIndex - Selects the row with a given file_hash and adds additional info about the position of that row using count.
	//The information about the position is used as Block Index
	selectByHashWithIndex string = `SELECT rd.file_hash,
                                           rd.consensus_start,
                                           rd.consensus_end,
                                           rd.prev_hash,
                                           rcd.block_index
                                    FROM   (SELECT *
                                            FROM   record_file
                                            WHERE  file_hash = ?) AS rd,
                                           (SELECT Count(*) - 1 AS block_index
                                            FROM   record_file
                                            WHERE  consensus_end <= (SELECT consensus_end
                                                                     FROM   record_file
                                                                     WHERE  file_hash = ?)) AS rcd`
	// selectSkippedRecordFilesCount - Selects the count of rows from the record_file table,
	// where each one's consensus_end is before the MIN consensus_end of address_book table (the first one added).
	// This way, record files before that timestamp are considered non-existent,
	// and the first record_file (block) will be considered equal or bigger
	// than the consensus_timestamp of the first account_balance
	selectSkippedRecordFilesCount string = `SELECT COUNT(*)
                                            FROM record_file
                                            WHERE consensus_end < (SELECT MIN(consensus_timestamp) FROM account_balance)`
)

type recordFile struct {
	ID             int64  `gorm:"type:bigint;primary_key"`
	Name           string `gorm:"size:250"`
	LoadStart      int64  `gorm:"type:bigint"`
	LoadEnd        int64  `gorm:"type:bigint"`
	FileHash       string `gorm:"size:96"`
	PrevHash       string `gorm:"size:96"`
	ConsensusStart int64  `gorm:"type:bigint"`
	ConsensusEnd   int64  `gorm:"type:bigint"`
	BlockIndex     int64  `gorm:"type:bigint"`
}

// TableName - Set table name to be `record_file`
func (recordFile) TableName() string {
	return "record_file"
}

// BlockRepository struct that has connection to the Database
type BlockRepository struct {
	dbClient                *gorm.DB
	recordFileStartingIndex *int64
}

// NewBlockRepository creates an instance of a BlockRepository struct
func NewBlockRepository(dbClient *gorm.DB) *BlockRepository {
	return &BlockRepository{dbClient: dbClient}
}

// FindByIndex retrieves a block by given Index
func (br *BlockRepository) FindByIndex(index int64) (*types.Block, *rTypes.Error) {
	startingIndex := br.getRecordFilesStartingIndex()

	rf := &recordFile{}
	if br.dbClient.Order("consensus_end asc").Offset(index + *startingIndex).First(rf).RecordNotFound() {
		return nil, errors.Errors[errors.BlockNotFound]
	}

	return br.constructBlockResponse(rf, index), nil
}

// FindByHash retrieves a block by a given Hash
func (br *BlockRepository) FindByHash(hash string) (*types.Block, *rTypes.Error) {
	rf, err := br.findRecordFileByHash(hash)
	if err != nil {
		return nil, err
	}

	return br.constructBlockResponse(rf, rf.BlockIndex), nil
}

// FindByIdentifier retrieves a block by Index && Hash
func (br *BlockRepository) FindByIdentifier(index int64, hash string) (*types.Block, *rTypes.Error) {
	rf, err := br.findRecordFileByHash(hash)
	if err != nil {
		return nil, err
	}
	if rf.BlockIndex != index {
		return nil, errors.Errors[errors.BlockNotFound]
	}
	return br.constructBlockResponse(rf, index), nil
}

// RetrieveGenesis retrieves the genesis block
func (br *BlockRepository) RetrieveGenesis() (*types.Block, *rTypes.Error) {
	startingIndex := br.getRecordFilesStartingIndex()

	rf := &recordFile{}
	if br.dbClient.Offset(*startingIndex).Limit(1).Find(rf).RecordNotFound() {
		return nil, errors.Errors[errors.BlockNotFound]
	}

	return &types.Block{
		Index:               0,
		Hash:                rf.FileHash,
		ConsensusStartNanos: rf.ConsensusStart,
		ConsensusEndNanos:   rf.ConsensusEnd,
	}, nil
}

// RetrieveLatest retrieves the latest block
func (br *BlockRepository) RetrieveLatest() (*types.Block, *rTypes.Error) {
	startingIndex := br.getRecordFilesStartingIndex()

	rf := &recordFile{}
	if br.dbClient.Raw(selectLatestWithIndex, *startingIndex).Scan(rf).RecordNotFound() {
		return nil, errors.Errors[errors.BlockNotFound]
	}

	return br.constructBlockResponse(rf, rf.BlockIndex), nil
}

func (br *BlockRepository) findRecordFileByHash(hash string) (*recordFile, *rTypes.Error) {
	startingIndex := br.getRecordFilesStartingIndex()

	rf := &recordFile{}
	if br.dbClient.Raw(selectByHashWithIndex, hash, hash).Scan(rf).RecordNotFound() {
		return nil, errors.Errors[errors.BlockNotFound]
	}

	rf.BlockIndex = rf.BlockIndex - *startingIndex
	if rf.BlockIndex < 0 {
		return nil, errors.Errors[errors.BlockNotFound]
	}

	return rf, nil
}

// constructBlockResponse returns the constructed Block. Takes into account genesis block. Block index is passed separately due to FindByIndex
func (br *BlockRepository) constructBlockResponse(rf *recordFile, blockIndex int64) *types.Block {
	parentIndex := blockIndex - 1
	parentHash := rf.PrevHash

	// Handle the edge case for querying first block
	if rf.PrevHash == genesisPreviousHash || parentIndex < 0 {
		parentIndex = 0          // Parent index should be 0, same as current block index
		parentHash = rf.FileHash // Parent hash should be same as current block hash
	}
	return &types.Block{
		Index:               blockIndex,
		Hash:                rf.FileHash,
		ParentIndex:         parentIndex,
		ParentHash:          parentHash,
		ConsensusStartNanos: rf.ConsensusStart,
		ConsensusEndNanos:   rf.ConsensusEnd,
	}
}

func (br *BlockRepository) getRecordFilesStartingIndex() *int64 {
	if br.recordFileStartingIndex == nil {
		br.dbClient.Raw(selectSkippedRecordFilesCount).Count(&br.recordFileStartingIndex)
		log.Printf(`Fetched Record Files starting index: %d`, *br.recordFileStartingIndex)
	}
	return br.recordFileStartingIndex
}
