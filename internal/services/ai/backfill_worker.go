package ai

import (
	"context"
	"log"
	"strings"
	"time"

	"stocet-core-service/internal/models"
)

// StartBackfillWorker runs a periodic loop that scans for rows missing an
// embedding (the fire-and-forget goroutine in the admin path may have failed,
// or rows pre-date the AI feature) and embeds them in small batches.
//
// It runs for the life of the process, so start it once from main().
func StartBackfillWorker(ctx context.Context, embedder *Embedder, interval time.Duration, batchSize int) {
	if embedder == nil || !embedder.Configured() {
		log.Println("[ai] backfill worker skipped (embedder not configured)")
		return
	}
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	if batchSize <= 0 {
		batchSize = 25
	}

	log.Printf("[ai] backfill worker starting (interval=%s batch=%d)", interval, batchSize)
	t := time.NewTicker(interval)
	defer t.Stop()

	// Run once immediately so freshly deployed instances start filling in
	// embeddings without waiting a full interval.
	RunBackfillOnce(ctx, embedder, batchSize)

	for {
		select {
		case <-ctx.Done():
			log.Println("[ai] backfill worker stopping")
			return
		case <-t.C:
			RunBackfillOnce(ctx, embedder, batchSize)
		}
	}
}

// BackfillStats counts rows successfully embedded in a single RunBackfillOnce
// pass (one batch per table).
type BackfillStats struct {
	CorporateNewsRows   int
	CompanyProfilesRows int
}

// RunBackfillOnce processes one batch each from corporate_news and
// company_profiles. Errors are logged, never propagated — the loop is
// best-effort. Returns how many rows were embedded per table.
func RunBackfillOnce(ctx context.Context, embedder *Embedder, batchSize int) BackfillStats {
	var stats BackfillStats
	if embedder == nil || !embedder.Configured() {
		return stats
	}
	if batchSize <= 0 {
		batchSize = 25
	}

	if n, err := backfillCorporateNews(ctx, embedder, batchSize); err != nil {
		log.Printf("[ai] backfill corporate_news: %v", err)
	} else {
		stats.CorporateNewsRows = n
		if n > 0 {
			log.Printf("[ai] backfill corporate_news: embedded %d rows", n)
		}
	}

	if n, err := backfillCompanyProfiles(ctx, embedder, batchSize); err != nil {
		log.Printf("[ai] backfill company_profiles: %v", err)
	} else {
		stats.CompanyProfilesRows = n
		if n > 0 {
			log.Printf("[ai] backfill company_profiles: embedded %d rows", n)
		}
	}
	return stats
}

func backfillCorporateNews(ctx context.Context, e *Embedder, batch int) (int, error) {
	type row struct {
		ID    string
		Title string
		Body  *string
	}

	var rows []row
	if err := e.DB().WithContext(ctx).
		Raw(`SELECT id, title, body FROM corporate_news WHERE embedding IS NULL LIMIT ?`, batch).
		Scan(&rows).Error; err != nil {
		return 0, err
	}

	count := 0
	for _, r := range rows {
		if ctx.Err() != nil {
			return count, ctx.Err()
		}
		// Reuse the same composer as the synchronous hook.
		news := &models.CorporateNews{ID: r.ID, Title: r.Title, Body: r.Body}
		text := ComposeNewsText(news)
		if strings.TrimSpace(text) == "" {
			continue
		}
		if err := e.GenerateAndSaveEmbedding(ctx, "corporate_news", "id", r.ID, text); err != nil {
			log.Printf("[ai] backfill corporate_news id=%s: %v", r.ID, err)
			continue
		}
		count++
	}
	return count, nil
}

func backfillCompanyProfiles(ctx context.Context, e *Embedder, batch int) (int, error) {
	type row struct {
		ID               string
		TickerSymbol     string
		CompanyName      string
		Sector           *string
		NatureOfBusiness *string `gorm:"column:nature_of_business"`
	}

	var rows []row
	if err := e.DB().WithContext(ctx).
		Raw(`SELECT id, ticker_symbol, company_name, sector, nature_of_business
		     FROM company_profiles WHERE embedding IS NULL LIMIT ?`, batch).
		Scan(&rows).Error; err != nil {
		return 0, err
	}

	count := 0
	for _, r := range rows {
		if ctx.Err() != nil {
			return count, ctx.Err()
		}
		profile := &models.CompanyProfile{
			ID:               r.ID,
			TickerSymbol:     r.TickerSymbol,
			CompanyName:      r.CompanyName,
			Sector:           r.Sector,
			NatureOfBusiness: r.NatureOfBusiness,
		}
		text := ComposeProfileText(profile)
		if strings.TrimSpace(text) == "" {
			continue
		}
		if err := e.GenerateAndSaveEmbedding(ctx, "company_profiles", "id", r.ID, text); err != nil {
			log.Printf("[ai] backfill company_profiles id=%s: %v", r.ID, err)
			continue
		}
		count++
	}
	return count, nil
}
