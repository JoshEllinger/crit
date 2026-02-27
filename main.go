package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

//go:embed frontend/*
var frontendFS embed.FS

//go:embed integrations/*
var integrationsFS embed.FS

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Handle "crit help" subcommand
	if len(os.Args) >= 2 && (os.Args[1] == "help" || os.Args[1] == "--help" || os.Args[1] == "-h") {
		printHelp()
		os.Exit(0)
	}

	// Handle "crit go [port]" subcommand — signals round-complete to a running crit server
	if len(os.Args) >= 2 && os.Args[1] == "go" {
		port := "3000" // default
		if len(os.Args) >= 3 {
			port = os.Args[2]
		}
		resp, err := http.Post("http://localhost:"+port+"/api/round-complete", "application/json", nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not reach crit on port %s: %v\n", port, err)
			os.Exit(1)
		}
		resp.Body.Close()
		if resp.StatusCode == 200 {
			fmt.Println("Round complete — crit will reload.")
		} else {
			fmt.Fprintf(os.Stderr, "Unexpected status: %d\n", resp.StatusCode)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Handle "crit install [agent]" subcommand
	if len(os.Args) >= 2 && os.Args[1] == "install" {
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: crit install <agent>")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Available agents:")
			for _, a := range availableIntegrations() {
				fmt.Fprintf(os.Stderr, "  %s\n", a)
			}
			fmt.Fprintln(os.Stderr, "  all")
			os.Exit(1)
		}
		target := os.Args[2]
		if target == "all" {
			for _, arg := range os.Args[3:] {
				if arg == "--global" || arg == "-g" {
					fmt.Fprintln(os.Stderr, "Error: --global is not supported with 'all'")
					os.Exit(1)
				}
			}
			for _, name := range availableIntegrations() {
				installIntegration(name)
			}
		} else {
			installIntegration(target)
		}
		os.Exit(0)
	}

	port := flag.Int("port", 0, "Port to listen on (default: random available port)")
	flag.IntVar(port, "p", 0, "Port to listen on (shorthand)")
	outputDir := flag.String("output", "", "Output directory for review files (default: same dir as input file)")
	flag.StringVar(outputDir, "o", "", "Output directory (shorthand)")
	noOpen := flag.Bool("no-open", false, "Don't auto-open browser")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.BoolVar(showVersion, "v", false, "Print version and exit (shorthand)")
	shareURL := flag.String("share-url", "", "Base URL of hosted Crit service for sharing reviews (overrides CRIT_SHARE_URL env var)")
	flag.Usage = func() {
		printHelp()
	}
	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	filePath := flag.Arg(0)
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalf("Error resolving path: %v", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if info.IsDir() {
		log.Fatalf("Error: %s is a directory, not a file", absPath)
	}

	outDir := *outputDir
	if outDir == "" {
		outDir = filepath.Dir(absPath)
	}

	doc, err := NewDocument(absPath, outDir)
	if err != nil {
		log.Fatalf("Error loading document: %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	addr := listener.Addr().(*net.TCPAddr)

	if *shareURL == "" {
		*shareURL = os.Getenv("CRIT_SHARE_URL")
	}
	srv, err := NewServer(doc, frontendFS, *shareURL, version, addr.Port)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	if os.Getenv("CRIT_NO_UPDATE_CHECK") == "" {
		go srv.checkForUpdates()
	}
	httpServer := &http.Server{
		Handler:     srv,
		ReadTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
		// No WriteTimeout — SSE connections need to stay open
	}

	status := newStatus(os.Stdout)
	srv.status = status
	doc.status = status

	url := fmt.Sprintf("http://localhost:%d", addr.Port)
	status.Listening(url)

	if !*noOpen {
		go openBrowser(url)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	watchStop := make(chan struct{})
	go doc.WatchFile(watchStop)

	go func() {
		if err := httpServer.Serve(listener); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	close(watchStop)
	fmt.Println()

	doc.Shutdown()
	doc.WriteFiles()

	shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(shutCtx)
}

func printHelp() {
	fmt.Fprintf(os.Stderr, `crit — inline code review for AI agent workflows

Usage:
  crit <file.md>              Open a file for review in your browser
  crit go [port]              Signal round-complete to a running crit instance
  crit install <agent>        Install integration files for an AI coding tool
  crit help                   Show this help message

Agents:
  claude-code, cursor, windsurf, github-copilot, cline, all

Options:
  -p, --port <port>           Port to listen on (default: random)
  -o, --output <dir>          Output directory for review files
      --no-open               Don't auto-open browser
      --share-url <url>       Share service URL (no default)
  -v, --version               Print version

Environment:
  CRIT_SHARE_URL              Override the share service URL
  CRIT_NO_UPDATE_CHECK        Disable update check on startup

`)
}

func printVersion() {
	line := "crit " + version
	var details []string
	if date != "unknown" {
		details = append(details, date)
	}
	if commit != "unknown" {
		short := commit
		if len(short) > 7 {
			short = short[:7]
		}
		details = append(details, short)
	}
	if len(details) > 0 {
		line += " (" + strings.Join(details, ", ") + ")"
	}
	fmt.Println(line)
	fmt.Println("Inline code review for AI agent workflows")
}

type integration struct {
	source string // path inside integrations/ embed
	dest   string // destination relative to cwd
	hint   string // usage hint printed after install
}

var integrationMap = map[string][]integration{
	"claude-code": {
		{source: "integrations/claude-code/crit.md", dest: ".claude/commands/crit.md", hint: "Run /crit in Claude Code to start a review loop"},
	},
	"cursor": {
		{source: "integrations/cursor/crit-command.md", dest: ".cursor/commands/crit.md", hint: "Run /crit in Cursor to start a review loop"},
	},
	"windsurf": {
		{source: "integrations/windsurf/crit.md", dest: ".windsurf/rules/crit.md", hint: "Windsurf will suggest Crit when writing plans"},
	},
	"github-copilot": {
		{source: "integrations/github-copilot/crit.prompt.md", dest: ".github/prompts/crit.prompt.md", hint: "Run /crit in GitHub Copilot to start a review loop"},
	},
	"cline": {
		{source: "integrations/cline/crit.md", dest: ".clinerules/crit.md", hint: "Cline will suggest Crit when writing plans"},
	},
}

func availableIntegrations() []string {
	return []string{"claude-code", "cursor", "windsurf", "github-copilot", "cline"}
}

func installIntegration(name string) {
	files, ok := integrationMap[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown agent: %s\n\nAvailable agents:\n", name)
		for _, a := range availableIntegrations() {
			fmt.Fprintf(os.Stderr, "  %s\n", a)
		}
		os.Exit(1)
	}

	force := false
	global := false
	for _, arg := range os.Args[3:] {
		if arg == "--force" || arg == "-f" {
			force = true
		}
		if arg == "--global" || arg == "-g" {
			global = true
		}
	}

	var homeDir string
	if global {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
	}

	for _, f := range files {
		dest := f.dest
		if global {
			dest = filepath.Join(homeDir, f.dest)
		}

		if !force {
			if _, err := os.Stat(dest); err == nil {
				fmt.Printf("  Skipped:   %s (already exists, use --force to overwrite)\n", dest)
				continue
			}
		}

		data, err := integrationsFS.ReadFile(f.source)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading embedded file %s: %v\n", f.source, err)
			os.Exit(1)
		}

		dir := filepath.Dir(dest)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}

		if err := os.WriteFile(dest, data, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", dest, err)
			os.Exit(1)
		}

		fmt.Printf("  Installed: %s\n", dest)
	}
	if files[0].hint != "" {
		fmt.Printf("  %s\n", files[0].hint)
	}
	fmt.Println()
}

func openBrowser(url string) {
	time.Sleep(200 * time.Millisecond)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}
	_ = cmd.Run()
}
