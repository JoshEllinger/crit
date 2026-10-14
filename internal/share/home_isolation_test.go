package share

import (
	"os"
	"testing"

	"github.com/JoshEllinger/crit/internal/testutil"
)

func TestMain(m *testing.M) { os.Exit(testutil.RunWithTempHome(m)) }
