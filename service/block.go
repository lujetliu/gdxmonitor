package service

import (
	"gdxmonitor/global"
	"gdxmonitor/model"

	"github.com/ethereum/go-ethereum/core/types"
)

func CreateBlock(header types.Header) (*model.Block, error) {
	var block model.Block
	block.Number = header.Number.String()
	block.ParentHash = header.ParentHash.String()
	block.UncleHash = header.UncleHash.String()
	block.Coinbase = header.Coinbase.String()
	block.Root = header.Root.String()
	block.TxHash = header.TxHash.String()
	block.ReceiptHash = header.ReceiptHash.String()
	block.GasLimit = header.GasLimit
	block.GasUsed = header.GasUsed
	block.Time = header.Time

	err := block.Create(global.DBEngine)
	if err != nil {
		return &block, err
	}
	return &block, nil
}
