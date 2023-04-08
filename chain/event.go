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
	SettleMakerOrderId = "0xa9d41f4c7e5cdf552e9bfe6d10327a427231e7905304c308dbf7455b6905556f"
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
		methodId := vLog.Topics[0].Hex()
		// 挂单
		if methodId == PlaceMakerOrderId {
			ParsePlaceMakerOrder(contractAbi, vLog)
		}

		// 结算
		if methodId == SettleMakerOrderId {
			ParseSettleMakerOrder(contractAbi, vLog)
		}
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

type SettleMakerOrderData struct {
	Timestamp         string
	MethodId          string
	OrderId           *big.Int
	MakerAmountOut    decimal.Decimal
	TakerAmountOut    decimal.Decimal
	TakerFeeAmountOut decimal.Decimal
	Token             string
	TxHash            string
	PlaceType         string
}

func (p PlaceMakerOrderData) Display() string {
	amount, _ := p.Amount.Float64()
	return fmt.Sprintf("%s: address %s %s %.4f %s", p.Timestamp, p.Recipient, p.PlaceType, amount, p.Token)
}

func ParsePlaceMakerOrder(contractAbi abi.ABI, vLog types.Log) error {
	var pmod PlaceMakerOrderData
	pmod.TxHash = vLog.TxHash.Hex()
	pmod.MethodId = PlaceMakerOrderId
	// fmt.Println(vLog.TxHash.Hex()) // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6

	orderId := vLog.Topics[1].Big()
	recipient := vLog.Topics[2].String()
	bundleId := vLog.Topics[3].Big()

	address := common.HexToAddress(recipient)
	if address.String() != "0xD35881A7DBE237F4a96e26b06C7264641970fcC4" && address.String() != "0x9020BED433ACeDD6afc1933d439Ea19ca33ae646" &&
		address.String() != "0x2B68c406306326633CE0Db0d7771A3904cEC31a7" && address.String() != "0x7299542910AaCf1F548343d2AC49236e44630cb9" {
		return nil
	}

	pmod.OrderId = orderId
	pmod.Recipient = address
	pmod.BundleId = bundleId

	var in = make(map[string]interface{})
	err := contractAbi.UnpackIntoMap(in, "PlaceMakerOrder", vLog.Data)
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
	return nil
}

func ParseSettleMakerOrder(contractAbi abi.ABI, vLog types.Log) error {
	var smod SettleMakerOrderData
	smod.TxHash = vLog.TxHash.Hex()
	smod.MethodId = SettleMakerOrderId
	// fmt.Println(vLog.TxHash.Hex()) // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6

	orderId := vLog.Topics[1].Big()
	smod.OrderId = orderId

	var in = make(map[string]interface{})
	err := contractAbi.UnpackIntoMap(in, "SettleMakerOrder", vLog.Data)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range in {
		amount := v.(*big.Int)
		wei := new(big.Int)
		wei.SetString(amount.String(), 10)
		decimalAmount := util.ToDecimal(wei, 18)

		switch k {
		case "makerAmountOut":
			smod.MakerAmountOut = decimalAmount
		case "takerAmountOut":
			smod.TakerAmountOut = decimalAmount
		case "takerFeeAmountOut":
			smod.TakerFeeAmountOut = decimalAmount
		default:
			log.Fatal("err log data")
		}

	}

	smod.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	// fmt.Printf("%#v\n", smod)
	return nil
}

// func ParseSwap(contractAbi abi.ABI, vLog types.Log) error {
// 	var smod SettleMakerOrderData
// 	smod.TxHash = vLog.TxHash.Hex()
// 	smod.MethodId = SettleMakerOrderId
// 	// fmt.Println(vLog.TxHash.Hex()) // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6

// 	orderId := vLog.Topics[1].Big()
// 	smod.OrderId = orderId

// 	var in = make(map[string]interface{})
// 	err := contractAbi.UnpackIntoMap(in, "SettleMakerOrder", vLog.Data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for k, v := range in {
// 		amount := v.(*big.Int)
// 		wei := new(big.Int)
// 		wei.SetString(amount.String(), 10)
// 		decimalAmount := util.ToDecimal(wei, 18)

// 		switch k {
// 		case "makerAmountOut":
// 			smod.MakerAmountOut = decimalAmount
// 		case "takerAmountOut":
// 			smod.TakerAmountOut = decimalAmount
// 		case "takerFeeAmountOut":
// 			smod.TakerFeeAmountOut = decimalAmount
// 		default:
// 			log.Fatal("err log data")
// 		}

// 	}

// 	smod.Timestamp = time.Now().Format("2006-01-02 15:04:05")
// 	fmt.Printf("%#v\n", smod)
// 	return nil
// }
