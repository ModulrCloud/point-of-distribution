package databases

import (
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
)

type Stores struct {
	CoreBlocksData                     *leveldb.DB
	AnchorsCoreBlocksData              *leveldb.DB
	AggregatedLeaderFinalizationProofs *leveldb.DB
}

func Init(basePath string) (*Stores, error) {
	coreBlocksData, err := leveldb.OpenFile(filepath.Join(basePath, "core_blocks_data"), nil)
	if err != nil {
		return nil, err
	}

	anchorsBlocksData, err := leveldb.OpenFile(filepath.Join(basePath, "anchors_core_blocks_data"), nil)
	if err != nil {
		coreBlocksData.Close()
		return nil, err
	}

	aggregatedLeaderFinalizationProofs, err := leveldb.OpenFile(filepath.Join(basePath, "aggregated_leaders_finalization_proofs"), nil)
	if err != nil {
		coreBlocksData.Close()
		anchorsBlocksData.Close()
		return nil, err
	}

	return &Stores{CoreBlocksData: coreBlocksData, AnchorsCoreBlocksData: anchorsBlocksData, AggregatedLeaderFinalizationProofs: aggregatedLeaderFinalizationProofs}, nil
}

func (s *Stores) Close() {
	if s == nil {
		return
	}
	if s.CoreBlocksData != nil {
		s.CoreBlocksData.Close()
	}
	if s.AnchorsCoreBlocksData != nil {
		s.AnchorsCoreBlocksData.Close()
	}
	if s.AggregatedLeaderFinalizationProofs != nil {
		s.AggregatedLeaderFinalizationProofs.Close()
	}
}
