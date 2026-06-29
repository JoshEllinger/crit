package auth

import "github.com/JoshEllinger/crit/internal/config"

type Config = config.Config

var (
	runAuth                = RunAuth
	lazyBackfillAuthUserID = LazyBackfillAuthUserID
)
