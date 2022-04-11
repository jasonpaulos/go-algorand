package main

import (
	"fmt"
	"log"
	"os"

	"github.com/algorand/go-algorand/data/avmdbg"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "amvdbg",
	Short: "AVM Debugger",
	Long:  `Debug an Algorand Virtual Machine input`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params, err := makeDebuggerParams()
		if err != nil {
			return fmt.Errorf("Error parsing inputs: %w", err)
		}

		ctx := avmdbg.MakeContext(&params)

		return mainLoop(ctx)

		//If no arguments passed, we should fallback to help
		// cmd.HelpFunc()(cmd, args)
	},
}

var flagTxnInputFile string
var flagAlgodURL string
var flagAlgodToken string

// var flagIgnoreTxnSignatures bool

func init() {
	rootCmd.Flags().StringVar(&flagTxnInputFile, "txns", "", "Path to transaction input file")
	rootCmd.Flags().StringVar(&flagAlgodURL, "algod-url", "", "Address of algod node to connect to")
	rootCmd.Flags().StringVar(&flagAlgodToken, "algod-token", "", "Token to use when connecting to algod node")
	// rootCmd.Flags().BoolVar(&flagIgnoreTxnSignatures, "ignore-signatures", false, "Do not validate transaction signatures")

	rootCmd.MarkFlagRequired("txns")
	rootCmd.MarkFlagRequired("algod-url")
}
