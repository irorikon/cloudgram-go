package channel

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/model/request"
	"github.com/irorikon/cloudgram-go/model/response"
	"github.com/irorikon/cloudgram-go/utils"
)

type ChannelApi struct{}

// 获取所有 Telegram Channel 列表
func (ch *ChannelApi) ListChannel(c *gin.Context) {
	channelList, err := channelService.FindTGChannels(false, false)
	if err != nil {
		response.FailWithMessage("Get channel failed", c)
		return
	}
	response.OKWithData(channelList, c)
}

// 获取一个可用的 Telegram Channel
func (ch *ChannelApi) GetChannel(c *gin.Context) {
	channel, err := channelService.FindTGChannels(true, false)
	if err != nil {
		response.FailWithMessage("Get channel failed", c)
		return
	}
	response.OKWithData(channel, c)
}

// 创建一个 Telegram Channel
func (ch *ChannelApi) CreateChannel(c *gin.Context) {
	var req request.ChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("Request body must be valid JSON", c)
		return
	}
	if err := utils.Verify(req, utils.ChannelVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	channel, err := channelService.CreateTGChannel(req.ChannelID, req.ChannelName)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OKWithData(channel, c)
}

// 修改频道名称
func (ch *ChannelApi) UpdateChannelName(c *gin.Context) {
	var req request.ChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("Request body must be valid JSON", c)
		return
	}
	if err := utils.Verify(req, utils.ChannelVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	channel, err := channelService.UpdateTGChannelName(req.ChannelID, req.ChannelName)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OKWithData(channel, c)
}

// 删除 TG 频道
func (ch *ChannelApi) DeleteTGChannel(c *gin.Context) {
	var req request.ChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("Request body must be valid JSON", c)
		return
	}
	if req.ChannelID == 0 {
		response.FailWithMessage("channel_id is required and must be a non-zero integer", c)
		return
	}
	err := channelService.DeleteTGChannel(req.ChannelID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OK(c)
}

// Check Telegram Channel
func (ch *ChannelApi) CheckTelegramChannel(c *gin.Context) {
	var req request.ChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("Request body must be valid JSON", c)
		return
	}
	if req.ChannelID == 0 {
		response.FailWithMessage("channel_id is required and must be a non-zero integer", c)
		return
	}
	fmt.Println(config.TelegramBotToken)
	err := telegramService.CheckChannel(c, config.TeleBot, req.ChannelID)
	if err != nil {
		response.FailWithMessage("Channel not found or deleted", c)
		return
	}
	response.OKWithData(gin.H{
		"channel_status": true,
	}, c)
}
