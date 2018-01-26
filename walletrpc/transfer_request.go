package walletrpc

// TransferRequest is the request body of the Transfer client rpc call.
type TransferRequest struct {
	// Destinations - array of destinations to receive XMR:
	Destinations []Destination `json:"destinations"`
	// Fee - unsigned int; Ignored, will be automatically calculated.
	Fee uint64 `json:"fee,omitempty"`
	// Mixin - unsigned int; Number of outpouts from the blockchain to mix with (0 means no mixing).
	Mixin uint64 `json:"mixin"`
	// unlock_time - unsigned int; Number of blocks before the monero can be spent (0 to not add a lock).
	UnlockTime uint64 `json:"unlock_time"`
	// payment_id - string; (Optional) Random 32-byte/64-character hex string to identify a transaction.
	PaymentID string `json:"payment_id,omitempty"`
	// get_tx_key - boolean; (Optional) Return the transaction key after sending.
	GetTxKey bool `json:"get_tx_key"`
	// priority - unsigned int; Set a priority for the transaction.
	// Accepted Values are: 0-3 for: default, unimportant, normal, elevated, priority.
	Priority Priority `json:"priority"`
	// do_not_relay - boolean; (Optional) If true, the newly created transaction will not be relayed to the monero network. (Defaults to false)
	DoNotRelay bool `json:"do_not_relay,omitempty"`
	// get_tx_hex - boolean; Return the transaction as hex string after sending
	GetTxHex bool `json:"get_tx_hex,omitempty"`
}

// Destination to receive XMR
type Destination struct {
	// Amount - unsigned int; Amount to send to each destination, in atomic units.
	Amount uint64 `json:"amount"`
	// Address - string; Destination public address.
	Address string `json:"address"`
}

// TransferResponse is the successful output of a Client.Transfer()
type TransferResponse struct {
	// fee - Integer value of the fee charged for the txn.
	Fee uint64 `json:"fee"`
	// tx_hash - String for the publically searchable transaction hash
	TxHash string `json:"tx_hash"`
	// tx_key - String for the transaction key if get_tx_key is true, otherwise, blank string.
	TxKey string `json:"tx_key,omitempty"`
	// tx_blob - Transaction as hex string if get_tx_hex is true
	TxBlob string `json:"tx_blob,omitempty"`
}

// TransferSplitResponse is the successful output of a Client.TransferSplit()
type TransferSplitResponse struct {
	// fee_list - array of: integer. The amount of fees paid for every transaction.
	FeeList []uint64 `json:"fee_list"`
	// tx_hash_list - array of: string. The tx hashes of every transaction.
	TxHashList []string `json:"tx_hash_list"`
	// tx_blob_list - array of: string. The tx as hex string for every transaction.
	TxBlobList []string `json:"tx_blob_list"`
	// amount_list - array of: integer. The amount transferred for every transaction..
	AmountList []uint64 `json:"amount_list"`
	// tx_key_list - array of: string. The transaction keys for every transaction.
	TxKeyList []string `json:"tx_key_list"`
}

// SweepAllRequest is the struct to send all unlocked balance to an address.
type SweepAllRequest struct {
	// address - string; Destination public address.
	Address string `json:"address"`
	// priority - unsigned int; (Optional)
	Priority Priority `json:"priority,omitempty"`
	// mixin - unsigned int; Number of outpouts from the blockchain to mix with (0 means no mixing).
	Mixin uint64 `json:"mixin"`
	// unlock_time - unsigned int; Number of blocks before the monero can be spent (0 to not add a lock).
	UnlockTime uint64 `json:"unlock_time"`
	// payment_id - string; (Optional) Random 32-byte/64-character hex string to identify a transaction.
	PaymentID string `json:"payment_id,omitempty"`
	// get_tx_keys - boolean; (Optional) Return the transaction keys after sending.
	GetTxKeys bool `json:"get_tx_keys,omitempty"`
	// below_amount - unsigned int; (Optional)
	BelowAmount uint64 `json:"below_amount"`
	// do_not_relay - boolean; (Optional)
	DoNotRelay bool `json:"do_not_relay,omitempty"`
	// get_tx_hex - boolean; (Optional) return the transactions as hex encoded string.
	GetTxHex bool `json:"get_tx_hex,omitempty"`
}

// SweepAllResponse is a tipical response of a SweepAllRequest
type SweepAllResponse struct {
	// tx_hash_list - array of: string. The tx hashes of every transaction.
	TxHashList []string `json:"tx_hash_list"`
	// tx_blob_list - array of: string. The tx as hex string for every transaction.
	TxBlobList []string `json:"tx_blob_list"`
	// tx_key_list - array of: string. The transaction keys for every transaction.
	TxKeyList []string `json:"tx_key_list"`
}

// Payment ...
type Payment struct {
	PaymentID   string `json:"payment_id"`
	TxHash      string `json:"tx_hash"`
	Amount      uint64 `json:"amount"`
	BlockHeight uint64 `json:"block_height"`
	UnlockTime  uint64 `json:"unlock_time"`
}

// GetTransfersRequest = GetTransfers body
type GetTransfersRequest struct {
	In             bool   `json:"in"`
	Out            bool   `json:"out"`
	Pending        bool   `json:"pending"`
	Failed         bool   `json:"failed"`
	Pool           bool   `json:"pool"`
	FilterByHeight bool   `json:"filter_by_height"`
	MinHeight      uint64 `json:"min_height"`
	MaxHeight      uint64 `json:"max_height"`
}

// GetTransfersResponse = GetTransfers output
type GetTransfersResponse struct {
	In      []Transfer `json:"in"`
	Out     []Transfer `json:"out"`
	Pending []Transfer `json:"pending"`
	Failed  []Transfer `json:"failed"`
	Pool    []Transfer `json:"pool"`
}

// Transfer is the transfer data of
type Transfer struct {
	TxID         string        `json:"txid"`
	PaymentID    string        `json:"payment_id"`
	Height       uint64        `json:"height"`
	Timestamp    uint64        `json:"timestamp"`
	Amount       uint64        `json:"amount"`
	Fee          uint64        `json:"fee"`
	Note         string        `json:"note"`
	Destinations []Destination `json:"destinations,omitempty"` // TODO: check if deprecated
	Type         string        `json:"type"`
}

// IncTransfer is returned by IncomingTransfers
type IncTransfer struct {
	Amount uint64 `json:"amount"`
	Spent  bool   `json:"spent"`
	// Mostly internal use, can be ignored by most users.
	GlobalIndex uint64 `json:"global_index"`
	// Several incoming transfers may share the same hash
	// if they were in the same transaction.
	TxHash string `json:"tx_hash"`
	TxSize uint64 `json:"tx_size"`
}
