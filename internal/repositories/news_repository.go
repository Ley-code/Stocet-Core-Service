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
	ID            *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title         string              `bson:"title" json:"title"`
	URL           string              `bson:"url" json:"url"`
	Source        string              `bson:"source" json:"source"`
	PublishedDate time.Time           `bson:"published_date" json:"published_date"`
	Content       string              `bson:"content" json:"content"`
	Summary       string              `bson:"summary,omitempty" json:"summary,omitempty"`
	Tags          []string            `bson:"tags" json:"tags"`
	ImageURL      string              `bson:"image_url,omitempty" json:"image_url,omitempty"`
}

type NewsRepository struct {
	collection *mongo.Collection
}

func NewNewsRepository(client *mongo.Client) *NewsRepository {
	db := client.Database(config.AppConfig.MongoDBDBName)
	return &NewsRepository{
		collection: db.Collection("articles"), // Collection name from news-aggregation service
	}
}

func (r *NewsRepository) GetNews(filters map[string]interface{}, limit, skip int) ([]NewsArticle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{}

	if category, ok := filters["category"].(string); ok && category != "" {
		// Category is stored in tags array in the news-aggregation service
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

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)).SetSort(bson.D{{Key: "published_date", Value: -1}})

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
