package genutil

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

type InitConfig struct {
	ChainID   string
	GenTxsDir string
	Name      string
	NodeID    string
	ValPubKey cryptotypes.PubKey
}

// NewInitConfig creates a new InitConfig object
func NewInitConfig(chainID, genTxsDir, name, nodeID string, valPubKey cryptotypes.PubKey) InitConfig {
	return InitConfig{
		ChainID:   chainID,
		GenTxsDir: genTxsDir,
		Name:      name,
		NodeID:    nodeID,
		ValPubKey: valPubKey,
	}
}
