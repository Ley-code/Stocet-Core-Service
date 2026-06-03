package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"

	"github.com/google/uuid"
)

const parentPreviewMaxLen = 80

type ChatRepository struct{}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{}
}

// Channel methods
func (r *ChatRepository) CreateChannel(channel *models.Channel) error {
	return database.DB.Create(channel).Error
}

func (r *ChatRepository) GetChannelByID(id string) (*models.Channel, error) {
	var channel models.Channel
	err := database.DB.Where("id = ?", id).First(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *ChatRepository) ListChannels() ([]models.Channel, error) {
	var channels []models.Channel
	err := database.DB.Find(&channels).Error
	return channels, err
}

func (r *ChatRepository) UpdateChannel(channel *models.Channel) error {
	return database.DB.Save(channel).Error
}

func (r *ChatRepository) DeleteChannel(id string) error {
	return database.DB.Delete(&models.Channel{}, "id = ?", id).Error
}

// Message methods
func (r *ChatRepository) CreateMessage(message *models.Message) error {
	return database.DB.Create(message).Error
}

func (r *ChatRepository) GetMessageByID(id string) (*models.Message, error) {
	var message models.Message
	err := database.DB.Where("id = ?", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *ChatRepository) UpdateMessage(message *models.Message) error {
	return database.DB.Save(message).Error
}

// enrichMessages loads user display names and parent-message previews in batch.
func (r *ChatRepository) enrichMessages(messages []models.Message) {
	if len(messages) == 0 {
		return
	}

	userIDSet := make(map[string]struct{})
	parentIDs := make([]string, 0)
	for i := range messages {
		if messages[i].UserID != uuid.Nil {
			userIDSet[messages[i].UserID.String()] = struct{}{}
		}
		if messages[i].ParentID != nil && *messages[i].ParentID != "" {
			parentIDs = append(parentIDs, *messages[i].ParentID)
		}
	}

	type parentPreview struct {
		content  string
		userID   string
	}
	parentPreviews := make(map[string]parentPreview)
	if len(parentIDs) > 0 {
		var parents []models.Message
		if err := database.DB.Select("id", "content", "user_id", "is_deleted").
			Where("id IN ?", parentIDs).Find(&parents).Error; err == nil {
			for _, p := range parents {
				if p.UserID != uuid.Nil {
					userIDSet[p.UserID.String()] = struct{}{}
				}
				var content string
				if p.IsDeleted {
					content = "This message was deleted"
				} else {
					content = p.Content
					if len(content) > parentPreviewMaxLen {
						content = content[:parentPreviewMaxLen] + "…"
					}
				}
				parentPreviews[p.ID] = parentPreview{content: content, userID: p.UserID.String()}
			}
		}
	}

	userNames := make(map[string]string)
	if len(userIDSet) > 0 {
		userIDs := make([]string, 0, len(userIDSet))
		for id := range userIDSet {
			userIDs = append(userIDs, id)
		}
		var users []models.User
		if err := database.DB.Select("id", "email", "full_name").Where("id IN ?", userIDs).Find(&users).Error; err == nil {
			for _, u := range users {
				userNames[u.ID] = u.DisplayName()
			}
		}
	}

	for i := range messages {
		if name, ok := userNames[messages[i].UserID.String()]; ok {
			messages[i].UserName = name
		}
		if messages[i].ParentID != nil {
			if prev, ok := parentPreviews[*messages[i].ParentID]; ok {
				messages[i].ParentContent = prev.content
				messages[i].ParentUser = userNames[prev.userID]
			}
		}
	}
}

func (r *ChatRepository) GetMessagesByChannel(channelID string, limit, offset int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	if err := database.DB.Model(&models.Message{}).Where("channel_id = ?", channelID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := database.DB.Where("channel_id = ?", channelID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	r.enrichMessages(messages)

	return messages, total, err
}

func (r *ChatRepository) GetRecentMessages(channelID string, limit int) ([]models.Message, error) {
	var messages []models.Message
	err := database.DB.Where("channel_id = ?", channelID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	r.enrichMessages(messages)

	return messages, err
}
