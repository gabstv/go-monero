# go-monero

This package is a hub of monero related tools for Go. At this time, only the Wallet RPC Client is available.

## Wallet RPC Client

[![GoDoc](https://godoc.org/github.com/gabstv/go-monero/walletrpc?status.svg)](https://godoc.org/github.com/gabstv/go-monero/walletrpc)

The ```go-monero/walletrpc``` package is a RPC client with all the methods of the v0.11.0.0 release.
It does support digest authentication, [however I don't recommend using it alone (without https).](https://en.wikipedia.org/wiki/Digest_access_authentication#Disadvantages) If there is a need to split the RPC client and server into separate instances, you could put a proxy on the instance that contains the RPC server and check the authenticity of the requests using https + X-API-KEY headers between the proxy and this RPC client (there is an example about this implementation below)

### Installation

```sh
go get -u github.com/gabstv/go-monero/walletrpc
```

### Usage

The simplest way to use walletrpc is if you have both the server (monero-wallet-rpc) and the client on the same machine.

#### Running monero-wallet-rpc:

```sh
monero-wallet-rpc --testnet --wallet-file ~/testnet/mywallet.bin --rpc-bind-port 18082 --disable-rpc-login
```

Go:

```Go
package main

import (
	"fmt"
	"os"

	"github.com/gabstv/go-monero/walletrpc"
	"github.com/gabstv/go-monero/walletrpc/unit"
)

func main() {
	// Start a wallet client instance
	client := walletrpc.New(walletrpc.Config{
		Address: "http://127.0.0.1:18082/json_rpc",
	})

	// check wallet balance
	balance, unlocked, err := client.Getbalance()

	// there are two types of error that can happen:
	//   connection errors
	//   monero wallet errors
	// connection errors are pretty much unicorns if everything is on the
	// same instance (unless your OS hit an open files limit or something)
	if err != nil {
		if iswerr, werr := walletrpc.GetWalletError(err); iswerr {
			// it is a monero wallet error
			fmt.Printf("Wallet error (id:%v) %v\n", werr.Code, werr.Message)
			os.Exit(1)
		}
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Balance:", walletrpc.XMRToDecimal(balance))
	fmt.Println("Unlocked balance:", walletrpc.XMRToDecimal(unlocked))

	// Make a transfer
	res, err := client.Transfer(walletrpc.TransferRequest{
		Destinations: []walletrpc.Destination{
			{
				Address: "45eoXYNHC4LcL2Hh42T9FMPTmZHyDEwDbgfBEuNj3RZUek8A4og4KiCfVL6ZmvHBfCALnggWtHH7QHF8426yRayLQq7MLf5",
				Amount:  10*unit.Millinero, // 0.01 XMR
			},
		},
		Priority: walletrpc.PriorityUnimportant,
		Mixin:    1,
	})
	if err != nil {
		if iswerr, werr := walletrpc.GetWalletError(err); iswerr {
			// insufficient funds return a monero wallet error
			// walletrpc.ErrGenericTransferError
			fmt.Printf("Wallet error (id:%v) %v\n", werr.Code, werr.Message)
			os.Exit(1)
		}
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Transfer success! Fee:", walletrpc.XMRToDecimal(res.Fee), "Hash:", res.TxHash)
}
```

### Using Digest Authentication

```sh
monero-wallet-rpc --testnet --rpc-bind-ip 127.0.0.1 --rpc-bind-port 29567 --rpc-login john:doe --wallet-file ~/testnet/wallet_03.bin
```

```Go
package main

import (
	"fmt"

	"github.com/gabstv/go-monero/walletrpc"
	"github.com/gabstv/httpdigest"
)

func main() {
	// username: john
	// password: doe
	t := httpdigest.New("john", "doe")

	client := walletrpc.New(walletrpc.Config{
		Address:   "http://127.0.0.1:29567/json_rpc",
		Transport: t,
	})

	balance, unlocked, err := client.Getbalance()

	if err != nil {
		panic(err)
	}
	fmt.Println("balance", walletrpc.XMRToDecimal(balance))
	fmt.Println("unlocked balance", walletrpc.XMRToDecimal(unlocked))
}
```

### Using a proxy

You can use a proxy to be in between this client and the monero RPC server. This way you can use a safe encryption tunnel around the network.

#### Starting the RPC server

```sh
monero-wallet-rpc --testnet --wallet-file ~/testnet/mywallet.bin --rpc-bind-port 18082 --disable-rpc-login
```

#### Starting a proxy server

This example uses sandpiper (```github.com/gabstv/sandpiper/sandpiper```) but you could also use nginx or apache

sandpiper config.yml:
```yaml
debug: true
#listen_addr:     :8084
listen_addr_tls: :23456
fallback_domain: moneroproxy
routes:
  - 
    domain:        moneroproxy
    out_conn_type: HTTP
    out_addr:      localhost:18082
    auth_mode:  apikey
    auth_key:   X-API-KEY
    auth_value: 55c12fca1b994455d3ec1795bdc82cca
    tls_cert_file: moneroproxy.cert.pem
    tls_key_file:  moneroproxy.key.pem
```

The Go program is similar, but it uses an API-KEY:

```Go
package main

import (
	"fmt"
    "os"
    "net/http"
    "crypto/tls"

	"github.com/gabstv/go-monero/walletrpc"
)

func main() {
	// Start a wallet client instance
	client := walletrpc.New(walletrpc.Config{
        Address: "http://127.0.0.1:23456/json_rpc",
        CustomHeaders: map[string]string{
			"X-API-KEY": "55c12fca1b994455d3ec1795bdc82cca", // we use the same key defined above
        },
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // WARNING: instead of this, you can
				// provide (or install) a certificate to make it
				// really secure with Certificates: []tls.Certificate{},
			},
		},
	})

	// check wallet balance
	balance, unlocked, err := client.Getbalance()

	// there are two types of error that can happen:
	//   connection errors
	//   monero wallet errors
	// connection errors are pretty much unicorns if everything is on the
	// same instance (unless your OS hit an open files limit or something)
	if err != nil {
		if iswerr, werr := walletrpc.GetWalletError(err); iswerr {
			// it is a monero wallet error
			fmt.Printf("Wallet error (id:%v) %v\n", werr.Code, werr.Message)
			os.Exit(1)
		}
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Balance:", walletrpc.XMRToDecimal(balance))
    fmt.Println("Unlocked balance:", walletrpc.XMRToDecimal(unlocked))
}
```

# Contributing

* You can contribute with pull requests.
* If could donate some XMR to the address below if you're feeling generous:

```
45eoXYNHC4LcL2Hh42T9FMPTmZHyDEwDbgfBEuNj3RZUek8A4og4KiCfVL6ZmvHBfCALnggWtHH7QHF8426yRayLQq7MLf5
```
