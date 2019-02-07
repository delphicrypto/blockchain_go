package main

import (
	"strings"
	"fmt"
	"log"
	"os"
	"bufio"
	"strconv"
	//"github.com/joho/godotenv"
)

// CLI responsible for processing command line arguments
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance ADDRESS - Get balance of ADDRESS")
	fmt.Println("  getbalances - Get balances of all addresses")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  printproblems - Print all the problems of the blockchain")
	fmt.Println("  printproblem HASH - Display problem with hash HASH")
	fmt.Println("  printlast - Print last block of the blockchain")
	fmt.Println("  printblock HEIGHT - Display block number HEIGHT")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send FROM TO AMOUNT - Send AMOUNT of coins from FROM address to TO.")
	fmt.Println("  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
	fmt.Println("  mineblock N- Mine N blocks with empty transactions. Default is 1")
	fmt.Println("  mineblockprob NODES DENSITY- Mine 1 block with empty transactions and NODES nodes and DENSITY density")
	fmt.Println("  mineblocksol HASH -  Mine 1 block with empty transactions and a solution to problem HASH")
	fmt.Println("  getdiff - Display current difficulty")
	fmt.Println("  creategraph - Create a new problem graph, with default 50 nodes and 620 edges")
	
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	//cli.validateArgs()
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		commands := strings.Fields(strings.TrimSpace(sendData))
		if len(commands) == 0 {
			continue 
		}
		command := strings.ToLower(commands[0])
		switch command {
			case "printchain":
				cli.printChain(nodeID)
			case "printlast":
				cli.printLast(nodeID)
			case "q", "quit":
				os.Exit(1)
			case "qs":
				cli.quickstart(nodeID)
			case "createwallet":
				cli.createWallet(nodeID)
			case "listaddresses":
				cli.listAddresses(nodeID)
			case "reindexutxo":
				cli.reindexUTXO(nodeID)
			case "getbalances":
				cli.getAllBalances(nodeID)
			case "getdiff":
				cli.getDifficulty(nodeID)
			case "creategraph":
				cli.createGraph(nodeID, 500, 110000)
			case "printproblems":
				cli.printProblemGraphs(nodeID)	
			case "getbalance":
				if len(commands) > 1 {
					address := commands[1]
					cli.getBalance(address, nodeID)
				 } else {
				 	fmt.Println("getbalance ADDRESS - Get balance of ADDRESS")
				 	fmt.Println("Missing argument ADDRESS")
				 }
			case "printproblem":
				if len(commands) > 1 {
					hash := commands[1]
					cli.printProblemGraph(nodeID, hash)
				 } else {
				 	fmt.Println("printproblem HASH - Display problem with hash HASH")
				 	fmt.Println("Missing argument HASH")
				 }
			case "printblock":
				if len(commands) > 1 {
					height,_ := strconv.Atoi(commands[1])
					cli.printHeight(nodeID, height)
				 } else {
				 	fmt.Println("printblock HEIGHT - Display block number HEIGHT")
				 	fmt.Println("Missing argument HEIGHT")
				 }

			case "send":
				if len(commands) > 3 {
					sendFrom := commands[1]
					sendTo   := commands[2]
					sendAmount,_ := strconv.Atoi(commands[3])
					sendMine := true
					cli.send(sendFrom, sendTo, sendAmount, nodeID, sendMine)
				 } else {
				 	fmt.Println("send FROM TO AMOUNT - Send AMOUNT of coins from FROM address to TO.")
				 	fmt.Println("Missing arguments")
				 }
				 
			case "createblockchain":
				if len(commands) > 1 {
					address := commands[1]
					cli.createBlockchain(address, nodeID)
				 } else {
				 	fmt.Println("createblockchain ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
				 	fmt.Println("Missing argument ADDRESS")
				 }
			case "startnode":
				if len(commands) > 1 {
					address := commands[1]
					cli.startNode(nodeID, address)
				 } else {
				 	fmt.Println("startnode ADDRESS - ")
				 	fmt.Println("Missing argument ADDRESS")
				 }
			case "mineblock":
				n := 1
				if len(commands) == 2 {
					m, err := strconv.Atoi(commands[1])
					if err != nil {
						fmt.Println("Invalid N argument.")
						fmt.Println("  mineblock N- Mine N blocks with empty transactions. Default is 1")
					} else {
						n = m
					}
				}
				for i := 0; i < n; i++ {
					block := cli.mineblock(nodeID)
					bc := NewBlockchain(nodeID)
					bc.AddBlock(block)
					bc.db.Close()
				}
			case "mineblockprob":
				if len(commands) == 3 {
					nodes, err1 := strconv.Atoi(commands[1])
					density, err2 := strconv.ParseFloat(commands[2], 64)
					if err1 != nil || err2 != nil {
						fmt.Println("Invalid arguments.")
						fmt.Println("  mineblockprob NODES DENSITY- Mine 1 block with empty transactions and NODES nodes and DENSITY density")
					} else {
						block :=  cli.mineblockWithNewProblem(nodeID, nodes, density)
						bc := NewBlockchain(nodeID)
						bc.AddBlock(block)
						bc.db.Close()
					}
				} else {
					fmt.Println("Invalid arguments.")
					fmt.Println("  mineblockprob NODES DENSITY- Mine 1 block with empty transactions and NODES nodes and DENSITY density")
				}
			case "mineblocksol":
				if len(commands) > 1 {
					pgHash := commands[1]
					block := cli.mineblockWithSolution(nodeID, pgHash)
					bc := NewBlockchain(nodeID)
					bc.AddBlock(block)
					bc.db.Close()
				 } else {
				 	fmt.Println("mineblocksol HASH -  Mine 1 block with empty transactions and a solution to problem HASH")
				 	fmt.Println("Missing argument HASH")
				 }
			case "minepar":
				if len(commands) > 1 {
					pgHash := commands[1]
					block := cli.mineblockParallel(nodeID, pgHash)
					bc := NewBlockchain(nodeID)
					bc.AddBlock(block)
					bc.db.Close()
				 } else {
				 	fmt.Println("mineblocksol HASH -  Mine 1 block with empty transactions and a solution to problem HASH")
				 	fmt.Println("Missing argument HASH")
				 }
				
			default:
				fmt.Println("Invalid option.")
				cli.printUsage()

		}
	}
}