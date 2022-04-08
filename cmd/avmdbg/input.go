package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func HandleUserInput(reporter EventReporter) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))
		parts := strings.Fields(line)
		reporter.reportUserInput(parts)
	}

	if err := scanner.Err(); err != nil {
		reporter.reportError(fmt.Errorf("Error reading standard input: %w", err))
		return
	}

	reporter.reportError(errors.New("Standard input ended"))
}

// type TransparentScanner struct {
// 	input         io.Reader
// 	pendingBuffer bytes.Buffer
// 	err           error
// }

// func MakeTransparentScanner(input io.Reader) *TransparentScanner {
// 	return &TransparentScanner{
// 		input: input,
// 	}
// }

// func (ts *TransparentScanner) Scan() bool {
// 	readBuffer := [32]byte{}
// 	for {
// 		n, err := io.ReadAtLeast(ts.input, readBuffer[:], 1)
// 		if n < 1 {
// 			ts.err = err
// 			return false
// 		}

// 		if bytes.ContainsRune(readBuffer[:], '\n') {
// 			pending := ts.pendingBuffer.Bytes()
// 		}

// 		ts.pendingBuffer.Write(readBuffer[:n])
// 	}
// 	// return ts.scanner.Scan()
// }

// func (ts *TransparentScanner) Err() error {
// 	return ts.err
// }

// func (ts *TransparentScanner) Text() string {
// 	return ts.scanner.Text()
// }

// func (ts *TransparentScanner) PendingText() string {
// 	return
// }
