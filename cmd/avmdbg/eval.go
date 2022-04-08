package main

import (
	"fmt"
	"strings"

	"github.com/algorand/go-algorand/data/avmdbg"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/data/transactions/logic"
	"github.com/algorand/go-algorand/ledger"
)

// EvalLocation represents the relative location of a transaction or TEAL op being evaluated
type EvalLocationRelative struct {
	// IsTransaction indicates if this is a transaction index or a TEAL op index
	IsTransaction  bool
	TxnIndex       int
	TxnGroup       []transactions.SignedTxn
	TEALDebugState *logic.DebugState
}

func (elr EvalLocationRelative) String() string {
	var template string
	var index int
	if elr.IsTransaction {
		template = "txn %d"
		index = elr.TxnIndex
	} else {
		template = "TEAL pc %d"
		index = elr.TEALDebugState.PC
	}
	return fmt.Sprintf(template, index)
}

type EvalLocationAbsolute []EvalLocationRelative

func (ela *EvalLocationAbsolute) PushTransactionIndex(index int, txnGroup []transactions.SignedTxn) {
	hop := EvalLocationRelative{
		IsTransaction: true,
		TxnIndex:      index,
		TxnGroup:      txnGroup,
	}
	*ela = append(*ela, hop)
}

func (ela *EvalLocationAbsolute) PushTealOpIndex(state *logic.DebugState) {
	hop := EvalLocationRelative{
		IsTransaction:  false,
		TEALDebugState: state,
	}
	*ela = append(*ela, hop)
}

func (ela *EvalLocationAbsolute) Pop() {
	*ela = (*ela)[:len(*ela)-1]
}

func (ela *EvalLocationAbsolute) Copy() EvalLocationAbsolute {
	newCopy := make(EvalLocationAbsolute, len(*ela))
	copy(newCopy, *ela)
	return newCopy
}

func (ela *EvalLocationAbsolute) IsTransaction() bool {
	lastIndex := len(*ela) - 1
	return (*ela)[lastIndex].IsTransaction
}

func (ela *EvalLocationAbsolute) String() string {
	strs := make([]string, len(*ela))
	for i, hop := range *ela {
		strs[i] = hop.String()
	}
	return strings.Join(strs, "->")
}

// EvalLayer represents a "layer" being evaluated, which means a transaction group and a TEAL program
// it may be executing
type EvalLayer struct {
	TxnGroup      []transactions.SignedTxn
	TxnGroupIndex int
	TxnGroupADs   []transactions.ApplyData

	TEALState *logic.DebugState
}

type EvalLayers []EvalLayer

func (el EvalLayers) ToLocation() EvalLocationAbsolute {
	location := EvalLocationAbsolute{}

	for _, layer := range el {
		location.PushTransactionIndex(layer.TxnGroupIndex, layer.TxnGroup)
		if layer.TEALState != nil {
			location.PushTealOpIndex(layer.TEALState)
		}
	}

	return location
}

type EvalDebugger struct {
	ctx      *avmdbg.DebuggerContext
	reporter EventReporter

	layers           EvalLayers
	previousLocation EvalLocationAbsolute

	// TODO: breakpoints, state for step over, step into, step out

	TxIDBreakpoints []transactions.Txid
}

type EvalDebuggerAction int

const (
	EvalDebuggerActionContinue EvalDebuggerAction = iota
	EvalDebuggerActionStepOver
	EvalDebuggerActionStepInto
	EvalDebuggerActionStepOut
)

type EvalDebuggerCursor chan EvalDebuggerAction

func (edc EvalDebuggerCursor) Continue() {
	edc <- EvalDebuggerActionContinue
}

func (edc EvalDebuggerCursor) StepOver() {
	edc <- EvalDebuggerActionStepOver
}

func (edc EvalDebuggerCursor) StepInto() {
	edc <- EvalDebuggerActionStepInto
}

func (edc EvalDebuggerCursor) StepOut() {
	edc <- EvalDebuggerActionStepOut
}

func (edc EvalDebuggerCursor) awaitEvalAction() EvalDebuggerAction {
	return <-edc
}

type logicDebuggerHelper struct {
	ed *EvalDebugger
}

func (ldh logicDebuggerHelper) Register(state *logic.DebugState) error {
	ldh.ed.layers[len(ldh.ed.layers)-1].TEALState = state
	return nil
}

func (ldh logicDebuggerHelper) Update(state *logic.DebugState) error {
	// called before every TEAL op in the program
	ldh.ed.layers[len(ldh.ed.layers)-1].TEALState = state

	location := ldh.ed.layers.ToLocation()

	shouldWait := true // TODO
	if shouldWait {
		ldh.ed.pauseAndWait(location)
	}

	ldh.ed.previousLocation = location
	return nil
}

func (ldh logicDebuggerHelper) Complete(state *logic.DebugState) error {
	ldh.ed.layers[len(ldh.ed.layers)-1].TEALState = nil
	return nil
}

func (ldh logicDebuggerHelper) EnterInners(ep *logic.EvalParams) error {
	txnGroup := make([]transactions.SignedTxn, len(ep.TxnGroup))
	for i, txn := range ep.TxnGroup {
		txnGroup[i] = txn.SignedTxn
	}
	ldh.ed.layers = append(ldh.ed.layers, EvalLayer{
		TxnGroup: txnGroup,
	})
	return nil
}

func (ldh logicDebuggerHelper) InnerTxn(groupIndex int, ep *logic.EvalParams) error {
	currentLayerIndex := len(ldh.ed.layers) - 1
	ldh.ed.layers[currentLayerIndex].TxnGroupIndex = groupIndex
	if groupIndex > 0 {
		prevTxnIndex := groupIndex - 1
		ldh.ed.layers[currentLayerIndex].TxnGroupADs = append(ldh.ed.layers[currentLayerIndex].TxnGroupADs, ep.TxnGroup[prevTxnIndex].ApplyData)
	}

	location := ldh.ed.layers.ToLocation()
	txn := ldh.ed.layers[currentLayerIndex].TxnGroup[groupIndex]

	if ldh.ed.shouldPauseForTxn(location, txn) {
		ldh.ed.pauseAndWait(location)
	}

	ldh.ed.previousLocation = location
	return nil
}

func (ldh logicDebuggerHelper) LeaveInners(ep *logic.EvalParams) error {
	currentLayerIndex := len(ldh.ed.layers) - 1
	ldh.ed.layers[currentLayerIndex].TxnGroupIndex = len(ldh.ed.layers[currentLayerIndex].TxnGroup)
	location := ldh.ed.layers.ToLocation()

	ldh.ed.layers = ldh.ed.layers[:currentLayerIndex]

	shouldWait := true // TODO
	if shouldWait {
		ldh.ed.pauseAndWait(location)
	}

	return nil
}

func (ed *EvalDebugger) GetLogicDebugger() logic.DebuggerHook {
	return logicDebuggerHelper{ed: ed}
}

func (ed *EvalDebugger) AboutToEvalTransaction(groupIndex int, prevTxibs []transactions.SignedTxnInBlock) {
	currentLayerIndex := len(ed.layers) - 1
	ed.layers[currentLayerIndex].TxnGroupIndex = groupIndex
	if groupIndex > 0 {
		prevTxnIndex := groupIndex - 1
		ed.layers[currentLayerIndex].TxnGroupADs = append(ed.layers[currentLayerIndex].TxnGroupADs, prevTxibs[prevTxnIndex].ApplyData)
	}

	location := ed.layers.ToLocation()
	txn := ed.layers[currentLayerIndex].TxnGroup[groupIndex]

	if ed.shouldPauseForTxn(location, txn) {
		ed.pauseAndWait(location)
	}

	ed.previousLocation = location
}

func (ed *EvalDebugger) shouldPauseForTxn(location EvalLocationAbsolute, txn transactions.SignedTxn) bool {
	// txid := txn.ID()
	// for _, bps := range ed.TxIDBreakpoints {
	// 	if txid == bps {
	// 		return true
	// 	}
	// }

	// return false
	return true
}

func (ed *EvalDebugger) pauseAndWait(location EvalLocationAbsolute) {
	cursor := make(EvalDebuggerCursor)

	ed.reporter.reportEvalEvent(EvalEventFields{
		Location: location.Copy(),
		Cursor:   cursor,
	})

	// TODO: use return value to inform when to break next
	cursor.awaitEvalAction()
}

func StartEvaluator(ctx *avmdbg.DebuggerContext, evalDebugger *EvalDebugger, reporter EventReporter) {
	sdelta, stibs, err := ledger.EvalForDebugger(ctx, evalDebugger, ctx.Params.InputTxns)
	if err != nil {
		reporter.reportEvalEvent(EvalEventFields{
			Location: evalDebugger.previousLocation.Copy(),
			Error:    err,
		})
		return
	}

	currentLayerIndex := len(evalDebugger.layers) - 1
	evalDebugger.layers[currentLayerIndex].TxnGroupIndex = len(evalDebugger.layers[currentLayerIndex].TxnGroup)

	reporter.reportEvalEvent(EvalEventFields{
		Location:                    evalDebugger.layers.ToLocation(),
		Succeeded:                   true,
		SuccessfulStateDelta:        sdelta,
		SuccessfulSignedTxnsInBlock: stibs,
	})
}
