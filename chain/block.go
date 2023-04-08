package chain

import (
	"context"
	"gdxmonitor/service"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client

func init() {
	cl, err := ethclient.Dial(wsUrl)
	if err != nil {
		panic(err)
	}
	client = cl
}

var initNum = 75351123

func ParseBlock(num *big.Int) {
	header, err := client.HeaderByNumber(context.Background(), num)
	if err != nil {
		log.Fatal(err)
	}
	service.CreateBlock(*header)
}

func SyncBlock() {
	for {
		initNum++
		num := big.NewInt(int64(initNum))
		ParseBlock(num)
	}
}
