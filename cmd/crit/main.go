package main

import (
	"os"

	integrationassets "github.com/JoshEllinger/crit/integrations"
	"github.com/JoshEllinger/crit/internal/clicmd"
	"github.com/JoshEllinger/crit/internal/session"
	webassets "github.com/JoshEllinger/crit/web"
)

var frontendFS = webassets.FS
var integrationsFS = integrationassets.FS

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		clicmd.Exit(session.RunReview(nil))
		return
	}
	if handled, err := dispatchCLI(os.Args[1:]); handled {
		clicmd.Exit(err)
		return
	}
	args := resolveAtPrefixedArgs(os.Args[1:])
	runPositionalCLI(args)
}
