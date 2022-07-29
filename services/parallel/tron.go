package parallel

import (
	"fmt"
	"github.com/darwinia-network/token/lib/web3"
	"github.com/darwinia-network/token/util"
	"math/big"
)

// var tronContract = "416e0d26adf5323f5b82d5714354dc3c6870adee7c"

type TronResponse struct {
	ConstantResult []string `json:"constant_result"`
}

func RingTronSupply(contract string) *big.Int {
	w := web3.New("tron")
	var e TronResponse

	if _ = w.Call(&e, fmt.Sprintf(util.TrxBase58toHexAddress(contract)), "totalSupply()"); len(e.ConstantResult) > 0 {
		return util.U256(e.ConstantResult[0])
	}
	return big.NewInt(0)
}

func RingTronBalance(contract, address string) *big.Int {
	w := web3.New("tron")
	var e TronResponse
	if _ = w.Call(&e, fmt.Sprintf(util.TrxBase58toHexAddress(contract)), "balanceOf(address)", util.TrimTronHex(address)); len(e.ConstantResult) > 0 {
		return util.U256(e.ConstantResult[0])
	}
	return big.NewInt(0)
}

type TronScan struct {
	Success bool             `json:"success"`
	Data    []TronScanResult `json:"data"`
}

type TronScanResult struct {
	TransactionId  string            `json:"transaction_id"`
	EventName      string            `json:"event_name"`
	Result         map[string]string `json:"result"`
	BlockNumber    int               `json:"block_number"`
	BlockTimestamp int64             `json:"block_timestamp"`
}

func TronScanLog(start int64, address string) (*TronScan, error) {
	w := web3.New("tron")
	var e TronScan
	if err := w.Event(&e, start, 0, address); err != nil || !e.Success {
		return nil, err
	}
	return &e, nil
}
