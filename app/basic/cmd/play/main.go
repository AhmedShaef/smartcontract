package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ardanlabs/smartcontract/app/basic/contracts/store"
	"github.com/ardanlabs/smartcontract/business/smart"
)

func main() {
	client, privateKey, err := smart.Connect()
	if err != nil {
		log.Fatal("dial ERROR:", err)
	}

	data, err := os.ReadFile("contract.env")
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	contractID := string(data)
	fmt.Println("contractID:", contractID)

	address := common.HexToAddress(contractID)
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal("NewStore ERROR:", err)
	}

	version, err := instance.Version(nil)
	if err != nil {
		log.Fatal("version ERROR:", err)
	}
	fmt.Println("version:", version)

	// =========================================================================

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("address:", fromAddress.String())

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("pending nonce at ERROR:", err)
	}
	fmt.Println("next nonce:", nonce)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("suggest gas price ERROR:", err)
	}
	fmt.Println("suggested gas price:", smart.Wei2Eth(gasPrice))

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)    // in wei
	auth.GasLimit = uint64(30000) // The maximum amount of Gas a user can consume to conduct this transaction.
	auth.GasPrice = gasPrice      // What you are willing to pay per unit to complete this transaction.

	var key [32]byte
	var value [32]byte
	copy(key[:], []byte("name"))
	copy(value[:], []byte("ale"))

	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx sent      :", tx.Hash().Hex())
	fmt.Println("tx gas units :", tx.Gas())
	fmt.Println("tx gas price :", smart.Wei2Eth(tx.GasPrice()))
	fmt.Println("tx cost      :", smart.Wei2Eth(tx.Cost()))

	// There is a delay from the time we set to the time we see. This
	// includes changes.

	var result [32]byte
	for {
		result, err = instance.Items(nil, key)
		if err != nil {
			log.Fatal("Items ERROR:", err)
		}

		if string(result[:]) == string(value[:]) {
			break
		}

		time.Sleep(time.Second)
	}

	fmt.Println("value:", string(result[:]))
}

func privateKey() (*ecdsa.PrivateKey, error) {
	inPath := "node/keystore/UTC--2022-05-12T14-47-50.112225000Z--6327a38415c53ffb36c11db55ea74cc9cb4976fd"
	password := "123"

	data, err := ioutil.ReadFile(inPath)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(data, password)
	if err != nil {
		return nil, err
	}

	return key.PrivateKey, nil
}
