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

type MarketPrice struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	TickerSymbol string             `bson:"ticker_symbol" json:"ticker_symbol"`
	Price        float64            `bson:"price" json:"price"`
	PercentChange float64            `bson:"percent_change" json:"percent_change"`
	Variation    float64            `bson:"variation" json:"variation"`
	Timestamp    time.Time          `bson:"timestamp" json:"timestamp"`
	Source       string             `bson:"source" json:"source"`
}

type MarketDataRepository struct {
	collection *mongo.Collection
}

func NewMarketDataRepository(client *mongo.Client) *MarketDataRepository {
	db := client.Database(config.AppConfig.MongoDBDBName)
	return &MarketDataRepository{
		collection: db.Collection("market_prices"),
	}
}

func (r *MarketDataRepository) GetLatestAll() ([]MarketPrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get the latest price for each ticker
	pipeline := []bson.M{
		{
			"$sort": bson.M{"timestamp": -1},
		},
		{
			"$group": bson.M{
				"_id":           "$ticker_symbol",
				"ticker_symbol": bson.M{"$first": "$ticker_symbol"},
				"price":         bson.M{"$first": "$price"},
				"percent_change": bson.M{"$first": "$percent_change"},
				"variation":      bson.M{"$first": "$variation"},
				"timestamp":      bson.M{"$first": "$timestamp"},
				"source":         bson.M{"$first": "$source"},
			},
		},
		{
			"$sort": bson.M{"ticker_symbol": 1},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prices []MarketPrice
	if err = cursor.All(ctx, &prices); err != nil {
		return nil, err
	}

	return prices, nil
}

func (r *MarketDataRepository) GetAllTickers() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	distinct, err := r.collection.Distinct(ctx, "ticker_symbol", bson.M{})
	if err != nil {
		return nil, err
	}

	tickers := make([]string, 0, len(distinct))
	for _, v := range distinct {
		if ticker, ok := v.(string); ok {
			tickers = append(tickers, ticker)
		}
	}

	return tickers, nil
}

func (r *MarketDataRepository) GetLatestByTicker(ticker string) (*MarketPrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var price MarketPrice
	err := r.collection.FindOne(ctx, bson.M{"ticker_symbol": ticker}, options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})).Decode(&price)
	if err != nil {
		return nil, err
	}

	return &price, nil
}

func (r *MarketDataRepository) GetByTickerAndDateRange(ticker string, startDate, endDate time.Time) ([]MarketPrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"ticker_symbol": ticker,
		"timestamp": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prices []MarketPrice
	if err = cursor.All(ctx, &prices); err != nil {
		return nil, err
	}

	return prices, nil
}
