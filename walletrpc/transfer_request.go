package walletrpc

/*
{
	"destinations":[
		{
			"amount":100000000,
			"address":"9wNgSYy2F9qPZu7KBjvsFgZLTKE2TZgEpNFbGka9gA5zPmAXS35QzzYaLKJRkYTnzgArGNX7TvSqZC87tBLwtaC5RQgJ8rm"
		},
		{
			"amount":200000000,
			"address":"9vH5D7Fv47mbpCpdcthcjU34rqiiAYRCh1tYywmhqnEk9iwCE9yppgNCXAyVHG5qJt2kExa42TuhzQfJbmbpeGLkVbg8xit"
		}
	],
	"mixin":4,
	"get_tx_key": true
}
*/

// TransferRequest is the request body of the Transfer client rpc call.
type TransferRequest struct {
	// Destinations - array of destinations to receive XMR:
	Destinations []Destination `json:"destinations"`
	// Fee - unsigned int; Ignored, will be automatically calculated.
	Fee uint64 `json:"fee,omitempty"`
	// Mixin - unsigned int; Number of outpouts from the blockchain to mix with (0 means no mixing).
	Mixin uint `json:"mixin"`
	// unlock_time - unsigned int; Number of blocks before the monero can be spent (0 to not add a lock).
	UnlockTime uint `json:"unlock_time"`
	// payment_id - string; (Optional) Random 32-byte/64-character hex string to identify a transaction.
	PaymentId string `json:"payment_id,omitempty"`
	// get_tx_key - boolean; (Optional) Return the transaction key after sending.
	GetTxKey bool `json:"get_tx_key"`
	// priority - unsigned int; Set a priority for the transaction. Accepted Values are: 0-3 for: default, unimportant, normal, elevated, priority.
	Priority uint `json:"priority"`
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
	Fee int64 `json:"fee"`
	// tx_hash - String for the publically searchable transaction hash
	TxHash string `json:"tx_hash"`
	// tx_key - String for the transaction key if get_tx_key is true, otherwise, blank string.
	TxKey string `json:"tx_key,omitempty"`
	// tx_blob - Transaction as hex string if get_tx_hex is true
	TxBlob string `json:"tx_blob,omitempty"`
}
