package main

import (
	"flag"
	"fmt"
	"github.com/pako8128/unchained"
	"os"
	"strconv"
)

type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
}

func (cli *CLI) Run() {
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	createBlockchainData := createBlockchainCmd.String("address", "", "address that gets the revenue for mining the first block")
	getBalanceData := getBalanceCmd.String("address", "", "Wallet address")
	sendCmdSender := sendCmd.String("from", "", "sets the Sender")
	sendCmdReceiver := sendCmd.String("to", "", "sets the Receiver")
	sendCmdAmount := sendCmd.Int("amount", 0, "sets the amount")

	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "createblockchain":
		createBlockchainCmd.Parse(os.Args[2:])
	case "printchain":
		printChainCmd.Parse(os.Args[2:])
	case "getbalance":
		getBalanceCmd.Parse(os.Args[2:])
	case "send":
		sendCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainData == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceData)
	}

	if sendCmd.Parsed() {
		if *sendCmdSender == "" {
			sendCmd.Usage()
			os.Exit(1)
		}
		if *sendCmdReceiver == "" {
			sendCmd.Usage()
			os.Exit(1)
		}
		if *sendCmdAmount == 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendCmdSender, *sendCmdReceiver, *sendCmdAmount)
	}
}

func (cli *CLI) createBlockchain(address string) {
	bc := unchained.CreateBlockchain(address)
	bc.Close()
	fmt.Println("Done!")
}

func (cli *CLI) printChain() {
	bc := unchained.NewBlockchain("")

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := unchained.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBalance(address string) {
	bc := unchained.NewBlockchain(address)
	defer bc.Close()

	balance := 0
	UTXOs := bc.FindUTXOs(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance for %s is %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	bc := unchained.NewBlockchain(from)
	defer bc.Close()

	tx := unchained.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*unchained.Transaction{tx})
	fmt.Println("Success!")
}

func main() {
	cli := CLI{}
	cli.Run()
}
