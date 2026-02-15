package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"syscall"
)

type Server struct {
	doc    *Document
	mux    *http.ServeMux
	assets fs.FS
}

func NewServer(doc *Document, frontendFS embed.FS) *Server {
	s := &Server{doc: doc}

	assets, _ := fs.Sub(frontendFS, "frontend")
	s.assets = assets

	mux := http.NewServeMux()
	mux.HandleFunc("/api/document", s.handleDocument)
	mux.HandleFunc("/api/comments", s.handleComments)
	mux.HandleFunc("/api/comments/", s.handleCommentByID)
	mux.HandleFunc("/api/finish", s.handleFinish)
	mux.HandleFunc("/api/stale", s.handleStale)
	mux.Handle("/", http.FileServer(http.FS(assets)))

	s.mux = mux
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := map[string]string{
		"filename": s.doc.FileName,
		"content":  s.doc.Content,
	}
	writeJSON(w, resp)
}

func (s *Server) handleStale(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		notice := s.doc.GetStaleNotice()
		writeJSON(w, map[string]string{"notice": notice})
	case http.MethodDelete:
		s.doc.ClearStaleNotice()
		writeJSON(w, map[string]string{"status": "ok"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleComments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		comments := s.doc.GetComments()
		writeJSON(w, comments)

	case http.MethodPost:
		var req struct {
			StartLine int    `json:"start_line"`
			EndLine   int    `json:"end_line"`
			Body      string `json:"body"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Body == "" {
			http.Error(w, "Comment body is required", http.StatusBadRequest)
			return
		}
		if req.StartLine < 1 || req.EndLine < req.StartLine {
			http.Error(w, "Invalid line range", http.StatusBadRequest)
			return
		}

		c := s.doc.AddComment(req.StartLine, req.EndLine, req.Body)
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, c)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleCommentByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	if id == "" {
		http.Error(w, "Comment ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		var req struct {
			Body string `json:"body"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Body == "" {
			http.Error(w, "Comment body is required", http.StatusBadRequest)
			return
		}
		c, ok := s.doc.UpdateComment(id, req.Body)
		if !ok {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		writeJSON(w, c)

	case http.MethodDelete:
		if !s.doc.DeleteComment(id) {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		writeJSON(w, map[string]string{"status": "deleted"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFinish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.doc.WriteFiles()

	writeJSON(w, map[string]string{
		"status":      "finished",
		"review_file": s.doc.reviewFilePath(),
	})

	go func() {
		fmt.Println("\nFinish review requested. Shutting down...")
		// Give time for the response to be sent
		<-r.Context().Done()
		// Use process signal to trigger graceful shutdown
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	}()
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
