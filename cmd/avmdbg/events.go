package main

import (
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/ledger/ledgercore"
)

type EventType int

const (
	ErrorEventType EventType = iota
	LoadingEventType
	UserInputEventType
	EvalEventType
)

type Event struct {
	Type EventType

	ErrorEventFields
	LoadingEventFields
	UserInputEventFields
	EvalEventFields
}

type ErrorEventFields struct {
	Error error
}

type LoadingEventFields struct {
	Finished bool
}

type UserInputEventFields struct {
	UserInputCommand []string
}

type EvalEventFields struct {
	Location                    EvalLocationAbsolute
	Cursor                      EvalDebuggerCursor
	Error                       error
	Succeeded                   bool
	SuccessfulStateDelta        ledgercore.StateDelta
	SuccessfulSignedTxnsInBlock []transactions.SignedTxnInBlock
}

type EventManager chan Event

func (em EventManager) getReporter() EventReporter {
	var c chan Event = em
	return c
}

type EventReporter chan<- Event

func (er EventReporter) reportError(err error) {
	er <- Event{
		Type: ErrorEventType,
		ErrorEventFields: ErrorEventFields{
			Error: err,
		},
	}
}

func (er EventReporter) reportLoadingEvent(finished bool) {
	er <- Event{
		Type: LoadingEventType,
		LoadingEventFields: LoadingEventFields{
			Finished: finished,
		},
	}
}

func (er EventReporter) reportUserInput(command []string) {
	er <- Event{
		Type: UserInputEventType,
		UserInputEventFields: UserInputEventFields{
			UserInputCommand: command,
		},
	}
}

func (er EventReporter) reportEvalEvent(fields EvalEventFields) {
	er <- Event{
		Type:            EvalEventType,
		EvalEventFields: fields,
	}
}
