package walletrpc

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/rpc/v2/json2"
)

// Client is a monero-wallet-rpc client.
type Client interface {
	// Getbalance - Return the wallet's balance.
	Getbalance() (balance, unlockedBalance uint64, err error)
	// Getaddress - Return the wallet's address.
	// address - string; The 95-character hex address string of the monero-wallet-rpc in session.
	Getaddress() (address string, err error)
	// Getheight - Returns the wallet's current block height.
	// height - unsigned int; The current monero-wallet-rpc's blockchain height.
	// If the wallet has been offline for a long time, it may need to catch up with the daemon.
	Getheight() (height uint, err error)
	// Transfer - Send monero to a number of recipients.
	Transfer(req TransferRequest) (resp *TransferResponse, err error)
}

// New returns a new monero-wallet-rpc client.
func New(cfg Config) Client {
	cl := &client{
		addr:    cfg.Address,
		headers: cfg.CustomHeaders,
	}
	if cfg.Transport == nil {
		cl.httpcl = http.DefaultClient
	} else {
		cl.httpcl = &http.Client{
			Transport: cfg.Transport,
		}
	}
	return cl
}

type client struct {
	httpcl  *http.Client
	addr    string
	headers map[string]string
}

func (c *client) do(method string, in, out interface{}) error {
	payload, err := json2.EncodeClientRequest(method, in)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.addr, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	if c.headers != nil {
		for k, v := range c.headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := c.httpcl.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	return json2.DecodeClientResponse(resp.Body, out)
}

func (c *client) Getbalance() (balance, unlockedBalance uint64, err error) {
	jd := struct {
		Balance         uint64 `json:"balance"`
		UnlockedBalance uint64 `json:"unlocked_balance"`
	}{}
	err = c.do("getbalance", nil, &jd)
	return jd.Balance, jd.UnlockedBalance, err
}

func (c *client) Getaddress() (address string, err error) {
	jd := struct {
		Address string `json:"address"`
	}{}
	err = c.do("getaddress", nil, &jd)
	if err != nil {
		return "", err
	}
	return jd.Address, nil
}

func (c *client) Getheight() (height uint, err error) {
	jd := struct {
		Height uint `json:"height"`
	}{}
	err = c.do("getheight", nil, &jd)
	if err != nil {
		return 0, err
	}
	return jd.Height, nil
}

func (c *client) Transfer(req TransferRequest) (resp *TransferResponse, err error) {
	resp = &TransferResponse{}
	err = c.do("transfer", &req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
