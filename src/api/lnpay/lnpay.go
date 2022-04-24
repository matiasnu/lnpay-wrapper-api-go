package lnpay

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/imroc/req"
)

const (
	BASE_URL = "https://api.lnpay.co/v1"
	TEST_KEY = "pak_O0iUMxk8kK_qUzkT4YKFvp1ZsUtp"
)

type Client struct {
	header req.Header
}

var transactionLock sync.Mutex

// NewClient is the first function you must call. Pass your main API key here.
// It will return a client you can later use to access wallets and transactions.
// You can find it at https://lnpay.co/developers/dashboard
func NewClient(key string) *Client {
	return &Client{
		header: req.Header{
			"X-Api-Key":    key,
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}
}

// Transaction
func (c *Client) Transaction(lntxId string) (lnTx LnTx, err error) {
	resp, err := req.Get(BASE_URL+"/lntx/"+lntxId, c.header)
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&lnTx)
	return
}
func (c *Client) QueryRoutes(pubKey, amt string) (routes QueryRoutes, err error) {
	resp, err := req.Get(fmt.Sprintf("%s/node/default/payments/queryroutes?pub_key=%s&amt=%s", BASE_URL, pubKey, amt), c.header)
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&routes)
	return
}
func (c *Client) DecodeInvoice(paymentRequest string) (invoice Invoice, err error) {
	resp, err := req.Get(fmt.Sprintf("%s/node/default/payments/decodeinvoice?payment_request=%s", BASE_URL, paymentRequest), c.header)
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&invoice)
	return
}

// CreateWallet creates a new wallet with a given descriptive label.
// It will return the wallet object which you can use to create invoices and payments.
// https://docs.lnpay.co/wallet/create-wallet
func (c *Client) CreateWallet(label string) (wal Wallet, err error) {
	resp, err := req.Post(BASE_URL+"/wallet", c.header, req.BodyJSON(struct {
		UserLabel string `json:"user_label"`
	}{label}))
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&wal)
	wal.Client = c
	wal.BaseUrl = BASE_URL + "/wallet/" + wal.AccessKeys.WalletAdmin[0]
	return
}

// Wallet returns a wallet that was already created.
// Pass the wallet key you got when creating it.
// I can be either the admin, invoice or read-only key.
func (c *Client) Wallet(key string) *Wallet {
	return &Wallet{
		Client:     c,
		StatusType: StatusType{},
		AccessKeys: AccessKeys{},
		BaseUrl:    BASE_URL + "/wallet/" + key,
	}
}

type WalletStatusType map[string]string

// Details returns basic information about a wallet, such as its id, label or balance.
// https://docs.lnpay.co/wallet/get-balance
func (w *Wallet) Details() (wal Wallet, err error) {
	resp, err := req.Get(w.BaseUrl, w.header)
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&wal)
	return
}
func (w *Wallet) UpdateBalance() error {
	wal, err := w.Details()
	if err != nil {
		return err
	}
	if wal.Balance != w.Balance {
		w.Balance = wal.Balance
	}
	return nil
}

// Transactions returns a list of the transactions associated with the wallet.
// https://docs.lnpay.co/wallet/get-transactions
func (w *Wallet) Transactions(page int) (txs []Wtx, header http.Header, err error) {
	resp, err := req.Get(w.BaseUrl+fmt.Sprintf("/transactions?per-page=10&page=%d", page), w.header)
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}
	header = resp.Response().Header
	err = resp.ToJSON(&txs)
	return
}

type InvoiceParams struct {
	Memo        string `json:"memo"`         // the invoice description.
	NumSatoshis int64  `json:"num_satoshis"` // the invoice amount.
	Expiry      int64  `json:"expiry"`       // seconds, default 86400 (1 day).

	// custom data you may want to associate with this invoice. optional.
	PassThru map[string]interface{} `json:"passThru"`

	// base64-encoded. If this is provided, memo will be ignored.
	// don't use this if you don't know what it means.
	DescriptionHash string `json:"description_hash"`
}

// Invoice creates an invoice associated with this wallet.
// https://docs.lnpay.co/wallet/generate-invoice
func (w *Wallet) Invoice(params InvoiceParams) (lntx LnTx, err error) {
	transactionLock.Lock()
	resp, err := req.Post(w.BaseUrl+"/invoice", w.header, req.BodyJSON(&params))
	transactionLock.Unlock()
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&lntx)
	return
}

type PayParams struct {
	// the BOLT11 payment request you want to pay.
	PaymentRequest string `json:"payment_request"`

	// custom data you may want to associate with this invoice. optional.
	PassThru map[string]interface{} `json:"passThru"`
}

// Pay pays a given invoice with funds from the wallet.
// https://docs.lnpay.co/wallet/pay-invoice
func (w *Wallet) Pay(params PayParams) (wtx Wtx, err error) {
	transactionLock.Lock()
	resp, err := req.Post(w.BaseUrl+"/withdraw", w.header, req.BodyJSON(&params))
	transactionLock.Unlock()
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&wtx)
	return
}

type TransferParams struct {
	Memo         string `json:"memo"`           // the transfer description.
	NumSatoshis  int64  `json:"num_satoshis"`   // the transfer amount.
	DestWalletId string `json:"dest_wallet_id"` // the key or id of the destination
}

// Transfer transfers between two lnpay.co wallets.
// https://docs.lnpay.co/wallet/transfers-between-wallets
func (w *Wallet) Transfer(params TransferParams) (wtx Wtx, err error) {
	transactionLock.Lock()
	resp, err := req.Post(w.BaseUrl+"/transfer", w.header, req.BodyJSON(&params))
	transactionLock.Unlock()
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&wtx)
	return
}
