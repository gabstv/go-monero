package walletrpc

import (
	"fmt"

	"github.com/gorilla/rpc/v2/json2"
)

// H is a helper map shortcut.
type H map[string]interface{}

// ErrorCode is a monero-wallet-rpc error code.
// I found them on https://github.com/monero-project/monero/blob/release-v0.11.0.0/src/wallet/wallet_rpc_server_error_codes.h
type ErrorCode int

const (
	E_UNKNOWN_ERROR          ErrorCode = -1
	E_WRONG_ADDRESS          ErrorCode = -2
	E_DAEMON_IS_BUSY         ErrorCode = -3
	E_GENERIC_TRANSFER_ERROR ErrorCode = -4
	E_WRONG_PAYMENT_ID       ErrorCode = -5
	E_TRANSFER_TYPE          ErrorCode = -6
	E_DENIED                 ErrorCode = -7
	E_WRONG_TXID             ErrorCode = -8
	E_WRONG_SIGNATURE        ErrorCode = -9
	E_WRONG_KEY_IMAGE        ErrorCode = -10
	E_WRONG_URI              ErrorCode = -11
	E_WRONG_INDEX            ErrorCode = -12
	E_NOT_OPEN               ErrorCode = -13
)

// WalletError is the error structured returned by the monero-wallet-rpc
type WalletError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (we *WalletError) Error() string {
	return fmt.Sprintf("%v: %v", we.Code, we.Message)
}

// GetWalletError checks if an erro interface is a wallet-rpc error.
func GetWalletError(err error) (isWalletError bool, werr *WalletError) {
	if err == nil {
		return false, nil
	}
	gerr, ok := err.(*json2.Error)
	if !ok {
		return false, nil
	}
	werr = &WalletError{
		Code:    ErrorCode(gerr.Code),
		Message: gerr.Message,
	}
	isWalletError = true
	return
}

// XMRToDecimal converts a raw atomic XMR balance to a more
// human readable format.
func XMRToDecimal(xmr uint64) string {
	str0 := fmt.Sprintf("%013d", xmr)
	l := len(str0)
	return str0[:l-12] + "." + str0[l-12:]
}
