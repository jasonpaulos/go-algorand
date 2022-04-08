package main

import (
	"fmt"

	"github.com/algorand/go-algorand/data/avmdbg"
)

func LoadState(ctx *avmdbg.DebuggerContext, reporter EventReporter) {
	reporter.reportLoadingEvent(false)

	var err error
	for attempt := 0; attempt < 5; attempt++ {
		// use multiple attempts in case REST responses are from different blocks
		err = ctx.GatherResources()
		if err == nil {
			break
		}
	}
	if err != nil {
		reporter.reportError(fmt.Errorf("Error gathering on-chain resources: %w", err))
		return
	}

	reporter.reportLoadingEvent(true)
}
