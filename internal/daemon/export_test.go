package daemon

import "github.com/JoshEllinger/crit/internal/config"

type commonDaemonFlags = CommonDaemonFlags

var (
	atomicWriteFile         = config.AtomicWriteFile
	appendCommonDaemonFlags = AppendCommonDaemonFlags
)
