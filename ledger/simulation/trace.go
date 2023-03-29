// Copyright (C) 2019-2023 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package simulation

import (
	"fmt"

	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/ledger/ledgercore"
)

// TxnPath is a "transaction path": e.g. [0, 0, 1] means the second inner txn of the first inner txn of the first txn.
// You can use this transaction path to find the txn data in the `TxnResults` list.
type TxnPath []uint64

// TxnResult contains the simulation result for a single transaction
type TxnResult struct {
	Txn                    transactions.SignedTxnWithAD
	MissingSignature       bool
	AppBudgetConsumed      uint64
	LogicSigBudgetConsumed uint64

	// Trace is the optional, highly detailed evaluation trace
	Trace TransactionTrace
}

// TxnGroupResult contains the simulation result for a single transaction group
type TxnGroupResult struct {
	Txns           []TxnResult
	FailureMessage string

	// FailedAt is the path to the txn that failed inside of this group
	FailedAt TxnPath

	// AppBudgetAdded is the total opcode budget for this group
	AppBudgetAdded uint64

	// AppBudgetConsumed is the total opcode cost used for this group
	AppBudgetConsumed uint64

	// FeeCredit is the fees left over after covering fees for this group
	FeeCredit uint64
}

func makeTxnGroupResult(txgroup []transactions.SignedTxn) TxnGroupResult {
	groupResult := TxnGroupResult{Txns: make([]TxnResult, len(txgroup))}
	for i, tx := range txgroup {
		groupResult.Txns[i] = TxnResult{Txn: transactions.SignedTxnWithAD{
			SignedTxn: tx,
		}}
	}
	return groupResult
}

// ResultLatestVersion is the latest version of the Result struct
const ResultLatestVersion = uint64(1)

// Result contains the result from a call to Simulator.Simulate
type Result struct {
	Version      uint64
	LastRound    basics.Round
	TxnGroups    []TxnGroupResult // this is a list so that supporting multiple in the future is not breaking
	WouldSucceed bool             // true iff no failure message, no missing signatures, and the budget was not exceeded
	Block        *ledgercore.ValidatedBlock
}

func makeSimulationResultWithVersion(lastRound basics.Round, txgroups [][]transactions.SignedTxn, version uint64) (Result, error) {
	if version != ResultLatestVersion {
		return Result{}, fmt.Errorf("invalid SimulationResult version: %d", version)
	}

	groups := make([]TxnGroupResult, len(txgroups))

	for i, txgroup := range txgroups {
		groups[i] = makeTxnGroupResult(txgroup)
	}

	return Result{
		Version:      version,
		LastRound:    lastRound,
		TxnGroups:    groups,
		WouldSucceed: true,
	}, nil
}

func makeSimulationResult(lastRound basics.Round, txgroups [][]transactions.SignedTxn) Result {
	result, err := makeSimulationResultWithVersion(lastRound, txgroups, ResultLatestVersion)
	if err != nil {
		// this should never happen, since we pass in ResultLatestVersion
		panic(err)
	}
	return result
}

// OpcodeTrace contains the trace effects of a single opcode evaluation
type OpcodeTrace struct {
	// The PC of the opcode being evaluated
	PC uint64

	// TODO: additional effects, like stack and scratch space changes
	// Ideally these effects contain only the "delta" information, unlike the
	// current dryrun trace, which returns the _entire_ stack and scratch space
	// for _every_ opcode.
	//
	// For stack, an interface that returns the first index that's changed and the rest of the stack
	// seems like a good idea. E.g. if the stack changes from [A, B, C] to [A, X, Y], the delta effect
	// could be something like { index: 1, elements: [X, Y] }. You can also think of it in terms of
	// pushes and pops if you want, so { pops: 2, pushes: [X, Y] }; the only difference is that index
	// is the absolute stack index while pops is a relative index from the end.
	//
	// * Note, this is not very efficient for opcodes like frame_bury and cover, which insert a single
	//   value potentially very deep into the stack, so perhaps an additional special case could be
	//   created to address that, maybe like { bury: X, depth: Y }? (I view this as nice to have but
	//   not critical for our first implementation).
	//
	// For scratch space, the good news is we don't have any opcodes which mutate more than 1 slot
	// at a time, and you can't ever delete a scratch value. So I think a more basic scheme is appropriate;
	// something that includes the index being set and the value should be sufficient.
	//
	// A note about types: stack and scratch space is "untyped", meaning either uint64s or byteslices
	// can be anywhere. Since JSON and msgpack arrays can also be untyped, I think it's a good idea
	// to embrace this and encode either a number or a byte string for each value.
}

// TransactionTrace contains the trace effects of a single transaction evaluation (including its inners)
type TransactionTrace struct {
	// LogicSigTrace contains the trace for LogicSig evaluation, if a LogicSig signed this transaction
	LogicSigTrace []OpcodeTrace
	// AppTrace contains the trace for an app evaluation, if this is an app call transaction
	AppTrace []OpcodeTrace
	// InnerTraces contains the traces for inner transactions, if this transaction spawned any. This
	// object only contains traces for inners that are immediate children of this transaction.
	// Grandchild traces will be present inside the TransactionTrace of their parent.
	InnerTraces []TransactionTrace
}

// Additional notes:
// * OpcodeTrace.PC is not very useful without the actual program bytecode. _For now_, I think that's
//   ok. In the future I think we will include the "initial state" of the simulation (e.g. every account,
//   app, local state, asset, holding, and box that's touched), and that might the problem for us.
//
