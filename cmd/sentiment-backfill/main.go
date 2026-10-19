// Command sentiment-backfill classifies articles that have no sentiment yet, using the same
// classifier as the Core service (real-ticker list, JSON-mode, no "neutral on error" fallback).
//
//	go run ./cmd/sentiment-backfill                 # dry run: shows what it would do
//	go run ./cmd/sentiment-backfill -apply          # classify everything waiting
//	go run ./cmd/sentiment-backfill -apply -reset   # first clear fake "neutral, no confidence" rows
//
// It stops at the first batch that is not fully processed (quota exhausted / API errors), so it never
// burns quota in a loop; just re-run it later to continue.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"stocet-core-service/internal/config"
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/services/ai"
)

func main() {
	apply := flag.Bool("apply", false, "write results (default is a dry run)")
	reset := flag.Bool("reset", false, "first clear sentiment stamped by a failure (neutral with no confidence)")
	batch := flag.Int("batch", 20, "articles per batch")
	max := flag.Int("max", 5000, "stop after this many articles")
	pause := flag.Duration("pause", 0, "wait between batches (free-tier limits are per minute; e.g. 45s with -batch 10)")
	model := flag.String("model", "", "override the chat model for this run only (e.g. gemini-flash-lite-latest); default: GEMINI_CHAT_MODEL")
	backup := flag.String("backup", "sentiment_reset_backup.json", "where to save reset rows before clearing them")
	flag.Parse()

	if err := config.LoadConfig(); err != nil {
		log.Fatal("config: ", err)
	}
	if err := database.InitDB(); err != nil {
		log.Fatal("database: ", err)
	}
	db := database.DB
	repo := repositories.NewPostgresNewsRepository(db)
	chatModel := config.AppConfig.GeminiChatModel
	if *model != "" {
		chatModel = *model
	}
	fmt.Println("chat model:", chatModel)
	client := ai.NewGeminiClient(config.AppConfig.GeminiAPIKey, chatModel,
		config.AppConfig.GeminiEmbedModel, time.Duration(config.AppConfig.AIRequestTimeoutSec)*time.Second)
	if !client.Configured() {
		log.Fatal("GEMINI_API_KEY not set")
	}

	var fake, waiting, total int64
	db.Raw(`SELECT count(*) FROM articles WHERE sentiment IS NOT NULL AND sentiment <> '' AND sentiment_confidence IS NULL`).Scan(&fake)
	db.Raw(`SELECT count(*) FROM articles WHERE sentiment IS NULL OR sentiment = ''`).Scan(&waiting)
	db.Raw(`SELECT count(*) FROM articles`).Scan(&total)
	known, _ := repo.GetKnownTickers()
	fmt.Printf("articles: %d | waiting for sentiment: %d | fake neutral rows to reset: %d | known tickers: %v\n", total, waiting, fake, known)
	if !*apply {
		fmt.Println("DRY RUN - nothing changed. Re-run with -apply (and -reset to redo the fake rows).")
		return
	}

	if *reset && fake > 0 {
		var rows []struct {
			ID               string
			Sentiment        string
			SentimentTickers string
		}
		db.Raw(`SELECT id::text AS id, sentiment, COALESCE(sentiment_tickers::text,'') AS sentiment_tickers FROM articles
		        WHERE sentiment IS NOT NULL AND sentiment <> '' AND sentiment_confidence IS NULL`).Scan(&rows)
		b, _ := json.MarshalIndent(rows, "", " ")
		if err := os.WriteFile(*backup, b, 0o644); err != nil {
			log.Fatal("could not write backup, not resetting: ", err)
		}
		res := db.Exec(`UPDATE articles SET sentiment=NULL, sentiment_confidence=NULL, sentiment_tickers=NULL, sentiment_analyzed_at=NULL
		                WHERE sentiment IS NOT NULL AND sentiment <> '' AND sentiment_confidence IS NULL`)
		if res.Error != nil {
			log.Fatal("reset failed: ", res.Error)
		}
		fmt.Printf("reset %d fake rows (backup: %s)\n", res.RowsAffected, *backup)
	}

	ctx := context.Background()
	done, stalls := 0, 0
	for done < *max {
		n, err := ai.RunSentimentBackfillOnce(ctx, client, repo, *batch)
		if err != nil {
			log.Printf("batch error: %v", err)
			break
		}
		done += n
		fmt.Printf("  classified %d so far\n", done)
		if n == 0 {
			// Nothing processed: either nothing is left, or the API is refusing us.
			var waitingNow int64
			db.Raw(`SELECT count(*) FROM articles WHERE sentiment IS NULL OR sentiment = ''`).Scan(&waitingNow)
			if waitingNow == 0 {
				break
			}
			stalls++
			if stalls >= 4 {
				fmt.Println("no progress for 4 batches in a row - daily quota is probably exhausted; stopping (re-run later)")
				break
			}
			fmt.Println("  rate limited, waiting 75s...")
			time.Sleep(75 * time.Second)
			continue
		}
		stalls = 0
		if n < *batch {
			// partial batch: the per-minute limit was hit mid-way; give it a minute to recover
			time.Sleep(75 * time.Second)
			continue
		}
		time.Sleep(*pause)
	}
	var left int64
	db.Raw(`SELECT count(*) FROM articles WHERE sentiment IS NULL OR sentiment = ''`).Scan(&left)
	fmt.Printf("done: classified %d this run | still waiting: %d\n", done, left)
}
