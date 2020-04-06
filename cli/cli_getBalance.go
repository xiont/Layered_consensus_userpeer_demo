package cli

import (
	"fmt"
	"github.com/corgi-kx/blockchain_golang/network"
)

func (cli *Cli) getBalance(address string) {
	balance := network.GetBalance(address)
	fmt.Printf("地址:%s的余额为：%s\n", address, balance)
}
