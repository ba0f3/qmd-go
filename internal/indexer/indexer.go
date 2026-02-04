package indexer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/tobi/qmd-go/internal/store"
)

func IndexFiles(s *store.Store, collectionName, rootPath, pattern string) error {
	fsys := os.DirFS(rootPath)

	now := time.Now()

	// Use doublestar to find files
	files, err := doublestar.Glob(fsys, pattern)
	if err != nil {
		return err
	}

	indexedCount := 0
	updatedCount := 0
	seenPaths := make(map[string]bool)

	for _, relPath := range files {
		fullPath := filepath.Join(rootPath, relPath)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		if info.IsDir() {
			continue
		}

		seenPaths[relPath] = true

		contentBytes, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		content := string(contentBytes)
		hash := store.HashContent(content)
		title := filepath.Base(relPath) // Simplified title extraction

		// Check if exists
		doc, err := s.FindActiveDocument(collectionName, relPath)
		if err == nil {
			// Update if changed
			if doc.Hash != hash {
				s.InsertContent(hash, content, now)
				s.UpdateDocument(doc.ID, title, hash, now)
				updatedCount++
			}
		} else {
			// Insert new
			s.InsertContent(hash, content, now)
			s.InsertDocument(collectionName, relPath, title, hash, info.ModTime(), now)
			indexedCount++
		}
	}

	// Handle deletions
	activePaths, err := s.GetActiveDocumentPaths(collectionName)
	removedCount := 0
	if err == nil {
		for _, path := range activePaths {
			if !seenPaths[path] {
				s.DeactivateDocument(collectionName, path)
				removedCount++
			}
		}
	}

	// Cleanup orphans
	s.CleanupOrphanedContent()

	fmt.Printf("Collection '%s': Indexed %d new, Updated %d, Removed %d.\n", collectionName, indexedCount, updatedCount, removedCount)
	return nil
}
