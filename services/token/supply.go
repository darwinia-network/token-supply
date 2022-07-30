package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/darwinia-network/token/config"
	"github.com/darwinia-network/token/services/parallel"
	"github.com/darwinia-network/token/util"
	"github.com/shopspring/decimal"
	"strings"
	"sync"
)

// note: Due to unstable network may be failed to fetch token balance
// so returning last successful supplied balance to front-end
var latestRingSupply Supply
var latestKtonSupply Supply

type Supply struct {
	CirculatingSupply decimal.Decimal `json:"circulatingSupply"`
	TotalSupply       decimal.Decimal `json:"totalSupply"`
	BondLockBalance   decimal.Decimal `json:"bond_lock_balance"`
	TreasuryBalance   decimal.Decimal `json:"treasury_balance"`
	MaxSupply         decimal.Decimal `json:"maxSupply"`
	Details           []*SupplyDetail `json:"details"`
}

type SupplyDetail struct {
	Network           string          `json:"network"`
	CirculatingSupply decimal.Decimal `json:"circulatingSupply"`
	TotalSupply       decimal.Decimal `json:"totalSupply"`
	Precision         int             `json:"precision"`
	Type              string          `json:"type,omitempty"`
	Contract          string          `json:"contract,omitempty"`
}

type Currency struct {
	Code          string
	EthContract   string
	TronContract  string
	MaxSupply     decimal.Decimal
	FilterAddress map[string][]string
}

func RingSupply() *Supply {
	ring := Currency{
		Code:         "ring",
		EthContract:  config.Cfg.Ring,
		TronContract: config.Cfg.TronRing,
		MaxSupply:    decimal.New(1, 10),
	}
	ring.FilterAddress = map[string][]string{
		"Tron":     {"TDWzV6W1L1uRcJzgg2uKa992nAReuDojfQ", "TSu1fQKFkTv95U312R6E94RMdixsupBZDS", "TTW2Vpr9TCu6gxGZ1yjwqy7R79hEH8iscC"},
		"Ethereum": {"0x5FD8bCC6180eCd977813465bDd0A76A5a9F88B47", "0xfA4FE04f69F87859fCB31dF3B9469f4E6447921c", "0x7f23e4a473db3d11d11b43d90b34f8a778753e34", "0x7f23e4a473db3d11d11b43d90b34f8a778753e34"},
	}
	supply, errFlag := ring.supply()
	if errFlag != false{
		return &latestRingSupply
	}
	latestRingSupply = *supply
	return  &latestRingSupply
}

func KtonSupply() *Supply {
	kton := Currency{
		Code:         "kton",
		EthContract:  config.Cfg.Kton,
		TronContract: config.Cfg.TronKton,
	}
	supply, errFlag := kton.supply()
	if errFlag != false{
		return &latestKtonSupply
	}
	latestKtonSupply = *supply
	return  &latestKtonSupply
}

func (c *Currency) supply() (*Supply, bool) {
	var supply Supply
	supply.MaxSupply = c.MaxSupply // 10 billion
	wg := sync.WaitGroup{}
	wg.Add(4)
	errflag := false
	go func() {
		ethSupply, er := c.ethSupply()
		if er != nil{
			errflag = true
			wg.Done()
			return
		}
		supply.CirculatingSupply = supply.CirculatingSupply.Add(ethSupply.CirculatingSupply)
		supply.Details = append(supply.Details, ethSupply)
		wg.Done()
	}()
	go func() {
		tronSupply := c.tronSupply()
		if tronSupply.CirculatingSupply.GreaterThan(decimal.NewFromInt(0)){
			supply.CirculatingSupply = supply.CirculatingSupply.Add(tronSupply.CirculatingSupply)
			supply.Details = append(supply.Details, tronSupply)
		}
		wg.Done()
	}()
	go func() {
		var err error
		supply.TreasuryBalance, err = c.TreasuryBalance(100, 0, "system")
		if err != nil{
			errflag = true
		}
		wg.Done()
	}()
	go func() {
		var err error
		supply.TotalSupply, supply.BondLockBalance, err = c.TotalSupply()
		if err != nil{
			errflag = true
		}
		wg.Done()
	}()
	wg.Wait()

	if supply.MaxSupply.IsZero() {
		if c.Code == "kton" {
			supply.MaxSupply = supply.TotalSupply
		} else {
			for _, one := range supply.Details {
				supply.MaxSupply = supply.MaxSupply.Add(one.TotalSupply)
			}
		}
	}

	supply.CirculatingSupply = supply.TotalSupply.Sub(supply.BondLockBalance).Sub(supply.TreasuryBalance).
		Sub(supply.CirculatingSupply)
	if  supply.CirculatingSupply.LessThan(decimal.NewFromInt(0)){
		errflag = false
	}
	return &supply, errflag
}

func (c *Currency) ethSupply() (*SupplyDetail, error) {
	var supply SupplyDetail
	supply.Precision = 18
	precision := decimal.New(1, int32(supply.Precision))
	s, err := parallel.RingEthSupply(c.EthContract)
	if err != nil{
		return nil, err
	}
	capDecimal := decimal.NewFromBigInt(s, 0).Div(precision)
	supply.Network = "Ethereum"
	supply.Contract = c.EthContract
	supply.CirculatingSupply = capDecimal.Sub(supply.filterBalance(c.FilterAddress).Div(precision))
	supply.TotalSupply = capDecimal
	supply.Type = "erc20"

	return &supply, nil
}

func (c *Currency) tronSupply() *SupplyDetail {
	var supply SupplyDetail
	supply.Precision = 18
	precision := decimal.New(1, int32(supply.Precision))
	capDecimal := decimal.NewFromBigInt(parallel.RingTronSupply(c.TronContract), 0).Div(precision)
	supply.Contract = c.TronContract
	supply.Network = "Tron"
	supply.CirculatingSupply = capDecimal.Sub(supply.filterBalance(c.FilterAddress).Div(precision))
	supply.TotalSupply = capDecimal
	supply.Type = "trc20"

	return &supply
}

func (c *Currency) TreasuryBalance(pageSize, pageIndex int64, filter string) (decimal.Decimal, error) {
	type AccountDetail struct {
		Balance     decimal.Decimal `json:"balance"`
		BalanceLock decimal.Decimal `json:"balance_lock"`
		KtonBalance decimal.Decimal `json:"kton_balance"`
		KtonLock    decimal.Decimal `json:"kton_lock"`
	}
	type AccountTokenRes struct {
		Data struct {
			Count int             `json:"count"`
			List  []AccountDetail `json:"list"`
		} `json:"data"`
	}

	params := make(map[string]interface{})
	params["row"] = pageSize
	params["page"] = pageIndex
	params["filter"] = filter

	b, _ := json.Marshal(params)
	var res AccountTokenRes
	data, _ := util.PostWithJson(fmt.Sprintf("%s/api/scan/accounts", config.Cfg.SubscanHost), bytes.NewReader(b))
	err := util.UnmarshalAny(&res, data)
	if err != nil{
		return decimal.Decimal{}, err
	}
	var token decimal.Decimal

	for _, a := range res.Data.List {
		if c.Code == "ring" {
			token = token.Add(a.Balance).Add(a.BalanceLock)
		}
		// kton has not treasure
	}
	return token, nil
}

func (c *Currency) TotalSupply() (decimal.Decimal, decimal.Decimal, error) {
	type TokenDetail struct {
		TotalIssuance       decimal.Decimal `json:"total_issuance"`
		TokenDecimals       int             `json:"token_decimals"`
		BondedLockedBalance decimal.Decimal `json:"bonded_locked_balance"`
	}
	type SubscanTokenRes struct {
		Data struct {
			Detail map[string]TokenDetail `json:"detail"`
		} `json:"data"`
	}
	var res SubscanTokenRes
	b, _ := util.HttpGet(fmt.Sprintf("%s/api/scan/token", config.Cfg.SubscanHost))
	err := util.UnmarshalAny(&res, b)
	if err != nil{
		return decimal.Decimal{}, decimal.Decimal{}, err
	}
	detail := res.Data.Detail[strings.ToUpper(c.Code)]
	return detail.TotalIssuance.Div(decimal.New(1, int32(detail.TokenDecimals))),
		detail.BondedLockedBalance.Div(decimal.New(1, int32(detail.TokenDecimals))), nil
}

func (s *SupplyDetail) filterBalance(filterAddress map[string][]string) decimal.Decimal {
	filter := filterAddress[s.Network]
	wg := sync.WaitGroup{}
	var sum decimal.Decimal
	for _, address := range filter {
		go func(address string) {
			defer wg.Done()
			switch s.Network {
			case "Tron":
				sum = sum.Add(decimal.NewFromBigInt(parallel.RingTronBalance(s.Contract, util.TrxBase58toHexAddress(address)), 0))
			case "Ethereum":
				s, err := parallel.RingEthBalance(s.Contract, address)
				if err != nil{
					return
				}
				sum = sum.Add(decimal.NewFromBigInt(s, 0))
				fmt.Println(sum, address)
			}
		}(address)
		wg.Add(1)
	}
	wg.Wait()
	return sum
}
