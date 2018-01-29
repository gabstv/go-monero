package walletrpc

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/rpc/v2/json2"
)

// Client is a monero-wallet-rpc client.
type Client interface {
	// Return the wallet's balance.
	GetBalance() (balance, unlockedBalance uint64, err error)
	// Return the wallet's address.
	// address - string; The 95-character hex address string of the monero-wallet-rpc in session.
	GetAddress() (address string, err error)
	// GetHeight - Returns the wallet's current block height.
	// height - unsigned int; The current monero-wallet-rpc's blockchain height.
	// If the wallet has been offline for a long time, it may need to catch up with the daemon.
	GetHeight() (height uint, err error)
	// Transfer - Send monero to a number of recipients.
	Transfer(req TransferRequest) (resp *TransferResponse, err error)
	// Same as transfer, but can split into more than one tx if necessary.
	TransferSplit(req TransferRequest) (resp *TransferSplitResponse, err error)
	// Send all dust outputs back to the wallet's, to make them easier to spend (and mix).
	SweepDust() (txHashList []string, err error)
	// Send all unlocked balance to an address.
	SweepAll(req SweepAllRequest) (resp *SweepAllResponse, err error)
	// Save the blockchain.
	Store() error
	// Get a list of incoming payments using a given payment id.
	GetPayments(paymentid string) (payments []Payment, err error)
	// Get a list of incoming payments using a given payment id, or a list of
	// payments ids, from a given height. This method is the preferred method
	// over get_payments because it has the same functionality but is more extendable.
	// Either is fine for looking up transactions by a single payment ID.
	// Inputs:
	//
	//	payment_ids - array of: string
	//	min_block_height - unsigned int; The block height at which to start looking for payments.
	GetBulkPayments(paymentids []string, minblockheight uint) (payments []Payment, err error)
	// Returns a list of transfers.
	GetTransfers(req GetTransfersRequest) (resp *GetTransfersResponse, err error)
	// Show information about a transfer to/from this address.
	GetTransferByTxID(txid string) (transfer *Transfer, err error)
	// Return a list of incoming transfers to the wallet.
	IncomingTransfers(transfertype GetTransferType) (transfers []IncTransfer, err error)
	// Return the spend or view private key (or mnemonic seed).
	QueryKey(keytype QueryKeyType) (key string, err error)
	// Make an integrated address from the wallet address and a payment id.
	// payment_id - string; hex encoded; can be empty, in which case a random payment id is generated
	MakeIntegratedAddress(paymentid string) (integratedaddr string, err error)
	// Retrieve the standard address and payment id corresponding to an integrated address.
	SplitIntegratedAddress(integratedaddr string) (paymentid, address string, err error)
	// Stops the wallet, storing the current state.
	StopWallet() error
	// Create a payment URI using the official URI spec.
	MakeURI(req URIDef) (uri string, err error)
	// Parse a payment URI to get payment information.
	ParseURI(uri string) (parsed *URIDef, err error)
	// Rescan blockchain from scratch.
	RescanBlockchain() error
	// Set arbitrary string notes for transactions.
	SetTxNotes(txids, notes []string) error
	// Get string notes for transactions.
	GetTxNotes(txids []string) (notes []string, err error)
	// Sign a string.
	Sign(data string) (signature string, err error)
	// Verify a signature on a string.
	Verify(data, address, signature string) (good bool, err error)
	// Export a signed set of key images.
	ExportKeyImages() (signedkeyimages []SignedKeyImage, err error)
	// Import signed key images list and verify their spent status.
	ImportKeyImages(signedkeyimages []SignedKeyImage) (resp *ImportKeyImageResponse, err error)
	// Retrieves entries from the address book.
	// indexes - array of unsigned int; indices of the requested address book entries
	GetAddressBook(indexes []uint64) (entries []AddressBookEntry, err error)
	// Add an entry to the address book.
	AddAddressBook(entry AddressBookEntry) (index uint64, err error)
	// Delete an entry from the address book.
	DeleteAddressBook(index uint64) error
	// Rescan the blockchain for spent outputs.
	RescanSpent() error
	// Start mining in the Monero daemon.
	// Inputs:
	//
	//	threads_count - unsigned int; Number of threads created for mining
	//	do_background_mining - boolean;
	//	ignore_battery - boolean;
	StartMining(threads uint, background, ignorebattery bool) error
	// Stop mining in the Monero daemon.
	StopMining() error
	// Get a list of available languages for your wallet's seed.
	GetLanguages() (languages []string, err error)
	// Create a new wallet. You need to have set the argument "–wallet-dir" when
	// launching monero-wallet-rpc to make this work.
	// Inputs:
	//
	//   filename - string;
	//    password - string;
	//    language - string; Language for your wallets' seed.
	CreateWallet(filename, password, language string) error
	// Open a wallet. You need to have set the argument "–wallet-dir" when
	// launching monero-wallet-rpc to make this work.
	OpenWallet(filename, password string) error
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

	// in theory this is only done to catch
	// any monero related errors if
	// we are not expecting any data back
	if out == nil {
		v := &json2.EmptyResponse{}
		return json2.DecodeClientResponse(resp.Body, v)
	}
	return json2.DecodeClientResponse(resp.Body, out)
}

func (c *client) GetBalance() (balance, unlockedBalance uint64, err error) {
	jd := struct {
		Balance         uint64 `json:"balance"`
		UnlockedBalance uint64 `json:"unlocked_balance"`
	}{}
	err = c.do("getbalance", nil, &jd)
	return jd.Balance, jd.UnlockedBalance, err
}

func (c *client) GetAddress() (address string, err error) {
	jd := struct {
		Address string `json:"address"`
	}{}
	err = c.do("getaddress", nil, &jd)
	if err != nil {
		return "", err
	}
	return jd.Address, nil
}

func (c *client) GetHeight() (height uint, err error) {
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

func (c *client) TransferSplit(req TransferRequest) (resp *TransferSplitResponse, err error) {
	resp = &TransferSplitResponse{}
	err = c.do("transfer_split", &req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *client) SweepDust() (txHashList []string, err error) {
	jd := struct {
		TxHashList []string `json:"tx_hash_list"`
	}{}
	err = c.do("sweep_dust", nil, &jd)
	if err != nil {
		return nil, err
	}
	return jd.TxHashList, nil
}

func (c *client) SweepAll(req SweepAllRequest) (resp *SweepAllResponse, err error) {
	resp = &SweepAllResponse{}
	err = c.do("sweep_all", &req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *client) Store() error {
	return c.do("store", nil, nil)
}

func (c *client) GetPayments(paymentid string) (payments []Payment, err error) {
	jin := struct {
		PaymentID string `json:"payment_id"`
	}{
		paymentid,
	}
	jd := struct {
		Payments []Payment `json:"payments"`
	}{}
	err = c.do("get_payments", &jin, &jd)
	if err != nil {
		return nil, err
	}
	return jd.Payments, nil
}

func (c *client) GetBulkPayments(paymentids []string, minblockheight uint) (payments []Payment, err error) {
	jin := struct {
		PaymentIDs     []string `json:"payment_ids"`
		MinBlockHeight uint     `json:"min_block_height"`
	}{
		paymentids,
		minblockheight,
	}
	jd := struct {
		Payments []Payment `json:"payments"`
	}{}
	err = c.do("get_bulk_payments", &jin, &jd)
	if err != nil {
		return nil, err
	}
	return jd.Payments, nil
}

func (c *client) GetTransfers(req GetTransfersRequest) (resp *GetTransfersResponse, err error) {
	resp = &GetTransfersResponse{}
	err = c.do("get_transfers", &req, resp)
	return
}

func (c *client) GetTransferByTxID(txid string) (transfer *Transfer, err error) {
	jin := struct {
		TxID string `json:"txid"`
	}{
		txid,
	}
	jd := struct {
		Transfer *Transfer `json:"transfer"`
	}{}
	err = c.do("get_transfer_by_txid", &jin, &jd)
	if err != nil {
		return
	}
	transfer = jd.Transfer
	return
}

func (c *client) IncomingTransfers(transfertype GetTransferType) (transfers []IncTransfer, err error) {
	jin := struct {
		TransferType GetTransferType `json:"transfer_type"`
	}{
		transfertype,
	}
	jd := struct {
		Transfers []IncTransfer `json:"transfers"`
	}{}
	err = c.do("incoming_transfers", &jin, &jd)
	if err != nil {
		return
	}
	transfers = jd.Transfers
	return
}

func (c *client) QueryKey(keytype QueryKeyType) (key string, err error) {
	jin := struct {
		KeyType QueryKeyType `json:"key_type"`
	}{
		keytype,
	}
	jd := struct {
		Key string `json:"key"`
	}{}
	err = c.do("query_key", &jin, &jd)
	if err != nil {
		return
	}
	key = jd.Key
	return
}

func (c *client) MakeIntegratedAddress(paymentid string) (integratedaddr string, err error) {
	jin := struct {
		PaymentID string `json:"payment_id"`
	}{
		paymentid,
	}
	jd := struct {
		Address string `json:"integrated_address"`
	}{}
	err = c.do("make_integrated_address", &jin, &jd)
	if err != nil {
		return
	}
	integratedaddr = jd.Address
	return
}

func (c *client) SplitIntegratedAddress(integratedaddr string) (paymentid, address string, err error) {
	jin := struct {
		IntegratedAddress string `json:"integrated_address"`
	}{
		integratedaddr,
	}
	jd := struct {
		Address   string `json:"standard_address"`
		PaymentID string `json:"payment_id"`
	}{}
	err = c.do("split_integrated_address", &jin, &jd)
	if err != nil {
		return
	}
	paymentid = jd.PaymentID
	address = jd.Address
	return
}

func (c *client) StopWallet() error {
	return c.do("stop_wallet", nil, nil)
}

func (c *client) MakeURI(req URIDef) (uri string, err error) {
	jd := struct {
		URI string `json:"uri"`
	}{}
	err = c.do("make_uri", &req, &jd)
	if err != nil {
		return
	}
	uri = jd.URI
	return
}

func (c *client) ParseURI(uri string) (parsed *URIDef, err error) {
	jin := struct {
		URI string `json:"uri"`
	}{
		uri,
	}
	parsed = &URIDef{}
	err = c.do("parse_uri", &jin, parsed)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) RescanBlockchain() error {
	return c.do("rescan_blockchain", nil, nil)
}

func (c *client) SetTxNotes(txids, notes []string) error {
	jin := struct {
		TxIDs []string `json:"txids"`
		Notes []string `json:"notes"`
	}{
		txids,
		notes,
	}
	return c.do("set_tx_notes", &jin, nil)
}

func (c *client) GetTxNotes(txids []string) (notes []string, err error) {
	jin := struct {
		TxIDs []string `json:"txids"`
	}{
		txids,
	}
	jd := struct {
		Notes []string `json:"notes"`
	}{}
	err = c.do("get_tx_notes", &jin, &jd)
	if err != nil {
		return nil, err
	}
	notes = jd.Notes
	return
}

func (c *client) Sign(data string) (signature string, err error) {
	jin := struct {
		Data string `json:"data"`
	}{
		data,
	}
	jd := struct {
		Signature string `json:"signature"`
	}{}
	err = c.do("sign", &jin, &jd)
	if err != nil {
		return "", err
	}
	signature = jd.Signature
	return
}

func (c *client) Verify(data, address, signature string) (good bool, err error) {
	jin := struct {
		Data      string `json:"data"`
		Address   string `json:"address"`
		Signature string `json:"signature"`
	}{
		data,
		address,
		signature,
	}
	jd := struct {
		Good bool `json:"good"`
	}{}
	err = c.do("verify", &jin, &jd)
	if err != nil {
		return false, err
	}
	good = jd.Good
	return
}

func (c *client) ExportKeyImages() (signedkeyimages []SignedKeyImage, err error) {
	jd := struct {
		SignedKeyImages []SignedKeyImage `json:"signed_key_images"`
	}{}
	err = c.do("export_key_images", nil, &jd)
	signedkeyimages = jd.SignedKeyImages
	return
}

func (c *client) ImportKeyImages(signedkeyimages []SignedKeyImage) (resp *ImportKeyImageResponse, err error) {
	jin := struct {
		SignedKeyImages []SignedKeyImage `json:"signed_key_images"`
	}{
		signedkeyimages,
	}
	resp = &ImportKeyImageResponse{}
	err = c.do("import_key_images", &jin, resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetAddressBook(indexes []uint64) (entries []AddressBookEntry, err error) {
	jin := struct {
		Indexes []uint64 `json:"entries"`
	}{
		indexes,
	}
	jd := struct {
		Entries []AddressBookEntry `json:"entries"`
	}{}
	err = c.do("get_address_book", &jin, &jd)
	if err != nil {
		return nil, err
	}
	entries = jd.Entries
	return
}

func (c *client) AddAddressBook(entry AddressBookEntry) (index uint64, err error) {
	entry.Index = 0
	jd := struct {
		Index uint64 `json:"index"`
	}{}
	err = c.do("add_address_book", &entry, &jd)
	if err != nil {
		return 0, err
	}
	index = jd.Index
	return
}

func (c *client) DeleteAddressBook(index uint64) error {
	jin := struct {
		Index uint64 `json:"index"`
	}{
		index,
	}
	return c.do("delete_address_book", &jin, nil)
}

func (c *client) RescanSpent() error {
	return c.do("rescan_spent", nil, nil)
}

func (c *client) StartMining(threads uint, background, ignorebattery bool) error {
	jin := struct {
		Threads       uint `json:"threads_count"`
		Background    bool `json:"do_background_mining"`
		IgnoreBattery bool `json:"ignore_battery"`
	}{
		threads,
		background,
		ignorebattery,
	}
	return c.do("start_mining", &jin, nil)
}

func (c *client) StopMining() error {
	return c.do("stop_mining", nil, nil)
}

func (c *client) GetLanguages() (languages []string, err error) {
	jd := struct {
		Languages []string `json:"languages"`
	}{}
	err = c.do("get_languages", nil, &jd)
	if err != nil {
		return nil, err
	}
	languages = jd.Languages
	return
}

func (c *client) CreateWallet(filename, password, language string) error {
	jin := struct {
		Filename string `json:"filename"`
		Password string `json:"password"`
		Language string `json:"language"`
	}{
		filename,
		password,
		language,
	}
	return c.do("create_wallet", &jin, nil)
}

func (c *client) OpenWallet(filename, password string) error {
	jin := struct {
		Filename string `json:"filename"`
		Password string `json:"password"`
	}{
		filename,
		password,
	}
	return c.do("open_wallet", &jin, nil)
}
