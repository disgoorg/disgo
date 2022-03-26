# httpserver

HTTPServer uses `crypto/ed25519` by default for signing verification. You can inject your own implementation by setting the `Verify` in this package.

## Example

For a simple command [`crypto/ed25519`](https://pkg.go.dev/crypto/ed25519) takes around 0.54ms on my machine and [`github.com/oasisprotocol/curve25519-voi`](https://pkg.go.dev/github.com/oasisprotocol/curve25519-voi) takes about 0.13ms.

```go
package main

import (
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/disgoorg/disgo/httpserver"
)
func main() {
	httpserver.Verify = func(publicKey httpserver.PublicKey, message, sig []byte) bool {
		return ed25519.Verify(publicKey, message, sig)
	}
}

```
