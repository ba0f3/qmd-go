package store

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func HashContent(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

func (s *Store) InsertContent(hash, content string, createdAt time.Time) error {
	_, err := s.DB.Exec(`INSERT OR IGNORE INTO content (hash, doc, created_at) VALUES (?, ?, ?)`,
		hash, content, createdAt.Format(time.RFC3339))
	return err
}

func (s *Store) InsertDocument(collection, path, title, hash string, createdAt, modifiedAt time.Time) error {
	_, err := s.DB.Exec(`
		INSERT INTO documents (collection, path, title, hash, created_at, modified_at, active)
		VALUES (?, ?, ?, ?, ?, ?, 1)
	`, collection, path, title, hash, createdAt.Format(time.RFC3339), modifiedAt.Format(time.RFC3339))
	return err
}

type Document struct {
	ID         int64
	Collection string
	Path       string
	Title      string
	Hash       string
	CreatedAt  time.Time
	ModifiedAt time.Time
	Active     bool
}

func (s *Store) FindActiveDocument(collection, path string) (*Document, error) {
	row := s.DB.QueryRow(`
		SELECT id, collection, path, title, hash, created_at, modified_at, active
		FROM documents
		WHERE collection = ? AND path = ? AND active = 1
	`, collection, path)

	var doc Document
	var active int
	var createdAt, modifiedAt string

	if err := row.Scan(&doc.ID, &doc.Collection, &doc.Path, &doc.Title, &doc.Hash, &createdAt, &modifiedAt, &active); err != nil {
		return nil, err
	}
	doc.Active = active == 1
	doc.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	doc.ModifiedAt, _ = time.Parse(time.RFC3339, modifiedAt)

	return &doc, nil
}

func (s *Store) UpdateDocument(id int64, title, hash string, modifiedAt time.Time) error {
	_, err := s.DB.Exec(`UPDATE documents SET title = ?, hash = ?, modified_at = ? WHERE id = ?`,
		title, hash, modifiedAt.Format(time.RFC3339), id)
	return err
}

func (s *Store) UpdateDocumentTitle(id int64, title string, modifiedAt time.Time) error {
	_, err := s.DB.Exec(`UPDATE documents SET title = ?, modified_at = ? WHERE id = ?`,
		title, modifiedAt.Format(time.RFC3339), id)
	return err
}

func (s *Store) DeactivateDocument(collection, path string) error {
	_, err := s.DB.Exec(`UPDATE documents SET active = 0 WHERE collection = ? AND path = ? AND active = 1`,
		collection, path)
	return err
}

func (s *Store) GetActiveDocumentPaths(collection string) ([]string, error) {
	rows, err := s.DB.Query(`SELECT path FROM documents WHERE collection = ? AND active = 1`, collection)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func (s *Store) CleanupOrphanedContent() (int64, error) {
	res, err := s.DB.Exec(`
		DELETE FROM content
		WHERE hash NOT IN (SELECT DISTINCT hash FROM documents WHERE active = 1)
	`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
