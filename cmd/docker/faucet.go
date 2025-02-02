package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/randymcmillan/plebnet-playground-rpc/internal/config"
	"github.com/randymcmillan/plebnet-playground-rpc/internal/docker"
)

var faucet = cli.Command{
	Name:      "faucet",
	Usage:     "generate and send bitcoin to given address",
	ArgsUsage: "<address> [amount] [asset]",
	Action:    faucetAction,
	Flags: []cli.Flag{
		&liquidFlag,
	},
}

func faucetAction(ctx *cli.Context) error {

	if isRunning, _ := nigiriState.GetBool("running"); !isRunning {
		return errors.New("nigiri is not running")
	}

	if ctx.NArg() < 1 || ctx.NArg() > 3 {
		return errors.New("wrong number of arguments")
	}

	isLiquid := ctx.Bool("liquid")
	datadir := ctx.String("datadir")
	composePath := filepath.Join(datadir, config.DefaultCompose)

	var serviceName string = "chopsticks"
	if isLiquid {
		serviceName = "chopsticks-liquid"
	}

	portSlice, err := docker.GetPortsForService(composePath, serviceName)
	if err != nil {
		return err
	}
	mappedPorts := strings.Split(portSlice[0], ":")

	request := map[string]interface{}{
		"address": ctx.Args().First(),
	}
	
	if ctx.Args().Len() >= 2 {
		amountFloat, err := strconv.ParseFloat(ctx.Args().Get(1), 64)
		if err != nil {
			return fmt.Errorf("invalid amount: %v", err)
		}
		request["amount"] = amountFloat
	}

	if isLiquid && ctx.Args().Len() == 3 {
		request["asset"] = ctx.Args().Get(2)
	}
	
	requestPort := mappedPorts[0]
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}
	res, err := http.Post("http://127.0.0.1:"+requestPort+"/faucet", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(string(data))
	}

	var dat map[string]string
	if err := json.Unmarshal([]byte(data), &dat); err != nil {
		return errors.New("internal error, please try again")
	}
	if dat["txId"] == "" {
		return errors.New("not successful")
	}
	fmt.Println("txId: " + dat["txId"])

	return nil
}
