package main

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"fmt"
	"math/rand"
	"flag"
)

const defaultWalletPath = "/tmp/ethwallet/"
const defaultWalletPassword = "123qwe"

func newTestAccount(walletPath string, walletPassword string) {
	wallet := keystore.NewKeyStore(walletPath,
		keystore.LightScryptN,keystore.LightScryptP)

	account, err := wallet.NewAccount(walletPassword)
	if err != nil {
		fmt.Println("Account error:",err)
		os.Exit(1)
	}
	fmt.Println(account.Address.Hex())
}

func main() {

	walletPath := flag.String("walletPath",defaultWalletPath,"Wallet storage directory")
	rpcUrl := flag.String("rpcUrl","http://localhost:8601","Geth json rpc or ipc url")
	numOfTx := flag.Int64("numOfTx",1000,"Number of transactions")
	toAddress := flag.String("toAddress","0xa363B708d266Cc7D5a07f8ce1f7f19b407f1a301","to address of transaction")
	newAccount := flag.Bool("newAccount",false,"Create new account")
	getBalance := flag.Bool("getBalance",false,"Get balance for the default account")

	flag.Parse()

	if *newAccount == true {
		newTestAccount(*walletPath,defaultWalletPassword)
		os.Exit(0)
	}

	wallet := keystore.NewKeyStore(*walletPath,
		keystore.LightScryptN,keystore.LightScryptP)
	if len(wallet.Accounts()) == 0 {
		fmt.Println("Empty wallet, create account first.")
		os.Exit(2)
	}
	account := wallet.Accounts()[0]
	wallet.Unlock(account, defaultWalletPassword)

	fmt.Println("account address:",account.Address.Hex())
	client,err := ethclient.Dial(*rpcUrl)

	if err != nil {
		fmt.Println("Dial error:",err)
		os.Exit(1)
	}
	//d := time.Now().Add(1000 * time.Millisecond)
	//ctx, cancel := context.WithDeadline(context.Background(),d)
	//defer cancel()
	ctx := context.Background()
	balance, err := client.BalanceAt(ctx,account.Address,nil)
	if err != nil {
		fmt.Println("Balance error:",err)
		os.Exit(1)
	}
	fmt.Println("Balance: ",balance.String())

	if *getBalance == true {
		os.Exit(0)
	}

	nonce, _ := client.NonceAt(ctx,account.Address,nil)
	fmt.Println("nonce: ", nonce)
	//toAddress := "0xfe38D1319C60cfBD89a969AB36160b4B73922096"
	var gasLimit uint64  = 200000
	gasPrice := big.NewInt(0)
	//nonce +=100
	//nonce++
	for i := int64(0); i< *numOfTx; i++ {
		amount := big.NewInt(rand.Int63n(10000))
		fmt.Println("nonce: ", nonce)
		tx := types.NewTransaction(nonce, common.HexToAddress(*toAddress), amount, gasLimit, gasPrice, nil)
		signTx, err := wallet.SignTx(account, tx, nil)
		//fmt.Println(signTx)
		err = client.SendTransaction(ctx, signTx)
		if err != nil {
			fmt.Println("err:", err)
			os.Exit(1)
		}
		nonce++
	}
}

