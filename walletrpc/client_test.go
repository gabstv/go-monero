package walletrpc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {

	testClientGetAddress(t)
	testClientGetBalance(t)
}

func testClientGetAddress(t *testing.T) {
	//
	// server setup
	sv0 := basicTestServer([]testfn{
		func(method string, params *json.RawMessage, w http.ResponseWriter, r *http.Request) bool {
			if method == "getaddress" {
				r0 := struct {
					Address string `json:"address"`
				}{
					"45eoXYNHC4LcL2Hh42T9FMPTmZHyDEwDbgfBEuNj3RZUek8A4og4KiCfVL6ZmvHBfCALnggWtHH7QHF8426yRayLQq7MLf5",
				}
				writerpcResponseOK(&r0, w)
				return true
			}
			return false
		},
	})
	defer sv0.Close()
	//
	// test starts here
	rpccl := New(Config{
		Address: sv0.URL + "/json_rpc",
	})
	addr, err := rpccl.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, "45eoXYNHC4LcL2Hh42T9FMPTmZHyDEwDbgfBEuNj3RZUek8A4og4KiCfVL6ZmvHBfCALnggWtHH7QHF8426yRayLQq7MLf5", addr)
}

func testClientGetBalance(t *testing.T) {
	//
	// server setup
	sv0 := basicTestServer([]testfn{
		func(method string, params *json.RawMessage, w http.ResponseWriter, r *http.Request) bool {
			if method == "getbalance" {
				r0 := struct {
					Balance  uint64 `json:"balance"`
					Unlocked uint64 `json:"unlocked_balance"`
				}{
					1e12,
					1e13,
				}
				writerpcResponseOK(&r0, w)
				return true
			}
			return false
		},
	})
	defer sv0.Close()
	//
	// test starts here
	rpccl := New(Config{
		Address: sv0.URL + "/json_rpc",
	})
	balance, unlocked, err := rpccl.GetBalance()
	assert.NoError(t, err)
	// 1 XMR
	assert.Equal(t, uint64(1000000000000), balance)
	// 10 XMR
	assert.Equal(t, uint64(10000000000000), unlocked)
}

//TODO: write more server stubs
//
//

type testfn = func(method string, params *json.RawMessage, w http.ResponseWriter, r *http.Request) bool

func basicTestServer(tests []testfn) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/json_rpc" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		var c clientRequest
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		for _, v := range tests {
			if v(c.Method, c.Params, w, r) {
				return
			}
		}
		// return method not found
		writerpcResponseError(ErrUnknown, "test this in curl with the real rpc", w)
	}))
}

func writerpcResponseOK(result interface{}, w http.ResponseWriter) {
	r := &clientResponse{
		Version: "2.0",
		Result:  result,
	}
	v, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(v)
}

func writerpcResponseError(code ErrorCode, message string, w http.ResponseWriter) {
	r := &clientResponse{
		Version: "2.0",
		Result:  nil,
		Error: &WalletError{
			Code:    code,
			Message: message,
		},
	}
	v, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(v)
}

// clientRequest represents a JSON-RPC request received by the server.
type clientRequest struct {
	// JSON-RPC protocol.
	Version string `json:"jsonrpc"`
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// Object to pass as request parameter to the method.
	Params *json.RawMessage `json:"params"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	Id uint64 `json:"id"`
}

// clientResponse represents a JSON-RPC response returned to a client.
type clientResponse struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}
