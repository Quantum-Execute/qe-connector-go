package qe_connector

import (
	"context"
	"encoding/json"
	"net/http"
)

// ─────────────────────────────────────────────────────────────────────────────
// 余额类
// ─────────────────────────────────────────────────────────────────────────────

// GetAccountBalanceService get Binance spot account balance
type GetAccountBalanceService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetAccountBalanceService) BindingId(v string) *GetAccountBalanceService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetAccountBalanceService) Do(ctx context.Context, opts ...RequestOption) (res *AccountBalanceReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/account-balance",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(AccountBalanceReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AccountBalanceReply Binance spot account balance response
type AccountBalanceReply struct {
	Balances    []SpotBalanceItem `json:"balances"`
	Exchange    string            `json:"exchange"`
	AccountType string            `json:"accountType"`
	UpdateTime  string            `json:"updateTime"`
}

// SpotBalanceItem single asset balance in spot account
type SpotBalanceItem struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetMarginBalanceService get Binance futures account balance
type GetMarginBalanceService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetMarginBalanceService) BindingId(v string) *GetMarginBalanceService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetMarginBalanceService) Do(ctx context.Context, opts ...RequestOption) (res *MarginBalanceReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/margin-balance",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(MarginBalanceReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// MarginBalanceReply Binance futures account balance response
type MarginBalanceReply struct {
	Balances    []FuturesBalanceItem `json:"balances"`
	Exchange    string               `json:"exchange"`
	AccountType string               `json:"accountType"`
	UpdateTime  string               `json:"updateTime"`
}

// FuturesBalanceItem single asset balance in futures account
type FuturesBalanceItem struct {
	Asset              string `json:"asset"`
	WalletBalance      string `json:"walletBalance"`
	AvailableBalance   string `json:"availableBalance"`
	CrossWalletBalance string `json:"crossWalletBalance"`
	CrossUnPnl         string `json:"crossUnPnl"`
	MarginBalance      string `json:"marginBalance"`
	MaxWithdrawAmount  string `json:"maxWithdrawAmount"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetPv1BalanceService get Binance PAPI PV1 balance
type GetPv1BalanceService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetPv1BalanceService) BindingId(v string) *GetPv1BalanceService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetPv1BalanceService) Do(ctx context.Context, opts ...RequestOption) (res *Pv1BalanceReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/pv1-balance",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Pv1BalanceReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Pv1BalanceReply Binance PAPI PV1 balance response
type Pv1BalanceReply struct {
	Exchange    string         `json:"exchange"`
	AccountType string         `json:"accountType"`
	Balances    []Pv1BalanceItem `json:"balances"`
}

// Pv1BalanceItem single asset balance in PAPI PV1 account
type Pv1BalanceItem struct {
	Asset               string `json:"asset"`
	TotalWalletBalance  string `json:"totalWalletBalance"`
	CrossMarginBorrowed string `json:"crossMarginBorrowed"`
	CrossMarginFree     string `json:"crossMarginFree"`
	CrossMarginInterest string `json:"crossMarginInterest"`
	CrossMarginLocked   string `json:"crossMarginLocked"`
	UmWalletBalance     string `json:"umWalletBalance"`
	UmUnrealizedPnl     string `json:"umUnrealizedPnl"`
	CmWalletBalance     string `json:"cmWalletBalance"`
	CmUnrealizedPnl     string `json:"cmUnrealizedPnl"`
	UpdateTime          int64  `json:"updateTime"`
	NegativeBalance     string `json:"negativeBalance"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetOkxAccountBalanceService get OKX account balance
type GetOkxAccountBalanceService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetOkxAccountBalanceService) BindingId(v string) *GetOkxAccountBalanceService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetOkxAccountBalanceService) Do(ctx context.Context, opts ...RequestOption) (res *OkxAccountBalanceReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/okx-account-balance",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(OkxAccountBalanceReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// OkxAccountBalanceReply OKX account balance response
type OkxAccountBalanceReply struct {
	Data     []OkxBalanceData `json:"data"`
	Exchange string           `json:"exchange"`
}

// OkxBalanceData OKX account balance data
type OkxBalanceData struct {
	TotalEq     string           `json:"totalEq"`
	AvailEq     string           `json:"availEq"`
	AdjEq       string           `json:"adjEq"`
	Imr         string           `json:"imr"`
	Mmr         string           `json:"mmr"`
	MgnRatio    string           `json:"mgnRatio"`
	NotionalUsd string           `json:"notionalUsd"`
	OrdFroz     string           `json:"ordFroz"`
	Upl         string           `json:"upl"`
	UTime       string           `json:"uTime"`
	Details     []OkxBalanceDetail `json:"details"`
}

// OkxBalanceDetail single currency balance detail in OKX account
type OkxBalanceDetail struct {
	Ccy       string `json:"ccy"`
	Eq        string `json:"eq"`
	EqUsd     string `json:"eqUsd"`
	AvailBal  string `json:"availBal"`
	AvailEq   string `json:"availEq"`
	CashBal   string `json:"cashBal"`
	FrozenBal string `json:"frozenBal"`
	Upl       string `json:"upl"`
	Liab      string `json:"liab"`
	Interest  string `json:"interest"`
}

// ─────────────────────────────────────────────────────────────────────────────
// 持仓类
// ─────────────────────────────────────────────────────────────────────────────

// GetFapiPositionSideDialService get Binance FAPI position side dual status
type GetFapiPositionSideDialService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetFapiPositionSideDialService) BindingId(v string) *GetFapiPositionSideDialService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetFapiPositionSideDialService) Do(ctx context.Context, opts ...RequestOption) (res *PositionSideDualReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/fapi-position-side-dial",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(PositionSideDualReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ─────────────────────────────────────────────────────────────────────────────

// GetPapiUmPositionSideDualService get Binance PAPI UM position side dual status
type GetPapiUmPositionSideDualService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetPapiUmPositionSideDualService) BindingId(v string) *GetPapiUmPositionSideDualService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetPapiUmPositionSideDualService) Do(ctx context.Context, opts ...RequestOption) (res *PositionSideDualReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/papi-um-position-side-dual",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(PositionSideDualReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PositionSideDualReply position side dual status response (shared by FAPI and PAPI UM)
type PositionSideDualReply struct {
	DualSidePosition bool   `json:"dualSidePosition"`
	Exchange         string `json:"exchange"`
	AccountType      string `json:"accountType"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetOkxAccountPositionsService get OKX account positions
type GetOkxAccountPositionsService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetOkxAccountPositionsService) BindingId(v string) *GetOkxAccountPositionsService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetOkxAccountPositionsService) Do(ctx context.Context, opts ...RequestOption) (res *OkxAccountPositionsReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/okx-account-positions",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(OkxAccountPositionsReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// OkxAccountPositionsReply OKX account positions response
type OkxAccountPositionsReply struct {
	Data     []OkxPositionItem `json:"data"`
	Exchange string            `json:"exchange"`
}

// OkxPositionItem single position in OKX account
type OkxPositionItem struct {
	InstId      string `json:"instId"`
	InstType    string `json:"instType"`
	Pos         string `json:"pos"`
	PosSide     string `json:"posSide"`
	AvgPx       string `json:"avgPx"`
	MarkPx      string `json:"markPx"`
	LiqPx       string `json:"liqPx"`
	Upl         string `json:"upl"`
	UplRatio    string `json:"uplRatio"`
	Lever       string `json:"lever"`
	MgnMode     string `json:"mgnMode"`
	Imr         string `json:"imr"`
	Mmr         string `json:"mmr"`
	Margin      string `json:"margin"`
	NotionalUsd string `json:"notionalUsd"`
	Adl         string `json:"adl"`
	Ccy         string `json:"ccy"`
	Pnl         string `json:"pnl"`
	RealizedPnl string `json:"realizedPnl"`
	Fee         string `json:"fee"`
	FundingFee  string `json:"fundingFee"`
	CTime       string `json:"cTime"`
	UTime       string `json:"uTime"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetOkxAccountMaxSizeService get OKX account max order size
type GetOkxAccountMaxSizeService struct {
	c         *Client
	bindingId string
	instId    string
	tdMode    string
}

// BindingId set bindingId
func (s *GetOkxAccountMaxSizeService) BindingId(v string) *GetOkxAccountMaxSizeService {
	s.bindingId = v
	return s
}

// InstId set instId (e.g. BTC-USDT-SWAP)
func (s *GetOkxAccountMaxSizeService) InstId(v string) *GetOkxAccountMaxSizeService {
	s.instId = v
	return s
}

// TdMode set tdMode (cross / isolated / cash)
func (s *GetOkxAccountMaxSizeService) TdMode(v string) *GetOkxAccountMaxSizeService {
	s.tdMode = v
	return s
}

// Do send request
func (s *GetOkxAccountMaxSizeService) Do(ctx context.Context, opts ...RequestOption) (res *OkxAccountMaxSizeReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/okx-account-max-size",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	r.setParam("instId", s.instId)
	r.setParam("tdMode", s.tdMode)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(OkxAccountMaxSizeReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// OkxAccountMaxSizeReply OKX account max order size response
type OkxAccountMaxSizeReply struct {
	Data     []OkxMaxSizeItem `json:"data"`
	Exchange string           `json:"exchange"`
}

// OkxMaxSizeItem max buy/sell size for a product in OKX account
type OkxMaxSizeItem struct {
	Ccy     string `json:"ccy"`
	InstId  string `json:"instId"`
	MaxBuy  string `json:"maxBuy"`
	MaxSell string `json:"maxSell"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetLtpPositionService get LTP account positions
type GetLtpPositionService struct {
	c         *Client
	bindingId string
	sym       *string
}

// BindingId set bindingId
func (s *GetLtpPositionService) BindingId(v string) *GetLtpPositionService {
	s.bindingId = v
	return s
}

// Sym set sym (optional, filter by trading pair)
func (s *GetLtpPositionService) Sym(v string) *GetLtpPositionService {
	s.sym = &v
	return s
}

// Do send request
func (s *GetLtpPositionService) Do(ctx context.Context, opts ...RequestOption) (res *LtpPositionReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/ltp-position",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	if s.sym != nil {
		r.setParam("sym", *s.sym)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(LtpPositionReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// LtpPositionReply LTP account positions response
type LtpPositionReply struct {
	Data     []LtpPositionItem `json:"data"`
	Exchange string            `json:"exchange"`
}

// LtpPositionItem single position in LTP account
type LtpPositionItem struct {
	PositionId        string `json:"positionId"`
	PortfolioId       string `json:"portfolioId"`
	Sym               string `json:"sym"`
	PositionSide      string `json:"positionSide"`
	PositionQty       string `json:"positionQty"`
	PositionValue     string `json:"positionValue"`
	PositionMargin    string `json:"positionMargin"`
	PositionMm        string `json:"positionMm"`
	UnrealizedPnl     string `json:"unrealizedPnl"`
	UnrealizedPnlRate string `json:"unrealizedPnlRate"`
	AvgPrice          string `json:"avgPrice"`
	MarkPrice         string `json:"markPrice"`
	LiqPrice          string `json:"liqPrice"`
	Leverage          string `json:"leverage"`
	MaxLeverage       string `json:"maxLeverage"`
	RiskLevel         string `json:"riskLevel"`
	Fee               string `json:"fee"`
	FundingFee        string `json:"fundingFee"`
	TpslOrder         string `json:"tpslOrder"`
	CreateAt          string `json:"createAt"`
	UpdateAt          string `json:"updateAt"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetDeribitPositionService get Deribit account positions
type GetDeribitPositionService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetDeribitPositionService) BindingId(v string) *GetDeribitPositionService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetDeribitPositionService) Do(ctx context.Context, opts ...RequestOption) (res *DeribitPositionReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/deribit-position",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(DeribitPositionReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeribitPositionReply Deribit account positions response
type DeribitPositionReply struct {
	Data     []DeribitPositionItem `json:"data"`
	Exchange string                `json:"exchange"`
}

// DeribitPositionItem single position in Deribit account
type DeribitPositionItem struct {
	InstrumentName             string  `json:"instrumentName"`
	Direction                  string  `json:"direction"`
	Size                       float64 `json:"size"`
	AveragePrice               float64 `json:"averagePrice"`
	MarkPrice                  float64 `json:"markPrice"`
	IndexPrice                 float64 `json:"indexPrice"`
	FloatingProfitLoss         float64 `json:"floatingProfitLoss"`
	TotalProfitLoss            float64 `json:"totalProfitLoss"`
	InitialMargin              float64 `json:"initialMargin"`
	MaintenanceMargin          float64 `json:"maintenanceMargin"`
	EstimatedLiquidationPrice  float64 `json:"estimatedLiquidationPrice"`
	Leverage                   int32   `json:"leverage"`
	Kind                       string  `json:"kind"`
	SizeCurrency               float64 `json:"sizeCurrency"`
	Delta                      float64 `json:"delta"`
	RealizedFunding            float64 `json:"realizedFunding"`
	RealizedProfitLoss         float64 `json:"realizedProfitLoss"`
	SettlementPrice            float64 `json:"settlementPrice"`
}

// ─────────────────────────────────────────────────────────────────────────────
// 账户类
// ─────────────────────────────────────────────────────────────────────────────

// GetUmAccountService get Binance PAPI UM account
type GetUmAccountService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetUmAccountService) BindingId(v string) *GetUmAccountService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetUmAccountService) Do(ctx context.Context, opts ...RequestOption) (res *UmAccountReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/um-account",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(UmAccountReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UmAccountReply Binance PAPI UM account response
type UmAccountReply struct {
	TradeGroupId int32                `json:"tradeGroupId"`
	Assets       []PapiAccountAsset   `json:"assets"`
	Positions    []PapiAccountPosition `json:"positions"`
	Exchange     string               `json:"exchange"`
	AccountType  string               `json:"accountType"`
	UpdateTime   string               `json:"updateTime"`
}

// PapiAccountAsset asset info in PAPI UM/CM account
type PapiAccountAsset struct {
	Asset                    string `json:"asset"`
	CrossWalletBalance       string `json:"crossWalletBalance"`
	CrossUnPnl               string `json:"crossUnPnl"`
	MaintMargin              string `json:"maintMargin"`
	InitialMargin            string `json:"initialMargin"`
	PositionInitialMargin    string `json:"positionInitialMargin"`
	OpenOrderInitialMargin   string `json:"openOrderInitialMargin"`
	UpdateTime               int64  `json:"updateTime"`
}

// PapiAccountPosition position info in PAPI UM/CM account
type PapiAccountPosition struct {
	Symbol          string `json:"symbol"`
	PositionAmt     string `json:"positionAmt"`
	PositionSide    string `json:"positionSide"`
	EntryPrice      string `json:"entryPrice"`
	BreakEvenPrice  string `json:"breakEvenPrice"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	Leverage        string `json:"leverage"`
	InitialMargin   string `json:"initialMargin"`
	MaintMargin     string `json:"maintMargin"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetCmAccountService get Binance PAPI CM account
type GetCmAccountService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetCmAccountService) BindingId(v string) *GetCmAccountService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetCmAccountService) Do(ctx context.Context, opts ...RequestOption) (res *CmAccountReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/cm-account",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CmAccountReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CmAccountReply Binance PAPI CM account response
type CmAccountReply struct {
	Assets      []PapiAccountAsset    `json:"assets"`
	Positions   []PapiAccountPosition `json:"positions"`
	Exchange    string                `json:"exchange"`
	AccountType string                `json:"accountType"`
	UpdateTime  string                `json:"updateTime"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetPv1AccountService get Binance PAPI PV1 account
type GetPv1AccountService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetPv1AccountService) BindingId(v string) *GetPv1AccountService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetPv1AccountService) Do(ctx context.Context, opts ...RequestOption) (res *Pv1AccountReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/pv1-account",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Pv1AccountReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Pv1AccountReply Binance PAPI PV1 account response
type Pv1AccountReply struct {
	Exchange                    string `json:"exchange"`
	AccountType                 string `json:"accountType"`
	UniMmr                      string `json:"uniMmr"`
	AccountEquity               string `json:"accountEquity"`
	ActualEquity                string `json:"actualEquity"`
	AccountInitialMargin        string `json:"accountInitialMargin"`
	AccountMaintMargin          string `json:"accountMaintMargin"`
	AccountStatus               string `json:"accountStatus"`
	VirtualMaxWithdrawAmount    string `json:"virtualMaxWithdrawAmount"`
	TotalAvailableBalance       string `json:"totalAvailableBalance"`
	TotalMarginOpenLoss         string `json:"totalMarginOpenLoss"`
	UpdateTime                  string `json:"updateTime"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetDapiAccountService get Binance DAPI account
type GetDapiAccountService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetDapiAccountService) BindingId(v string) *GetDapiAccountService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetDapiAccountService) Do(ctx context.Context, opts ...RequestOption) (res *DapiAccountReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/dapi-account",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(DapiAccountReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DapiAccountReply Binance DAPI account response
type DapiAccountReply struct {
	Assets      []DapiAsset    `json:"assets"`
	Positions   []DapiPosition `json:"positions"`
	CanDeposit  bool           `json:"canDeposit"`
	CanTrade    bool           `json:"canTrade"`
	CanWithdraw bool           `json:"canWithdraw"`
	FeeTier     int32          `json:"feeTier"`
	Exchange    string         `json:"exchange"`
	AccountType string         `json:"accountType"`
	UpdateTime  string         `json:"updateTime"`
}

// DapiAsset asset info in DAPI account
type DapiAsset struct {
	Asset              string `json:"asset"`
	WalletBalance      string `json:"walletBalance"`
	UnrealizedProfit   string `json:"unrealizedProfit"`
	MarginBalance      string `json:"marginBalance"`
	AvailableBalance   string `json:"availableBalance"`
	MaxWithdrawAmount  string `json:"maxWithdrawAmount"`
	CrossWalletBalance string `json:"crossWalletBalance"`
	CrossUnPnl         string `json:"crossUnPnl"`
}

// DapiPosition position info in DAPI account
type DapiPosition struct {
	Symbol          string `json:"symbol"`
	PositionAmt     string `json:"positionAmt"`
	PositionSide    string `json:"positionSide"`
	EntryPrice      string `json:"entryPrice"`
	BreakEvenPrice  string `json:"breakEvenPrice"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	Leverage        string `json:"leverage"`
	Isolated        bool   `json:"isolated"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetFapiAccountService get Binance FAPI account
type GetFapiAccountService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetFapiAccountService) BindingId(v string) *GetFapiAccountService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetFapiAccountService) Do(ctx context.Context, opts ...RequestOption) (res *FapiAccountReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/fapi-account",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(FapiAccountReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FapiAccountReply Binance FAPI account response
type FapiAccountReply struct {
	TotalWalletBalance       string         `json:"totalWalletBalance"`
	TotalUnrealizedProfit    string         `json:"totalUnrealizedProfit"`
	TotalMarginBalance       string         `json:"totalMarginBalance"`
	AvailableBalance         string         `json:"availableBalance"`
	MaxWithdrawAmount        string         `json:"maxWithdrawAmount"`
	TotalInitialMargin       string         `json:"totalInitialMargin"`
	TotalMaintMargin         string         `json:"totalMaintMargin"`
	TotalCrossWalletBalance  string         `json:"totalCrossWalletBalance"`
	TotalCrossUnPnl          string         `json:"totalCrossUnPnl"`
	Assets                   []FapiAsset    `json:"assets"`
	Positions                []FapiPosition `json:"positions"`
	Exchange                 string         `json:"exchange"`
	AccountType              string         `json:"accountType"`
}

// FapiAsset asset info in FAPI account
type FapiAsset struct {
	Asset             string `json:"asset"`
	WalletBalance     string `json:"walletBalance"`
	UnrealizedProfit  string `json:"unrealizedProfit"`
	MarginBalance     string `json:"marginBalance"`
	AvailableBalance  string `json:"availableBalance"`
	MaxWithdrawAmount string `json:"maxWithdrawAmount"`
	MarginAvailable   bool   `json:"marginAvailable"`
}

// FapiPosition position info in FAPI account
type FapiPosition struct {
	Symbol           string `json:"symbol"`
	PositionAmt      string `json:"positionAmt"`
	PositionSide     string `json:"positionSide"`
	EntryPrice       string `json:"entryPrice"`
	BreakEvenPrice   string `json:"breakEvenPrice"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	Leverage         string `json:"leverage"`
	Isolated         bool   `json:"isolated"`
	Notional         string `json:"notional"`
	IsolatedWallet   string `json:"isolatedWallet"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetCrossMarginAccountDetailService get Binance cross margin account detail
type GetCrossMarginAccountDetailService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetCrossMarginAccountDetailService) BindingId(v string) *GetCrossMarginAccountDetailService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetCrossMarginAccountDetailService) Do(ctx context.Context, opts ...RequestOption) (res *CrossMarginAccountDetailReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/cross-margin-account-detail",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CrossMarginAccountDetailReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CrossMarginAccountDetailReply Binance cross margin account detail response
type CrossMarginAccountDetailReply struct {
	Exchange            string                 `json:"exchange"`
	AccountType         string                 `json:"accountType"`
	BorrowEnabled       bool                   `json:"borrowEnabled"`
	TradeEnabled        bool                   `json:"tradeEnabled"`
	TransferEnabled     bool                   `json:"transferEnabled"`
	MarginLevel         string                 `json:"marginLevel"`
	TotalAssetOfBtc     string                 `json:"totalAssetOfBtc"`
	TotalLiabilityOfBtc string                 `json:"totalLiabilityOfBtc"`
	TotalNetAssetOfBtc  string                 `json:"totalNetAssetOfBtc"`
	UserAssets          []CrossMarginUserAsset `json:"userAssets"`
}

// CrossMarginUserAsset single asset in cross margin account
type CrossMarginUserAsset struct {
	Asset    string `json:"asset"`
	Free     string `json:"free"`
	Locked   string `json:"locked"`
	Borrowed string `json:"borrowed"`
	Interest string `json:"interest"`
	NetAsset string `json:"netAsset"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetLtpAccountService get LTP account info
type GetLtpAccountService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetLtpAccountService) BindingId(v string) *GetLtpAccountService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetLtpAccountService) Do(ctx context.Context, opts ...RequestOption) (res *LtpAccountReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/ltp-account",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(LtpAccountReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// LtpAccountReply LTP account info response
type LtpAccountReply struct {
	Data     []LtpAccountItem `json:"data"`
	Exchange string           `json:"exchange"`
}

// LtpAccountItem single portfolio account info in LTP
type LtpAccountItem struct {
	PortfolioId     string `json:"portfolioId"`
	ExchangeType    string `json:"exchangeType"`
	Equity          string `json:"equity"`
	MaintainMargin  string `json:"maintainMargin"`
	PositionValue   string `json:"positionValue"`
	UniMmr          string `json:"uniMmr"`
	RiskRatio       string `json:"riskRatio"`
	AccountStatus   string `json:"accountStatus"`
	AvailableMargin string `json:"availableMargin"`
	ValidMargin     string `json:"validMargin"`
	FrozenMargin    string `json:"frozenMargin"`
	Upnl            string `json:"upnl"`
	PositionMode    string `json:"positionMode"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetLtpPortfolioAssetService get LTP portfolio assets
type GetLtpPortfolioAssetService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetLtpPortfolioAssetService) BindingId(v string) *GetLtpPortfolioAssetService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetLtpPortfolioAssetService) Do(ctx context.Context, opts ...RequestOption) (res *LtpPortfolioAssetReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/ltp-portfolio-asset",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(LtpPortfolioAssetReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// LtpPortfolioAssetReply LTP portfolio assets response
type LtpPortfolioAssetReply struct {
	Data     []LtpPortfolioAssetItem `json:"data"`
	Exchange string                  `json:"exchange"`
}

// LtpPortfolioAssetItem single portfolio asset in LTP account
type LtpPortfolioAssetItem struct {
	PortfolioId     string `json:"portfolioId"`
	Coin            string `json:"coin"`
	ExchangeType    string `json:"exchangeType"`
	Available       string `json:"available"`
	Frozen          string `json:"frozen"`
	Equity          string `json:"equity"`
	Balance         string `json:"balance"`
	Borrow          string `json:"borrow"`
	Debt            string `json:"debt"`
	MarginValue     string `json:"marginValue"`
	IndexPrice      string `json:"indexPrice"`
	MaxTransferable string `json:"maxTransferable"`
	Upnl            string `json:"upnl"`
	EquityValue     string `json:"equityValue"`
	CreateAt        string `json:"createAt"`
	UpdateAt        string `json:"updateAt"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetDeribitAccountService get Deribit account info
type GetDeribitAccountService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetDeribitAccountService) BindingId(v string) *GetDeribitAccountService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetDeribitAccountService) Do(ctx context.Context, opts ...RequestOption) (res *DeribitAccountReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/deribit-account",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(DeribitAccountReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeribitAccountReply Deribit account info response
type DeribitAccountReply struct {
	Data     []DeribitAccountItem `json:"data"`
	Exchange string               `json:"exchange"`
}

// DeribitAccountItem account info per currency in Deribit
type DeribitAccountItem struct {
	Currency                    string  `json:"currency"`
	Equity                      float64 `json:"equity"`
	Balance                     float64 `json:"balance"`
	AvailableFunds              float64 `json:"availableFunds"`
	AvailableWithdrawalFunds    float64 `json:"availableWithdrawalFunds"`
	MarginBalance               float64 `json:"marginBalance"`
	InitialMargin               float64 `json:"initialMargin"`
	MaintenanceMargin           float64 `json:"maintenanceMargin"`
	LockedBalance               float64 `json:"lockedBalance"`
	TotalPl                     float64 `json:"totalPl"`
	SessionUpl                  float64 `json:"sessionUpl"`
	SessionRpl                  float64 `json:"sessionRpl"`
	FuturesPl                   float64 `json:"futuresPl"`
	OptionsValue                float64 `json:"optionsValue"`
	OptionsDelta                float64 `json:"optionsDelta"`
	OptionsGamma                float64 `json:"optionsGamma"`
	OptionsVega                 float64 `json:"optionsVega"`
	OptionsTheta                float64 `json:"optionsTheta"`
	DeltaTotal                  float64 `json:"deltaTotal"`
	MarginModel                 string  `json:"marginModel"`
	PortfolioMarginingEnabled   bool    `json:"portfolioMarginingEnabled"`
	CrossCollateralEnabled      bool    `json:"crossCollateralEnabled"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Hyperliquid
// ─────────────────────────────────────────────────────────────────────────────

// GetHyperliquidSpotBalanceService get Hyperliquid spot account balance
type GetHyperliquidSpotBalanceService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetHyperliquidSpotBalanceService) BindingId(v string) *GetHyperliquidSpotBalanceService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetHyperliquidSpotBalanceService) Do(ctx context.Context, opts ...RequestOption) (res *HyperliquidSpotBalanceReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/hyperliquid-spot-balance",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(HyperliquidSpotBalanceReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// HyperliquidSpotBalanceReply Hyperliquid spot balance response
type HyperliquidSpotBalanceReply struct {
	Balances        []HyperliquidSpotBalanceItem `json:"balances"`
	Exchange        string                       `json:"exchange"`
	AvailableMargin string                       `json:"availableMargin"`
}

// HyperliquidSpotBalanceItem single asset balance in Hyperliquid spot account
type HyperliquidSpotBalanceItem struct {
	Coin       string `json:"coin"`
	Total      string `json:"total"`
	Hold       string `json:"hold"`
	Available  string `json:"available"`
	TotalValue string `json:"totalValue"`
	Price      string `json:"price"`
}

// ─────────────────────────────────────────────────────────────────────────────

// GetHyperliquidPositionsService get Hyperliquid perpetual positions
type GetHyperliquidPositionsService struct {
	c         *Client
	bindingId string
}

// BindingId set bindingId
func (s *GetHyperliquidPositionsService) BindingId(v string) *GetHyperliquidPositionsService {
	s.bindingId = v
	return s
}

// Do send request
func (s *GetHyperliquidPositionsService) Do(ctx context.Context, opts ...RequestOption) (res *HyperliquidPositionsReply, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/user/exchange-apis/hyperliquid-positions",
		secType:  secTypeSigned,
	}
	r.setParam("bindingId", s.bindingId)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(HyperliquidPositionsReply)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// HyperliquidPositionsReply Hyperliquid perpetual positions response
type HyperliquidPositionsReply struct {
	Positions       []HyperliquidPositionItem `json:"positions"`
	Exchange        string                    `json:"exchange"`
	Withdrawable    string                    `json:"withdrawable"`
	AccountValue    string                    `json:"accountValue"`
	TotalMarginUsed string                    `json:"totalMarginUsed"`
	TotalNtlPos     string                    `json:"totalNtlPos"`
}

// HyperliquidPositionItem single position in Hyperliquid account
type HyperliquidPositionItem struct {
	Coin           string `json:"coin"`
	Szi            string `json:"szi"`
	PositionValue  string `json:"positionValue"`
	EntryPx        string `json:"entryPx"`
	UnrealizedPnl  string `json:"unrealizedPnl"`
	LeverageType   string `json:"leverageType"`
	LeverageValue  int32  `json:"leverageValue"`
	LiquidationPx  string `json:"liquidationPx"`
	MarginUsed     string `json:"marginUsed"`
	ReturnOnEquity string `json:"returnOnEquity"`
}
