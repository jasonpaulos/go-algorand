package main

import (
	"fmt"

	"github.com/algorand/go-algorand/data/avmdbg"
	tm "github.com/buger/goterm"
)

type DebuggerStatus int

const (
	DebuggerStatusLoadingStart DebuggerStatus = iota
	DebuggerStatusLoadingFinished
	DebuggerStatusRunning
	DebuggerStatusPaused
	DebuggerStatusError
	DebuggerStatusFinished
	DebuggerStatusExit
)

type DebuggerState struct {
	Status DebuggerStatus

	FinishedLoading bool

	EvalDebugger   EvalDebugger
	PausedLocation EvalLocationAbsolute
	PausedCursor   EvalDebuggerCursor

	EvalError error

	LastCommand         []string
	LastCommandResponse string
}

func mainLoop(ctx *avmdbg.DebuggerContext) error {
	eventManager := make(EventManager)

	go HandleUserInput(eventManager.getReporter())
	go LoadState(ctx, eventManager.getReporter())

	var state DebuggerState
	state.EvalDebugger = EvalDebugger{
		ctx:      ctx,
		reporter: eventManager.getReporter(),
		layers: EvalLayers{
			EvalLayer{
				TxnGroup: ctx.Params.InputTxns,
			},
		},
	}

	for state.Status != DebuggerStatusExit {
		event := <-eventManager

		switch event.Type {
		case LoadingEventType:
			if event.LoadingEventFields.Finished {
				state.Status = DebuggerStatusLoadingFinished
				state.PausedLocation = EvalLocationAbsolute{
					EvalLocationRelative{
						IsTransaction: true,
						TxnIndex:      -1,
						TxnGroup:      ctx.Params.InputTxns,
					},
				}
			} else {
				state.Status = DebuggerStatusLoadingStart
			}
		case UserInputEventType:
			// execute a command

			// clear last command response
			state.LastCommandResponse = ""

			if len(event.UserInputCommand) == 0 {
				if len(state.LastCommand) != 0 {
					// replay last command
					event.UserInputCommand = state.LastCommand
				} else {
					// ignore
					break
				}
			}

			foundCommand := false
			for _, command := range allDebuggerCommands {
				if command.FirstWord == event.UserInputCommand[0] {
					command.Action(&state, ctx, eventManager, event.UserInputCommand[1:])
					foundCommand = true
					break
				}
			}

			if !foundCommand {
				state.LastCommandResponse = "Unknown command: " + event.UserInputCommand[0]
			}

			state.LastCommand = event.UserInputCommand
		case EvalEventType:
			if event.ErrorEventFields.Error != nil {
				state.Status = DebuggerStatusError
				state.PausedLocation = event.Location
				state.EvalError = event.ErrorEventFields.Error
				state.PausedCursor = nil
			} else if event.EvalEventFields.Succeeded {
				state.Status = DebuggerStatusFinished
				state.PausedLocation = event.Location
				state.PausedCursor = nil
			} else {
				state.Status = DebuggerStatusPaused
				state.PausedLocation = event.Location
				state.PausedCursor = event.Cursor
			}
		case ErrorEventType:
			return event.ErrorEventFields.Error
		default:
			return fmt.Errorf("Unknown event type: %v", event.Type)
		}

		render(state)
	}
	tm.Clear()

	// TODO: check that 0 < len(txnx) <= 16
	// TODO: check that txns are well formed
	// TODO: check that txns have a valid group
	// TODO: check txns signatures
	// verify.TxnGroup(params.InputTxns, )

	return nil
}
