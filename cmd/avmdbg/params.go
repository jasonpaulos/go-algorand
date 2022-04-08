package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	algod "github.com/algorand/go-algorand/daemon/algod/api/client"
	"github.com/algorand/go-algorand/data/avmdbg"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/protocol"
)

func makeDebuggerParams() (avmdbg.DebuggerParams, error) {
	params := avmdbg.DebuggerParams{
		CheckTxnSignatures: !flagIgnoreTxnSignatures,
	}

	if flagTxnInputFile != "" {
		data, err := ioutil.ReadFile(flagTxnInputFile)
		if err != nil {
			return avmdbg.DebuggerParams{}, fmt.Errorf("Could not read file %s: %w", flagTxnInputFile, err)
		}

		txns, err := decodeSignedTxns(data)
		if err != nil {
			return avmdbg.DebuggerParams{}, fmt.Errorf("Could not decode transactions from file %s: %w", flagTxnInputFile, err)
		}
		params.InputTxns = txns
	}

	if flagAlgodURL != "" {
		parsedURL, err := url.Parse(flagAlgodURL)
		if err != nil {
			return avmdbg.DebuggerParams{}, fmt.Errorf("Could not parse algod-url %s: %w", flagAlgodURL, err)
		}
		params.AlgodClient = algod.MakeRestClient(*parsedURL, flagAlgodToken)
		params.AlgodClient.SetAPIVersionAffinity(algod.APIVersionV2)
		// params.algodClient.HealthCheck()
	} else if flagAlgodToken != "" {
		return avmdbg.DebuggerParams{}, errors.New("algod-token cannot be specified without algod-url")
	}

	return params, nil
}

func decodeSignedTxns(data []byte) ([]transactions.SignedTxn, error) {
	var txnGroup []transactions.SignedTxn

	// 1. Attempt json - a single transactions.SignedTxn
	var txn transactions.SignedTxn
	err1 := protocol.DecodeJSON(data, &txn)
	if err1 == nil {
		txnGroup = append(txnGroup, txn)
		return txnGroup, nil
	}

	// 2. Attempt json - array of transactions.SignedTxn
	err2 := protocol.DecodeJSON(data, &txnGroup)
	if err2 == nil {
		return txnGroup, nil
	}

	// 3. Attempt msgp - array of transactions.SignedTxn
	var err3 error
	dec := protocol.NewDecoderBytes(data)
	for {
		var txn transactions.SignedTxn
		err3 = dec.Decode(&txn)
		if err3 == io.EOF {
			err3 = nil
			break
		}
		if err3 != nil {
			break
		}
		txnGroup = append(txnGroup, txn)
	}

	if err3 == nil {
		return txnGroup, nil
	}

	// TODO: include messages from err1, err2, err3
	return nil, errors.New("Could not decode transactions")
}
