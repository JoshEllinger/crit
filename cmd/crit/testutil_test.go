package main

import (
	"github.com/JoshEllinger/crit/internal/vcs"
)

func resetDefaultBranchOnce() {
	vcs.ResetDefaultBranchOnceForTest()
}
