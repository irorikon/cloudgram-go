package channel

import (
	"errors"
	"strings"

	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/model"
	"gorm.io/gorm"
)

type ChannelService struct{}

var ChannelApp = new(ChannelService)

// CreateTGChannel 创建Telegram频道
func (cs *ChannelService) CreateTGChannel(channelId int64, channelName string) (*model.Channel, error) {
	channelName = strings.TrimSpace(channelName)

	// 检查是否已存在
	var existing model.Channel
	if err := config.DB.Where("channel_id = ?", channelId).First(&existing).Error; err == nil {
		return nil, errors.New("channel already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	channel := &model.Channel{
		ChannelID: channelId,
		Name:      channelName,
		Limited:   false,
		MessageID: 0,
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(channel).Error
	})

	return channel, err
}

// DeleteTGChannel 删除Telegram频道
func (cs *ChannelService) DeleteTGChannel(channelId int64) error {
	return config.DB.Where("channel_id = ?", channelId).Delete(&model.Channel{}).Error
}

// UpdateTGChannelMessageId 更新频道消息ID
func (cs *ChannelService) UpdateTGChannelMessageId(channelId int64, messageId int, updateLimit bool) (*model.Channel, error) {
	limited := messageId > 980000 && updateLimit

	err := config.DB.Model(&model.Channel{}).
		Where("channel_id = ?", channelId).
		Updates(map[string]any{
			"message_id": messageId,
			"limited":    limited,
		}).Error

	if err != nil {
		return nil, err
	}

	var channel model.Channel
	if err := config.DB.Where("channel_id = ?", channelId).First(&channel).Error; err != nil {
		return nil, errors.New("channel not found")
	}

	return &channel, nil
}

// UpdateTGChannelName 更新频道名称
func (cs *ChannelService) UpdateTGChannelName(channelId int64, channelName string) (*model.Channel, error) {
	channelName = strings.TrimSpace(channelName)

	if err := config.DB.Model(&model.Channel{}).
		Where("channel_id = ?", channelId).
		Update("name", channelName).Error; err != nil {
		return nil, err
	}

	var channel model.Channel
	if err := config.DB.Where("channel_id = ?", channelId).First(&channel).Error; err != nil {
		return nil, errors.New("channel not found")
	}

	return &channel, nil
}

// FindTGChannels 查询频道
func (cs *ChannelService) FindTGChannels(one, limited bool) (any, error) {
	if one {
		var channel model.Channel
		db := config.DB.Where("limited = ?", limited)
		err := db.First(&channel).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return &channel, err
	}

	var channels []model.Channel
	err := config.DB.Find(&channels).Error
	return channels, err
}
