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

type ArticleLike struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ArticleID string             `bson:"article_id"    json:"article_id"`
	UserID    string             `bson:"user_id"       json:"user_id"`
	CreatedAt time.Time          `bson:"created_at"    json:"created_at"`
}

type ArticleComment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ArticleID string             `bson:"article_id"    json:"article_id"`
	UserID    string             `bson:"user_id"       json:"user_id"`
	UserName  string             `bson:"user_name"     json:"user_name"`
	Content   string             `bson:"content"       json:"content"`
	CreatedAt time.Time          `bson:"created_at"    json:"created_at"`
}

type ArticleInteractionsRepository struct {
	likesCollection    *mongo.Collection
	commentsCollection *mongo.Collection
}

func NewArticleInteractionsRepository(client *mongo.Client) *ArticleInteractionsRepository {
	db := client.Database(config.AppConfig.MongoDBDBName)
	return &ArticleInteractionsRepository{
		likesCollection:    db.Collection("article_likes"),
		commentsCollection: db.Collection("article_comments"),
	}
}

func (r *ArticleInteractionsRepository) ToggleLike(articleID, userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"article_id": articleID,
		"user_id":    userID,
	}

	var existing ArticleLike
	err := r.likesCollection.FindOne(ctx, filter).Decode(&existing)
	if err == nil {
		_, delErr := r.likesCollection.DeleteOne(ctx, filter)
		return false, delErr
	}
	if err != mongo.ErrNoDocuments {
		return false, err
	}

	like := ArticleLike{
		ArticleID: articleID,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}
	_, err = r.likesCollection.InsertOne(ctx, like)
	return true, err
}

func (r *ArticleInteractionsRepository) CountLikes(articleID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.likesCollection.CountDocuments(ctx, bson.M{"article_id": articleID})
}

func (r *ArticleInteractionsRepository) IsLikedByUser(articleID, userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := r.likesCollection.CountDocuments(ctx, bson.M{
		"article_id": articleID,
		"user_id":    userID,
	})
	return count > 0, err
}

func (r *ArticleInteractionsRepository) CreateComment(articleID, userID, userName, content string) (*ArticleComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	comment := ArticleComment{
		ArticleID: articleID,
		UserID:    userID,
		UserName:  userName,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	result, err := r.commentsCollection.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		comment.ID = oid
	}

	return &comment, nil
}

func (r *ArticleInteractionsRepository) GetComments(articleID string, limit, skip int) ([]ArticleComment, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"article_id": articleID}

	total, err := r.commentsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(skip)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.commentsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var comments []ArticleComment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, 0, err
	}

	if comments == nil {
		comments = make([]ArticleComment, 0)
	}

	return comments, total, nil
}

func (r *ArticleInteractionsRepository) GetCommentByID(commentID string) (*ArticleComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, err
	}

	var comment ArticleComment
	err = r.commentsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *ArticleInteractionsRepository) DeleteComment(commentID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	_, err = r.commentsCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
