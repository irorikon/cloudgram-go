package model

type Channel struct {
	BaseModel
	ChannelID int64  `json:"channel_id" gorm:"primaryKey;comment:频道ID"`
	Name      string `json:"name" gorm:"comment:频道名称"`
	Limited   bool   `json:"limited" gorm:"comment:是否限制"`
	MessageID int    `json:"message_id" gorm:"comment:消息ID"`
}

func (Channel) TableName() string {
	return "channels"
}
