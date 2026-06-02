package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"stocet-core-service/internal/config"
)

type NewsArticle struct {
	ID                   *primitive.ObjectID `bson:"_id,omitempty"                  json:"id,omitempty"`
	Title                string              `bson:"title"                          json:"title"`
	URL                  string              `bson:"url"                            json:"url"`
	Source               string              `bson:"source"                         json:"source"`
	PublishedDate        time.Time           `bson:"published_date"                 json:"published_date"`
	Content              string              `bson:"content"                        json:"content"`
	Summary              string              `bson:"summary,omitempty"              json:"summary,omitempty"`
	Tags                 []string            `bson:"tags"                           json:"tags"`
	ImageURL             string              `bson:"image_url,omitempty"            json:"image_url,omitempty"`
	// Sentiment fields — set by the sentiment backfill worker
	Sentiment            string              `bson:"sentiment,omitempty"            json:"sentiment,omitempty"`
	SentimentTickers     []string            `bson:"sentiment_tickers,omitempty"    json:"sentiment_tickers,omitempty"`
	SentimentAnalyzedAt  *time.Time          `bson:"sentiment_analyzed_at,omitempty" json:"sentiment_analyzed_at,omitempty"`
}

type NewsRepository struct {
	collection *mongo.Collection
}

func NewNewsRepository(client *mongo.Client) *NewsRepository {
	db := client.Database(config.AppConfig.MongoDBDBName)
	return &NewsRepository{
		collection: db.Collection("articles"),
	}
}

func buildNewsFilter(filters map[string]interface{}) bson.M {
	filter := bson.M{}

	if category, ok := filters["category"].(string); ok && category != "" {
		filter["tags"] = category
	}

	if source, ok := filters["source"].(string); ok && source != "" {
		filter["source"] = source
	}

	if startDate, ok := filters["start_date"].(time.Time); ok {
		filter["published_date"] = bson.M{"$gte": startDate}
	}

	if endDate, ok := filters["end_date"].(time.Time); ok {
		if existingFilter, ok := filter["published_date"].(bson.M); ok {
			existingFilter["$lte"] = endDate
		} else {
			filter["published_date"] = bson.M{"$lte": endDate}
		}
	}

	if searchText, ok := filters["search_text"].(string); ok && searchText != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": searchText, "$options": "i"}},
			{"content": bson.M{"$regex": searchText, "$options": "i"}},
		}
	}

	// Sentiment filter: "bullish", "bearish", "neutral", or "analyzed" (any tagged article)
	if sentiment, ok := filters["sentiment"].(string); ok && sentiment != "" {
		if sentiment == "analyzed" {
			filter["sentiment"] = bson.M{"$exists": true, "$ne": ""}
		} else {
			filter["sentiment"] = sentiment
		}
	}

	// Ticker filter: match articles that tagged this ticker as affected
	if ticker, ok := filters["ticker"].(string); ok && ticker != "" {
		filter["sentiment_tickers"] = ticker
	}

	return filter
}

func (r *NewsRepository) CountNews(filters map[string]interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := buildNewsFilter(filters)
	return r.collection.CountDocuments(ctx, filter)
}

// GetNews returns articles ordered by published_date desc, then _id desc.
// When cursorPublished and cursorID are non-nil, returns the next page after that composite key (keyset).
// Otherwise uses skip/limit offset pagination.
func (r *NewsRepository) GetNews(filters map[string]interface{}, limit, skip int, cursorPublished *time.Time, cursorID *primitive.ObjectID) ([]NewsArticle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := buildNewsFilter(filters)

	if cursorPublished != nil && cursorID != nil {
		keyset := bson.M{
			"$or": []bson.M{
				{"published_date": bson.M{"$lt": *cursorPublished}},
				{
					"$and": []bson.M{
						{"published_date": *cursorPublished},
						{"_id": bson.M{"$lt": *cursorID}},
					},
				},
			},
		}
		filter = bson.M{"$and": []bson.M{filter, keyset}}
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "published_date", Value: -1}, {Key: "_id", Value: -1}})

	if cursorPublished == nil || cursorID == nil {
		opts = opts.SetSkip(int64(skip))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []NewsArticle
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, err
	}

	return articles, nil
}

// GetUnanalyzed returns articles without a sentiment field, oldest first.
// Used by the sentiment backfill worker.
func (r *NewsRepository) GetUnanalyzed(limit int) ([]NewsArticle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Articles that have no sentiment field OR an empty one
	filter := bson.M{"$or": []bson.M{
		{"sentiment": bson.M{"$exists": false}},
		{"sentiment": ""},
	}}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "published_date", Value: -1}}) // process newest first

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []NewsArticle
	return articles, cursor.All(ctx, &articles)
}

// UpdateSentiment writes the AI-classified sentiment back to an article document.
func (r *NewsRepository) UpdateSentiment(ctx context.Context, id primitive.ObjectID, sentiment string, tickers []string) error {
	if tickers == nil {
		tickers = []string{}
	}
	now := time.Now().UTC()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"sentiment":             sentiment,
			"sentiment_tickers":     tickers,
			"sentiment_analyzed_at": now,
		}},
	)
	return err
}

func (r *NewsRepository) DeleteNews(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
