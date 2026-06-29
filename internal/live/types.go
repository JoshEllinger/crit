package live

import (
	"github.com/JoshEllinger/crit/internal/comment"
	"github.com/JoshEllinger/crit/internal/config"
	"github.com/JoshEllinger/crit/internal/github"
	"github.com/JoshEllinger/crit/internal/review"
	"github.com/JoshEllinger/crit/internal/server"
	"github.com/JoshEllinger/crit/internal/session"
	"github.com/JoshEllinger/crit/internal/share"
	"github.com/JoshEllinger/crit/internal/testutil"
)

var writeFile = testutil.WriteFile

type (
	Config       = config.Config
	CritJSON     = session.CritJSON
	CritJSONFile = session.CritJSONFile
	Session      = session.Session
	FileEntry    = session.FileEntry
	Comment      = session.Comment
	DOMAnchor    = session.DOMAnchor
	SSEEvent     = session.SSEEvent
)

var (
	looksLikeLiveArgs      = LooksLikeLiveArgs
	saveCritJSON           = review.SaveCritJSON
	loadCritJSON           = review.LoadCritJSON
	appendReply            = comment.AppendReply
	checkShareAllowed      = share.CheckShareAllowed
	checkGitHubSyncAllowed = share.CheckGitHubSyncAllowed
	checkCommentCLIAllowed = comment.CheckCommentCLIAllowed
	carryForwardComment    = session.CarryForwardComment
	NewServer              = server.NewServer
	frontendFS             = server.FrontendFS
)

type GhComment = github.GhComment

var mergeGHComments = github.MergeGHComments
