package chain

import (
	"context"
	"fmt"
	"gdxmonitor/util"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

const (
	PlaceMakerOrderId  = "0xd0caf2c4677b4d382504eb6b0f15d030d6887b69dbcb3fc302a62fea5a9fc7a7"
	PlaceOrderTypeSell = "sell"
	PlaceOrderTypeBuy  = "buy"
)

var gdxAddress = "0x8Eb76679F7eD2a2Ec0145A87fE35d67ff6e19aa6"

// logPlaceMakerOrder := []byte("PlaceMakerOrder(uint256, address, uint64, bool, int24, uint128)")
// logSettleMakerOrder := []byte("SettleMakerOrder(uint256, uint128, uint128, uint128)")
// logPlaceOrderHash := crypto.Keccak256Hash(logPlaceMakerOrder)
// logSettleOrderHash := crypto.Keccak256Hash(logSettleMakerOrder)

func SubEvent() {
	client, err := ethclient.Dial(wsUrl)
	if err != nil {
		log.Fatal(err)
	}

	gdxContract := common.HexToAddress(gdxAddress)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{gdxContract},
	}

	contractAbi, err := abi.JSON(strings.NewReader(GdxABI))
	if err != nil {
		log.Fatal(err)
	}

	logs := make(chan types.Log)
	_, err = client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for vLog := range logs {
		// fmt.Println(vLog.BlockHash.Hex()) // 0x3404b8c050aa0aacd0223e91b5c32fee6400f357764771d0684fa7b3f448f1a8
		// fmt.Println(vLog.BlockNumber)     // 2394201

		// event := struct {
		// 	OrderId        *big.Int
		// 	Recipient      common.Address
		// 	BundleId       uint64
		// 	Zero           bool
		// 	BboundaryLower *big.Int
		// 	Amount         *big.Int
		// }{}

		sign := vLog.Topics[0].Hex()
		// 挂单
		if sign != PlaceMakerOrderId {
			continue
		}

		var pmod PlaceMakerOrderData
		pmod.TxHash = vLog.TxHash.Hex()
		pmod.MethodId = PlaceMakerOrderId
		// fmt.Println(vLog.TxHash.Hex()) // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6

		orderId := vLog.Topics[1].Big()
		recipient := vLog.Topics[2].String()
		bundleId := vLog.Topics[3].Big()
		address := common.HexToAddress(recipient)

		pmod.OrderId = orderId
		pmod.Recipient = address
		pmod.BundleId = bundleId

		var in = make(map[string]interface{})
		err = contractAbi.UnpackIntoMap(in, "PlaceMakerOrder", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range in {
			if k == "amount" {
				amount := v.(*big.Int)
				wei := new(big.Int)
				wei.SetString(amount.String(), 10)
				decimalAmount := util.ToDecimal(wei, 18)

				pmod.Amount = decimalAmount
				continue
			}
			if k == "boundaryLower" {
				boundaryLower := v.(*big.Int)
				pmod.BoundaryLower = boundaryLower
			}
			if k == "zero" {
				zero := v.(bool)
				pmod.Zero = zero
			}
		}

		if pmod.Zero {
			pmod.PlaceType = PlaceOrderTypeSell
			pmod.Token = "GDX"
		} else {
			pmod.PlaceType = PlaceOrderTypeBuy
			pmod.Token = "ETH"
		}

		pmod.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(pmod.Display())

	}
}

type PlaceMakerOrderData struct {
	Timestamp     string
	MethodId      string
	OrderId       *big.Int
	Recipient     common.Address
	BundleId      *big.Int
	Zero          bool
	BoundaryLower *big.Int
	Amount        decimal.Decimal
	Token         string
	TxHash        string
	PlaceType     string
}

func (p PlaceMakerOrderData) Display() string {
	amount, _ := p.Amount.Float64()
	return fmt.Sprintf("%s: address %s %s %.4f %s", p.Timestamp, p.Recipient, p.PlaceType, amount, p.Token)
}
