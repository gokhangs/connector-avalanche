package main

import (
	"context"

	"fmt"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/accounts/keystore"
	// "github.com/ethereum/go-ethereum/common/hexutil"
	//"github.com/ethereum/go-ethereum/core/types"
)

// var avalancheTestnetCchain = "https://rpc.ankr.com/avalanche_fuji-c"
var avalancheTestnetCchain2 = "https://api.avax-test.network/ext/bc/C/rpc"

//var avalancheTestnetCchain3 = "https://api.avax-test.network:443"

func main() {
	//connect
	client, err := ethclient.DialContext(context.Background(), avalancheTestnetCchain2)
	if err != nil {
		log.Fatalf("Error while creating ether client for avalanche testnet chain %v", err)
	}
	defer client.Close()

	// block, err := client.BlockByNumber(context.Background(), nil) //to get last block of the avalanche blockchain
	// if err != nil {
	// 	log.Fatalf("Error to get block %v", err)
	// }
	// fmt.Println(block.Number())

	addr := "0x8bda0c667503ec64216ebcc906341ed468873d0d"
	addr2 := "0xf23d812bb3a57478d681d6ccbc50469ca237c05e"

	addr1_privKey := "c50dfc42f8661b965a56bd7c305d34ed20c997591e9305df9d22ee6e1b9035e7"

	add1 := common.HexToAddress(addr)
	add2 := common.HexToAddress(addr2)

	balance1, err := client.BalanceAt(context.Background(), add1, nil)
	if err != nil {
		log.Fatal(err)
	}
	floatBalance := new(big.Float)
	floatBalance.SetString(balance1.String())

	balance2, err := client.BalanceAt(context.Background(), add2, nil)
	if err != nil {
		log.Fatal(err)
	}

	floatBalance2 := new(big.Float)
	floatBalance2.SetString(balance2.String())

	// value := new(big.Float).Quo(floatBalance, big.NewFloat(math.Pow10(18)))

	fmt.Println("pre-transfer balance, addr1: ", new(big.Float).Quo(floatBalance, big.NewFloat(math.Pow10(18))))
	fmt.Println("pre-transfer balance, addr2: ", new(big.Float).Quo(floatBalance2, big.NewFloat(math.Pow10(18))))

	nonce, err := client.PendingNonceAt(context.Background(), add1) //hexutil.DecodeUint64("0x1216")
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// gaslimit, err := client.SuggestGasTipCap(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	v, err := hexutil.DecodeBig("0x1b")
	if err != nil {
		log.Fatal(err)
	}

	r, err := hexutil.DecodeBig("0x56b5bf9222ce26c3239492173249696740bc7c28cd159ad083a0f4940baf6d03")
	if err != nil {
		panic(err)
	}
	s, err := hexutil.DecodeBig("0x5fcd608b3b638950d3fe007b19ca8c4ead37237eaf89a8426777a594fd245c2a")
	if err != nil {
		panic(err)
	}

	data := []byte("This a data that is extremely useless. Really, really useless.")

	//new transaction between accounts
	txData := types.TxData(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      21992,
		To:       &add2,
		Value:    big.NewInt(1),
		Data:     data,
		V:        v,
		R:        r,
		S:        s,
	})
	tx := types.NewTx(txData)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(addr1_privKey)
	if err != nil {
		log.Fatal(err)

	}

	txSgined, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), txSgined)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("tx sent: ", txSgined.Hash().Hex())
	time.Sleep(20 * time.Second)
	receipt, err := client.TransactionReceipt(context.Background(), txSgined.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx receipt:", receipt)

	balance1, err = client.BalanceAt(context.Background(), add1, nil)
	if err != nil {
		log.Fatal(err)
	}
	floatBalance = new(big.Float)
	floatBalance.SetString(balance1.String())

	balance2, err = client.BalanceAt(context.Background(), add2, nil)
	if err != nil {
		log.Fatal(err)
	}

	floatBalance2 = new(big.Float)
	floatBalance2.SetString(balance2.String())

	fmt.Println("post-transfer balance, addr1: ", new(big.Float).Quo(floatBalance, big.NewFloat(math.Pow10(18))))
	fmt.Println("post-transfer balance, addr2: ", new(big.Float).Quo(floatBalance2, big.NewFloat(math.Pow10(18))))
}
