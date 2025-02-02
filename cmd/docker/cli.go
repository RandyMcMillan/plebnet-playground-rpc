package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

var command = cli.Command{
	Name:   "cli",
	Usage:  "invoke bitcoin-cli",
	Action: cliAction,
	Flags: []cli.Flag{
		&liquidFlag,
		&cli.StringFlag{
			Name:  "rpcwallet",
			Usage: "rpcwallet to be used for node JSONRPC commands",
			Value: "",
		},
	},
}

func cliAction(ctx *cli.Context) error {

	if isRunning, _ := nigiriState.GetBool("running"); !isRunning {
		return errors.New("nigiri is not running")
	}

    isLiquid := ctx.Bool("liquid")
    rpcWallet := ctx.String("playground-wallet")
    network := ctx.String("signet")
    //datadir := ctx.String("/root/.bitcoin")

    rpcArgs := []string{"exec", "playground-bitcoind", "bitcoin-cli", "-" + network, "-rpcwallet=" + rpcWallet }
	if isLiquid {
		rpcArgs = []string{"exec", "liquid", "elements-cli", "-datadir=config", "-rpcwallet=" + rpcWallet}
	}
	cmdArgs := append(rpcArgs, ctx.Args().Slice()...)
	bashCmd := exec.Command("docker", cmdArgs...)
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr

	if err := bashCmd.Run(); err != nil {
		return err
	}

	return nil
}
