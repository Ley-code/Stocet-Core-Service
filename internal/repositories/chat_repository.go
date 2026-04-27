package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

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

// enrichMessages loads user display names and parent-message previews.
func (r *ChatRepository) enrichMessages(messages []models.Message) {
	for i := range messages {
		// Load user display name
		var user models.User
		if err := database.DB.Select("email", "full_name").Where("id = ?", messages[i].UserID).First(&user).Error; err == nil {
			messages[i].UserName = user.DisplayName()
		}

		// Load parent message preview for replies
		if messages[i].ParentID != nil && *messages[i].ParentID != "" {
			var parent models.Message
			if err := database.DB.Select("id", "content", "user_id", "is_deleted").Where("id = ?", *messages[i].ParentID).First(&parent).Error; err == nil {
				if parent.IsDeleted {
					messages[i].ParentContent = "This message was deleted"
				} else {
					content := parent.Content
					if len(content) > 80 {
						content = content[:80] + "…"
					}
					messages[i].ParentContent = content
				}
				var parentUser models.User
				if err := database.DB.Select("email", "full_name").Where("id = ?", parent.UserID).First(&parentUser).Error; err == nil {
					messages[i].ParentUser = parentUser.DisplayName()
				}
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
