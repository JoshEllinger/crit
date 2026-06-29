package main

import (
	"github.com/JoshEllinger/crit/internal/auth"
	"github.com/JoshEllinger/crit/internal/clicmd"
	"github.com/JoshEllinger/crit/internal/comment"
	"github.com/JoshEllinger/crit/internal/config"
	"github.com/JoshEllinger/crit/internal/github"
	"github.com/JoshEllinger/crit/internal/live"
	"github.com/JoshEllinger/crit/internal/preview"
	"github.com/JoshEllinger/crit/internal/review"
	"github.com/JoshEllinger/crit/internal/session"
	"github.com/JoshEllinger/crit/internal/share"
)

func runShare(args []string)     { clicmd.Exit(share.RunShare(args)) }
func runFetch(args []string)     { clicmd.Exit(share.RunFetch(args)) }
func runUnpublish(args []string) { clicmd.Exit(share.RunUnpublish(args)) }
func runConfig(args []string)    { clicmd.Exit(config.RunConfig(args)) }
func runPull(args []string)      { clicmd.Exit(github.RunPull(args)) }
func runPush(args []string)      { clicmd.Exit(github.RunPush(args)) }
func runComment(args []string)   { clicmd.Exit(comment.RunComment(args)) }
func runComments(args []string)  { clicmd.Exit(comment.RunComments(args)) }
func runReview(args []string)    { clicmd.Exit(session.RunReview(args)) }
func runPlan(args []string)      { clicmd.Exit(session.RunPlan(args)) }
func runStop(args []string)      { clicmd.Exit(session.RunStop(args)) }
func runStatus(args []string)    { clicmd.Exit(session.RunStatus(args)) }
func runCleanup(args []string)   { clicmd.Exit(review.RunCleanup(args)) }
func runPR(args []string)        { clicmd.Exit(github.RunPR(args)) }

func runLive(args []string)    { live.RunLive(args) }
func runPreview(args []string) { preview.RunPreview(args) }
func runAuth(args []string)    { auth.RunAuth(args) }
func runStats(args []string)   { session.RunStats(args) }

func runPlanHook()      { clicmd.Exit(session.RunPlanHook()) }
func runCodexPlanHook() { clicmd.Exit(session.RunCodexPlanHook()) }
