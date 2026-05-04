package ai

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"stocet-core-service/internal/models"

	"gorm.io/gorm"
)

// Embedder owns the Gemini client + DB handle and persists embeddings.
//
// We keep this separate from the chat/RAG services because:
//   - it's used both synchronously (admin write hooks fire-and-forget) and
//     by the periodic backfill worker
//   - it can be safely nil when GEMINI_API_KEY is unset (no-op)
type Embedder struct {
	db     *gorm.DB
	client *GeminiClient
}

// NewEmbedder returns nil if the gemini client is not configured. Callers
// should treat a nil Embedder as "skip embedding work".
func NewEmbedder(db *gorm.DB, client *GeminiClient) *Embedder {
	if db == nil || client == nil || !client.Configured() {
		return nil
	}
	return &Embedder{db: db, client: client}
}

// Configured reports whether the embedder will attempt to run.
func (e *Embedder) Configured() bool { return e != nil && e.client != nil && e.client.Configured() }

// Client exposes the underlying GeminiClient (used by RAG search to embed
// user queries with the same model).
func (e *Embedder) Client() *GeminiClient { return e.client }

// DB exposes the underlying gorm.DB. Used by the backfill worker.
func (e *Embedder) DB() *gorm.DB { return e.db }

// GenerateAndSaveEmbedding embeds `textToEmbed` and writes the vector to
// `tableName` for the row identified by `idColumn = id`. The query uses an
// explicit `::vector` cast so we don't need a custom GORM type.
//
// On error, the caller (typically an admin-write goroutine) should log only,
// because failures here are self-healing via the backfill worker.
func (e *Embedder) GenerateAndSaveEmbedding(ctx context.Context, tableName, idColumn, id, textToEmbed string) error {
	if !e.Configured() {
		return errors.New("embedder: not configured")
	}
	tableName = strings.TrimSpace(tableName)
	idColumn = strings.TrimSpace(idColumn)
	if !isSafeIdent(tableName) || !isSafeIdent(idColumn) {
		return fmt.Errorf("embedder: unsafe identifier")
	}
	textToEmbed = strings.TrimSpace(textToEmbed)
	if textToEmbed == "" {
		return errors.New("embedder: empty text")
	}

	vec, err := e.client.Embed(ctx, textToEmbed)
	if err != nil {
		return err
	}

	// pgvector wire format: [v1,v2,...]
	literal := vectorLiteral(vec)

	// SQL is constructed from validated identifiers above; values are still parameterized.
	sql := fmt.Sprintf("UPDATE %s SET embedding = ?::vector WHERE %s = ?", tableName, idColumn)
	return e.db.WithContext(ctx).Exec(sql, literal, id).Error
}

// ComposeNewsText builds the embedding text for a corporate_news row.
// Order: Title, then Body. Body may be nil; if so we fall back to title only.
func ComposeNewsText(n *models.CorporateNews) string {
	if n == nil {
		return ""
	}
	parts := []string{strings.TrimSpace(n.Title)}
	if n.Body != nil {
		body := strings.TrimSpace(*n.Body)
		if body != "" {
			parts = append(parts, body)
		}
	}
	return strings.Join(parts, "\n\n")
}

// ComposeProfileText builds the embedding text for a company_profiles row.
// Concatenates ticker, name, and the most "description-like" fields the
// schema offers (sector, nature_of_business). Skips empty/nil parts.
func ComposeProfileText(p *models.CompanyProfile) string {
	if p == nil {
		return ""
	}
	ticker := strings.TrimSpace(p.TickerSymbol)
	name := strings.TrimSpace(p.CompanyName)

	var b strings.Builder
	if ticker != "" {
		b.WriteString(ticker)
	}
	if name != "" {
		if b.Len() > 0 {
			b.WriteString(" - ")
		}
		b.WriteString(name)
	}
	if p.Sector != nil {
		s := strings.TrimSpace(*p.Sector)
		if s != "" {
			b.WriteString(" (")
			b.WriteString(s)
			b.WriteString(")")
		}
	}
	if p.NatureOfBusiness != nil {
		nb := strings.TrimSpace(*p.NatureOfBusiness)
		if nb != "" {
			if b.Len() > 0 {
				b.WriteString(". ")
			}
			b.WriteString(nb)
		}
	}
	return b.String()
}

// vectorLiteral renders a []float32 as the pgvector text format: [v1,v2,...].
// Uses 'g' formatting with reasonable precision to keep payloads small.
func vectorLiteral(v []float32) string {
	var b strings.Builder
	b.Grow(len(v) * 10)
	b.WriteByte('[')
	for i, x := range v {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatFloat(float64(x), 'g', -1, 32))
	}
	b.WriteByte(']')
	return b.String()
}

// isSafeIdent permits only [A-Za-z0-9_]. Used as a defensive check before
// interpolating identifiers into raw SQL.
func isSafeIdent(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '_':
		default:
			return false
		}
	}
	return true
}
