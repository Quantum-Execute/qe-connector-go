package qe_connector

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Quantum-Execute/qe-connector-go/handlers"
	"github.com/bitly/go-simplejson"
)

// TimeInForceType define time in force type of order
type TimeInForceType string

// UserDataEventType define spot user data event type
type UserDataEventType string

// Client define API client
type Client struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
}

type doFunc func(req *http.Request) (*http.Response, error)

// Globals
const (
	timestampKey  = "timestamp"
	signatureKey  = "signature"
	recvWindowKey = "recvWindow"
)

func currentTimestamp() int64 {
	return FormatTimestamp(time.Now())
}

// FormatTimestamp formats a time into Unix timestamp in milliseconds, as requested by Binance.
func FormatTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

// NewClient Create client function for initialising new QE client
func NewClient(apiKey string, secretKey string, baseURL ...string) *Client {
	u := "https://api.quantumexecute.com"

	if len(baseURL) > 0 {
		u = baseURL[0]
	}

	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    u,
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, Name, log.LstdFlags),
	}
}

// NewTestClient Create client function for initialising new QE Test client
func NewTestClient(apiKey string, secretKey string, baseURL ...string) *Client {
	u := "https://testapi.quantumexecute.com"

	if len(baseURL) > 0 {
		u = baseURL[0]
	}

	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    u,
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, Name, log.LstdFlags),
		Debug:      true,
	}
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)
	if r.recvWindow > 0 {
		r.setParam(recvWindowKey, r.recvWindow)
	}
	if r.secType == secTypeSigned {
		r.setParam(timestampKey, currentTimestamp()-c.TimeOffset)
	}
	queryString := r.query.Encode()
	body := &bytes.Buffer{}
	bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	header.Set("User-Agent", fmt.Sprintf("%s/%s", Name, Version))
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}
	if r.secType == secTypeAPIKey || r.secType == secTypeSigned {
		header.Set("X-MBX-APIKEY", c.APIKey)
	}

	if r.secType == secTypeSigned {
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return err
		}
		v := url.Values{}
		v.Set(signatureKey, fmt.Sprintf("%x", (mac.Sum(nil))))
		if queryString == "" {
			queryString = v.Encode()
		} else {
			queryString = fmt.Sprintf("%s&%s", queryString, v.Encode())
		}
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	c.debug("full url: %s, body: %s", fullURL, bodyString)
	r.fullURL = fullURL
	r.header = header
	r.body = body
	return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.debug("request: %#v", req)
	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	c.debug("response: %#v", res)
	c.debug("response body: %s", string(data))
	c.debug("response status code: %d", res.StatusCode)

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(handlers.APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			c.debug("failed to unmarshal json: %s", e)
		}
		return nil, apiErr
	}
	var respData *handlers.APISuccess
	err = json.Unmarshal(data, &respData)
	if err != nil {
		c.debug("failed to unmarshal json: %s", err)
		return nil, err
	}
	if respData.Code != 200 {
		c.debug("response status code: %d", respData.Code)
		return nil, errors.New(respData.Reason)
	}
	return json.Marshal(respData.Message)
}

func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (c *Client) NewListExchangeApisService() *ListExchangeApisService {
	return &ListExchangeApisService{c: c}
}
func (c *Client) NewGetMasterOrdersService() *GetMasterOrdersService {
	return &GetMasterOrdersService{c: c}
}
func (c *Client) NewGetOrderFillsService() *GetOrderFillsService {
	return &GetOrderFillsService{c: c}
}
func (c *Client) NewGetMasterOrderDetailService() *GetMasterOrderDetailService {
	return &GetMasterOrderDetailService{c: c}
}
func (c *Client) NewCreateMasterOrderService() *CreateMasterOrderService {
	return &CreateMasterOrderService{c: c}
}
func (c *Client) NewCancelMasterOrderService() *CancelMasterOrderService {
	return &CancelMasterOrderService{c: c}
}
func (c *Client) NewPauseMasterOrderService() *PauseMasterOrderService {
	return &PauseMasterOrderService{c: c}
}
func (c *Client) NewResumeMasterOrderService() *ResumeMasterOrderService {
	return &ResumeMasterOrderService{c: c}
}
func (c *Client) NewUpdateMasterOrderParamsService() *UpdateMasterOrderParamsService {
	return &UpdateMasterOrderParamsService{c: c}
}
func (c *Client) NewCreateListenKeyService() *CreateListenKeyService {
	return &CreateListenKeyService{c: c}
}
func (c *Client) NewGetTcaAnalysisService() *GetTcaAnalysisService {
	return &GetTcaAnalysisService{c: c}
}
func (c *Client) NewTradingPairsService() *TradingPairsService {
	return &TradingPairsService{c: c}
}
func (c *Client) NewPingServer() *PingService {
	return &PingService{c: c}
}
func (c *Client) NewTimestampService() *TimestampService {
	return &TimestampService{c: c}
}

// NewWebSocketService create WebSocket service for real-time data streaming
func (c *Client) NewWebSocketService(host ...string) *WebSocketService {
	return NewWebSocketService(c, host...)
}

// NewGetAccountBalanceService create service for getting Binance spot account balance
func (c *Client) NewGetAccountBalanceService() *GetAccountBalanceService {
	return &GetAccountBalanceService{c: c}
}

// NewGetMarginBalanceService create service for getting Binance futures account balance
func (c *Client) NewGetMarginBalanceService() *GetMarginBalanceService {
	return &GetMarginBalanceService{c: c}
}

// NewGetPv1BalanceService create service for getting Binance PAPI PV1 balance
func (c *Client) NewGetPv1BalanceService() *GetPv1BalanceService {
	return &GetPv1BalanceService{c: c}
}

// NewGetOkxAccountBalanceService create service for getting OKX account balance
func (c *Client) NewGetOkxAccountBalanceService() *GetOkxAccountBalanceService {
	return &GetOkxAccountBalanceService{c: c}
}

// NewGetFapiPositionSideDialService create service for getting Binance FAPI position side dual status
func (c *Client) NewGetFapiPositionSideDialService() *GetFapiPositionSideDialService {
	return &GetFapiPositionSideDialService{c: c}
}

// NewGetPapiUmPositionSideDualService create service for getting Binance PAPI UM position side dual status
func (c *Client) NewGetPapiUmPositionSideDualService() *GetPapiUmPositionSideDualService {
	return &GetPapiUmPositionSideDualService{c: c}
}

// NewGetOkxAccountPositionsService create service for getting OKX account positions
func (c *Client) NewGetOkxAccountPositionsService() *GetOkxAccountPositionsService {
	return &GetOkxAccountPositionsService{c: c}
}

// NewGetOkxAccountMaxSizeService create service for getting OKX account max order size
func (c *Client) NewGetOkxAccountMaxSizeService() *GetOkxAccountMaxSizeService {
	return &GetOkxAccountMaxSizeService{c: c}
}

// NewGetLtpPositionService create service for getting LTP account positions
func (c *Client) NewGetLtpPositionService() *GetLtpPositionService {
	return &GetLtpPositionService{c: c}
}

// NewGetDeribitPositionService create service for getting Deribit account positions
func (c *Client) NewGetDeribitPositionService() *GetDeribitPositionService {
	return &GetDeribitPositionService{c: c}
}

// NewGetUmAccountService create service for getting Binance PAPI UM account
func (c *Client) NewGetUmAccountService() *GetUmAccountService {
	return &GetUmAccountService{c: c}
}

// NewGetCmAccountService create service for getting Binance PAPI CM account
func (c *Client) NewGetCmAccountService() *GetCmAccountService {
	return &GetCmAccountService{c: c}
}

// NewGetPv1AccountService create service for getting Binance PAPI PV1 account
func (c *Client) NewGetPv1AccountService() *GetPv1AccountService {
	return &GetPv1AccountService{c: c}
}

// NewGetDapiAccountService create service for getting Binance DAPI account
func (c *Client) NewGetDapiAccountService() *GetDapiAccountService {
	return &GetDapiAccountService{c: c}
}

// NewGetFapiAccountService create service for getting Binance FAPI account
func (c *Client) NewGetFapiAccountService() *GetFapiAccountService {
	return &GetFapiAccountService{c: c}
}

// NewGetCrossMarginAccountDetailService create service for getting Binance cross margin account detail
func (c *Client) NewGetCrossMarginAccountDetailService() *GetCrossMarginAccountDetailService {
	return &GetCrossMarginAccountDetailService{c: c}
}

// NewGetLtpAccountService create service for getting LTP account info
func (c *Client) NewGetLtpAccountService() *GetLtpAccountService {
	return &GetLtpAccountService{c: c}
}

// NewGetLtpPortfolioAssetService create service for getting LTP portfolio assets
func (c *Client) NewGetLtpPortfolioAssetService() *GetLtpPortfolioAssetService {
	return &GetLtpPortfolioAssetService{c: c}
}

// NewGetDeribitAccountService create service for getting Deribit account info
func (c *Client) NewGetDeribitAccountService() *GetDeribitAccountService {
	return &GetDeribitAccountService{c: c}
}
