package share

import "github.com/JoshEllinger/crit/internal/testutil"

var (
	tokenFromHostedURL   = TokenFromHostedURL
	mustMkdirAll         = testutil.MustMkdirAll
	buildSharePayload    = BuildSharePayload
	loadCommentsForShare = LoadCommentsForShare
	unpublishFromWeb     = UnpublishFromWeb
	shareScope           = ShareScope
)
