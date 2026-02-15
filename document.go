package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Comment struct {
	ID        string `json:"id"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CommentsFile struct {
	File      string    `json:"file"`
	FileHash  string    `json:"file_hash"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type Document struct {
	FilePath    string
	FileName    string
	FileDir     string
	Content     string
	FileHash    string
	OutputDir   string
	Comments    []Comment
	mu          sync.RWMutex
	nextID      int
	writeTimer  *time.Timer
	staleNotice string
}

func NewDocument(filePath, outputDir string) (*Document, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	content := string(data)
	hash := fmt.Sprintf("sha256:%x", sha256.Sum256(data))

	doc := &Document{
		FilePath:  filePath,
		FileName:  filepath.Base(filePath),
		FileDir:   filepath.Dir(filePath),
		Content:   content,
		FileHash:  hash,
		OutputDir: outputDir,
		Comments:  []Comment{},
		nextID:    1,
	}

	doc.loadComments()
	return doc, nil
}

func (d *Document) commentsFilePath() string {
	return filepath.Join(d.OutputDir, "."+d.FileName+".comments.json")
}

func (d *Document) reviewFilePath() string {
	ext := filepath.Ext(d.FileName)
	base := strings.TrimSuffix(d.FileName, ext)
	return filepath.Join(d.OutputDir, base+".review"+ext)
}

func (d *Document) loadComments() {
	data, err := os.ReadFile(d.commentsFilePath())
	if err != nil {
		return
	}

	var cf CommentsFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return
	}

	if cf.FileHash != d.FileHash {
		d.staleNotice = "The source file has changed since the last review session. Previous comments may not align with the current content."
		return
	}

	d.Comments = cf.Comments
	for _, c := range d.Comments {
		id := 0
		fmt.Sscanf(c.ID, "c%d", &id)
		if id >= d.nextID {
			d.nextID = id + 1
		}
	}
}

func (d *Document) AddComment(startLine, endLine int, body string) Comment {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	c := Comment{
		ID:        fmt.Sprintf("c%d", d.nextID),
		StartLine: startLine,
		EndLine:   endLine,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}
	d.nextID++
	d.Comments = append(d.Comments, c)
	d.scheduleWrite()
	return c
}

func (d *Document) UpdateComment(id, body string) (Comment, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, c := range d.Comments {
		if c.ID == id {
			d.Comments[i].Body = body
			d.Comments[i].UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			d.scheduleWrite()
			return d.Comments[i], true
		}
	}
	return Comment{}, false
}

func (d *Document) DeleteComment(id string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, c := range d.Comments {
		if c.ID == id {
			d.Comments = append(d.Comments[:i], d.Comments[i+1:]...)
			d.scheduleWrite()
			return true
		}
	}
	return false
}

func (d *Document) GetComments() []Comment {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]Comment, len(d.Comments))
	copy(result, d.Comments)
	return result
}

func (d *Document) GetStaleNotice() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.staleNotice
}

func (d *Document) ClearStaleNotice() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.staleNotice = ""
}

func (d *Document) scheduleWrite() {
	if d.writeTimer != nil {
		d.writeTimer.Stop()
	}
	d.writeTimer = time.AfterFunc(200*time.Millisecond, func() {
		d.WriteFiles()
	})
}

func (d *Document) WriteFiles() {
	d.mu.RLock()
	comments := make([]Comment, len(d.Comments))
	copy(comments, d.Comments)
	d.mu.RUnlock()

	d.writeCommentsJSON(comments)
	d.writeReviewMD(comments)
}

func (d *Document) writeCommentsJSON(comments []Comment) {
	cf := CommentsFile{
		File:      d.FileName,
		FileHash:  d.FileHash,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Comments:  comments,
	}

	data, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling comments: %v\n", err)
		return
	}

	if err := os.WriteFile(d.commentsFilePath(), data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing comments file: %v\n", err)
	}
}

func (d *Document) writeReviewMD(comments []Comment) {
	if len(comments) == 0 {
		os.Remove(d.reviewFilePath())
		return
	}

	reviewContent := GenerateReviewMD(d.Content, comments)

	if err := os.WriteFile(d.reviewFilePath(), []byte(reviewContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing review file: %v\n", err)
	}
}
