package cli

import (
	"github.com/corgi-kx/blockchain_golang/network"
)

func (cli Cli) startWebsocketServer() {
	network.StartWebsocketServer(cli)
}
