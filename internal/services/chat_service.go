package services

import (
	"errors"
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
)

type ChatService struct {
	repo *repositories.ChatRepository
}

func NewChatService(repo *repositories.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) CreateChannel(channel *models.Channel) error {
	if channel.Name == "" {
		return errors.New("channel name is required")
	}
	return s.repo.CreateChannel(channel)
}

func (s *ChatService) GetChannelByID(id string) (*models.Channel, error) {
	return s.repo.GetChannelByID(id)
}

func (s *ChatService) ListChannels() ([]models.Channel, error) {
	return s.repo.ListChannels()
}

func (s *ChatService) CreateMessage(channelID, userID, content string) (*models.Message, error) {
	if len(content) > config.CHAT_MESSAGE_MAX_LENGTH {
		return nil, errors.New("message too long")
	}

	message := &models.Message{
		ChannelID: channelID,
		UserID:    userID,
		Content:   content,
	}

	if err := s.repo.CreateMessage(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *ChatService) GetMessages(channelID string, limit, offset int) ([]models.Message, int64, error) {
	return s.repo.GetMessagesByChannel(channelID, limit, offset)
}

func (s *ChatService) GetRecentMessages(channelID string, limit int) ([]models.Message, error) {
	return s.repo.GetRecentMessages(channelID, limit)
}
