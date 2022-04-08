package main

import (
	"fmt"
	"strings"

	"github.com/algorand/go-algorand/data/basics"
	tm "github.com/buger/goterm"
)

func render(state DebuggerState) {
	// var message string

	// switch state.Status {
	// case DebuggerStatusLoadingStart:
	// 	message = "Loading on-chain resources..."
	// case DebuggerStatusLoadingFinished:
	// 	message = "Resources loaded, type 'run' to start debugging"
	// case DebuggerStatusRunning:
	// 	message = "Program running"
	// case DebuggerStatusPaused:
	// 	message = fmt.Sprintf("Paused at %s", state.PausedLocation.String())
	// case DebuggerStatusError:
	// 	message = fmt.Sprintf("Encountered error at %s: %v", state.PausedLocation.String(), state.EvalError)
	// case DebuggerStatusFinished:
	// 	message = "Succeeded without errors. Want to examine the results?"
	// }

	// By moving cursor to top-left position we ensure that console output
	// will be overwritten each time, instead of adding new.
	tm.MoveCursor(1, 1)
	tm.Clear()

	// tm.Println(tm.Background(tm.Color(tm.Bold(message), tm.RED), tm.WHITE))

	if state.Status == DebuggerStatusLoadingStart {
		welcomeBox := tm.NewBox(60|tm.PCT, 60|tm.PCT, 0)

		fmt.Fprintf(welcomeBox, "Welcome to the AVM debugger\n\n")
		fmt.Fprintf(welcomeBox, "Loading resources from algod...\n")

		tm.Print(tm.MoveTo(welcomeBox.String(), 20|tm.PCT, 20|tm.PCT))
	} else {
		topBox := tm.NewBox(tm.Width(), tm.Height()-1, 0)

		// statusBox := tm.NewBox(tm.Width(), 3, 0)
		var debuggerStatus string
		switch state.Status {
		case DebuggerStatusLoadingFinished:
			debuggerStatus = "ready. Type 'run' to start debugging"
		case DebuggerStatusRunning:
			debuggerStatus = "running"
		case DebuggerStatusPaused:
			debuggerStatus = fmt.Sprintf("paused at %s", state.PausedLocation.String())
		case DebuggerStatusError:
			debuggerStatus = fmt.Sprintf("encountered error: %v", state.EvalError)
		case DebuggerStatusFinished:
			debuggerStatus = "succeeded without errors"
		default:
			debuggerStatus = "status unknown"
		}
		fmt.Fprintf(topBox, "Debugger %s\n\n", debuggerStatus)

		// tm.MoveCursor(1, 4)

		if len(state.PausedLocation) == 0 || state.PausedLocation.IsTransaction() {
			txnTable := tm.NewTable(0, 10, 5, ' ', 0)
			fmt.Fprintf(txnTable, "Index\tTxID\tTxType\tSender\tExecuting\n")

			// selectedTxn := -1
			// if len(state.PausedLocation) > 0 {
			// 	selectedTxn = state.PausedLocation[len(state.PausedLocation)-1].TxnIndex
			// }
			location := state.PausedLocation[len(state.PausedLocation)-1]

			for i, txn := range location.TxnGroup {
				txid := txn.Txn.ID()
				txtype := txn.Txn.Type
				sender := txn.Txn.Sender
				selected := ""
				if i == location.TxnIndex {
					selected = "<--"
				}
				fmt.Fprintf(txnTable, "%d\t%s\t%s\t%s\t%s\n", i, txid, txtype, sender, selected)
			}

			if location.TxnIndex == len(location.TxnGroup) {
				fmt.Fprintf(txnTable, "\t\t\t\t<--\n")
			}

			fmt.Fprintf(topBox, txnTable.String())
		} else {
			debugState := state.PausedLocation[len(state.PausedLocation)-1].TEALDebugState

			programTable := tm.NewTable(0, 10, 5, ' ', 0)
			fmt.Fprintf(programTable, "Line\tPC\tOp\tExecuting\tStack (Top ... Bottom)\n")

			approxTableHeight := tm.Height() - 3

			lines := strings.Split(strings.TrimSpace(debugState.Disassembly), "\n")

			linesStart := 0
			linesLength := approxTableHeight + 2 // add 2 more to be safe

			if debugState.Line > linesLength/2 {
				linesStart = debugState.Line - linesLength/2
			}

			linesEnd := linesStart + linesLength
			if linesEnd > len(lines) {
				linesEnd = len(lines)
			}

			for i, line := range lines[linesStart:linesEnd] {
				lineNumber := linesStart + i
				pc := debugState.LineToPC(lineNumber)
				safeLine := strings.ReplaceAll(line, "%", "%%")
				stack := ""
				selected := ""
				if debugState.Line == lineNumber {
					selected = "<--"

					stackAsStrs := make([]string, len(debugState.Stack))
					for i, value := range debugState.Stack {
						prefix := ""
						if value.Type == basics.TealBytesType {
							prefix = "0x"
						}
						mirroredIndex := len(stackAsStrs) - i - 1
						stackAsStrs[mirroredIndex] = prefix + value.String()
					}
					stack = fmt.Sprintf("[%s]", strings.Join(stackAsStrs, ","))
				}
				fmt.Fprintf(programTable, "%d\t%d\t%s\t%s\t%s\n", lineNumber, pc, safeLine, selected, stack)
			}

			fmt.Fprintf(topBox, programTable.String())
		}

		tm.Print(tm.MoveTo(topBox.String(), 0, 1))
	}

	// bottomBox := tm.NewBox(tm.Width(), 50|tm.PCT, 0)

	// fmt.Fprintf(bottomBox, "Hello from the bottom box")

	// tm.Print(tm.MoveTo(bottomBox.String(), 0, 50|tm.PCT))

	tm.MoveCursor(2, tm.Height()-1)
	if len(state.LastCommandResponse) != 0 {
		tm.Print(tm.Color(state.LastCommandResponse+". ", tm.RED))
	}
	// tm.Print("Available commands: TODO")

	tm.MoveCursor(1, tm.Height())
	tm.Print("(avmdbg) ")

	tm.Flush()
}
