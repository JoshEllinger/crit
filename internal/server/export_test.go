package server

import (
	"path/filepath"
	"testing"

	"github.com/JoshEllinger/crit/internal/config"
	"github.com/JoshEllinger/crit/internal/review"
	"github.com/JoshEllinger/crit/internal/session"
)

var (
	resolveHost             = config.ResolveHost
	randomUUID              = session.RandomUUID
	maxAttachmentBytes      = session.MaxAttachmentBytes
	parseDaemonFlagsForTest = parseDaemonFlags
)

func writeCritJSONForTest(t *testing.T, dir string, cj CritJSON) {
	t.Helper()
	critPath := filepath.Join(dir, ".crit")
	if err := review.SaveCritJSON(critPath, cj); err != nil {
		t.Fatal(err)
	}
}
