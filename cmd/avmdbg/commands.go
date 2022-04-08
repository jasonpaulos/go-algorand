package main

import (
	"fmt"
	"strings"

	"github.com/algorand/go-algorand/data/avmdbg"
	"github.com/algorand/go-algorand/data/transactions"
)

type Command struct {
	FirstWord string
	Help      string
	Action    func(state *DebuggerState, ctx *avmdbg.DebuggerContext, eventManager EventManager, args []string)
}

var exitCommand = Command{
	FirstWord: "exit",
	Help:      `Exists the program`,
	Action: func(state *DebuggerState, ctx *avmdbg.DebuggerContext, eventManager EventManager, args []string) {
		if len(args) > 0 {
			state.LastCommandResponse = "'exit' expects no arguments"
			return
		}
		state.Status = DebuggerStatusExit
	},
}

var runCommand = Command{
	FirstWord: "run",
	Help:      `Begin executing the transaction group`,
	Action: func(state *DebuggerState, ctx *avmdbg.DebuggerContext, eventManager EventManager, args []string) {
		if len(args) > 0 {
			state.LastCommandResponse = "'run' expects no arguments"
			return
		}

		if state.Status == DebuggerStatusLoadingStart {
			state.LastCommandResponse = "Cannot process 'run' command before state is loaded"
			return
		}

		if state.Status != DebuggerStatusLoadingFinished {
			state.LastCommandResponse = "Cannot process 'run' command after run has started"
			return
		}

		state.Status = DebuggerStatusRunning
		go StartEvaluator(ctx, &state.EvalDebugger, eventManager.getReporter())
	},
}

var stepCommand = Command{
	FirstWord: "step",
	Help:      `Advance the evaluator by a single step, moving into any child steps of the current step`,
	Action: func(state *DebuggerState, ctx *avmdbg.DebuggerContext, eventManager EventManager, args []string) {
		if len(args) > 0 {
			state.LastCommandResponse = "'step' expects no arguments"
			return
		}

		if state.Status != DebuggerStatusPaused {
			state.LastCommandResponse = "Cannot process 'step' command when debugger is not paused"
			return
		}

		state.Status = DebuggerStatusRunning
		state.PausedCursor.StepInto()
	},
}

var nextCommand = Command{
	FirstWord: "next",
	Help:      `Advance the evaluator by a single step, skipping any child steps of the current step`,
	Action: func(state *DebuggerState, ctx *avmdbg.DebuggerContext, eventManager EventManager, args []string) {
		if len(args) > 0 {
			state.LastCommandResponse = "'next' expects no arguments"
			return
		}

		if state.Status != DebuggerStatusPaused {
			state.LastCommandResponse = "Cannot process 'next' command when debugger is not paused"
			return
		}

		state.Status = DebuggerStatusRunning
		state.PausedCursor.StepOver()
	},
}

var finishCommand = Command{
	FirstWord: "finish",
	Help:      `Move the evaluator to the end of the current unit in the evaluator`,
	Action: func(state *DebuggerState, ctx *avmdbg.DebuggerContext, eventManager EventManager, args []string) {
		if len(args) > 0 {
			state.LastCommandResponse = "'finish' expects no arguments"
			return
		}

		if state.Status != DebuggerStatusPaused {
			state.LastCommandResponse = "Cannot process 'finish' command when debugger is not paused"
			return
		}

		state.Status = DebuggerStatusRunning
		state.PausedCursor.StepOut()
	},
}

var continueCommand = Command{
	FirstWord: "continue",
	Help:      `Resume execution until the next breakpoint`,
	Action: func(state *DebuggerState, ctx *avmdbg.DebuggerContext, eventManager EventManager, args []string) {
		if len(args) > 0 {
			state.LastCommandResponse = "'continue' expects no arguments"
			return
		}

		if state.Status != DebuggerStatusPaused {
			state.LastCommandResponse = "Cannot process 'continue' command when debugger is not paused"
			return
		}

		state.Status = DebuggerStatusRunning
		state.PausedCursor.Continue()
	},
}

var breakpointCommand = Command{
	FirstWord: "breakpoint",
	Help:      `Toggle a breakpoint at the given location`,
	Action: func(state *DebuggerState, ctx *avmdbg.DebuggerContext, eventManager EventManager, args []string) {
		if state.Status == DebuggerStatusRunning {
			state.LastCommandResponse = "Cannot set a breakpoint when the evaluator is running"
			return
		}

		if state.Status == DebuggerStatusFinished || state.Status == DebuggerStatusError {
			state.LastCommandResponse = "Cannot set a breakpoint because the evaluator has terminated"
			return
		}

		if len(args) == 0 {
			state.LastCommandResponse = "'breakpoint' expects at least one argument"
			return
		}

		switch args[0] {
		case "transaction", "txn":
			if len(args) != 2 {
				state.LastCommandResponse = fmt.Sprintf("Proper usage: 'breakpoint %s <TxID>'", args[0])
				break
			}
			textTxID := strings.ToUpper(args[1])
			var txid transactions.Txid
			err := txid.UnmarshalText([]byte(textTxID))
			if err != nil {
				state.LastCommandResponse = fmt.Sprintf("Could not parse TxID '%s': %s", textTxID, err.Error())
				break
			}
			state.EvalDebugger.TxIDBreakpoints = append(state.EvalDebugger.TxIDBreakpoints, txid)
			state.LastCommandResponse = fmt.Sprintf("Setting a breakpoint for txn %s", txid)
		default:
			state.LastCommandResponse = fmt.Sprintf("Unknown argument for breakpoint: '%s'", args[0])
		}
	},
}

var allDebuggerCommands = [...]Command{
	exitCommand,
	runCommand,
	stepCommand,
	nextCommand,
	finishCommand,
	continueCommand,
	breakpointCommand,
}
