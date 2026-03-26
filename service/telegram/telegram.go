package telegram

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type TelegramBotService struct{}

var UserTelegramBotApp = new(TelegramBotService)

// UploadFileToTelegram 上传文件到 Telegram
// ctx: 上下文
// chatID: 目标聊天 ID
// fileStream: 文件流
// fileName: 文件名
// isChunked: 是否为分块文件（当前未使用，保留扩展性）
func (t *TelegramBotService) UploadFileToTelegram(
	ctx context.Context,
	bot *telego.Bot,
	chatID int64,
	fileStream io.Reader,
	fileName string,
) (*telego.Message, error) {

	// 参数校验
	if bot == nil {
		return nil, errors.New("bot is not initialized")
	}
	if chatID == 0 {
		return nil, errors.New("chatID is required")
	}
	if fileStream == nil {
		return nil, errors.New("fileStream is required")
	}
	if strings.TrimSpace(fileName) == "" {
		return nil, errors.New("fileName is required")
	}

	// Telegram Bot API 文件大小限制
	const maxSize = 20 * 1024 * 1024 // 20MB

	// 使用 io.LimitReader 限制读取大小
	// 创建带大小检查的 Reader
	sizeCheckingReader := &sizeCheckingReader{
		reader:    fileStream,
		maxSize:   maxSize,
		totalRead: 0,
	}

	// 使用 telegoutil 工具包的正确方式上传文件
	params := tu.Document(
		tu.ID(chatID),
		tu.FileFromReader(sizeCheckingReader, fileName),
	)

	// 调用 API
	msg, err := bot.SendDocument(ctx, params.WithDisableContentTypeDetection())
	if err != nil {
		return nil, fmt.Errorf("Telegram API 错误: %w", err)
	}

	return msg, nil
}

// 带大小检查的 Reader
type sizeCheckingReader struct {
	reader    io.Reader
	maxSize   int64
	totalRead int64
}

func (sr *sizeCheckingReader) Read(p []byte) (n int, err error) {
	// 计算还能读取的字节数
	remaining := sr.maxSize - sr.totalRead
	if remaining <= 0 {
		return 0, fmt.Errorf("文件大小超过 %dMB 限制", sr.maxSize/(1024*1024))
	}

	// 限制本次读取的大小
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}

	n, err = sr.reader.Read(p)
	sr.totalRead += int64(n)

	return n, err
}

func (t *TelegramBotService) GetTelegramFilePath(ctx context.Context, bot *telego.Bot, fileID string) (path string, err error) {
	// 校验参数
	if strings.TrimSpace(fileID) == "" {
		err = errors.New("fileName is required")
		return
	}
	fileInfo, err := bot.GetFile(ctx, &telego.GetFileParams{FileID: fileID})
	if err != nil {
		return
	}

	return fileInfo.FilePath, nil
}

func (t *TelegramBotService) GetTelegramFileUrl(ctx context.Context, bot *telego.Bot, fileID string) (url string, err error) {
	path, err := t.GetTelegramFilePath(ctx, bot, fileID)
	if err != nil {
		return
	}
	return bot.FileDownloadURL(path), err
}

func (t *TelegramBotService) DownloadTelegramFile(ctx context.Context, bot *telego.Bot, fileID string) ([]byte, error) {
	fileUrl, err := t.GetTelegramFileUrl(ctx, bot, fileID)
	if err != nil {
		return nil, err
	}
	return tu.DownloadFile(fileUrl)
}

func (t *TelegramBotService) DeleteTelegramFile(ctx context.Context, bot *telego.Bot, chatID int64, messageId int) (err error) {
	if messageId == 0 {
		return errors.New("messageId is empty")
	}
	if chatID == 0 {
		return errors.New("chatID is required")
	}

	err = bot.DeleteMessage(ctx, &telego.DeleteMessageParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageId,
	})
	if err != nil {
		return
	}
	return
}

func (t *TelegramBotService) DeleteTelegramFiles(ctx context.Context, bot *telego.Bot, chatID int64, messageIds []int) (err error) {
	if messageIds == nil {
		return errors.New("messageIds is empty")
	}
	if chatID == 0 {
		return errors.New("chatID is required")
	}

	err = bot.DeleteMessages(ctx, &telego.DeleteMessagesParams{
		ChatID:     tu.ID(chatID),
		MessageIDs: messageIds,
	})
	if err != nil {
		return
	}
	return
}
func (t *TelegramBotService) CheckMe(ctx context.Context, bot *telego.Bot) (err error) {
	_, err = bot.GetMe(ctx)
	if err != nil {
		return
	}
	return
}

func (t *TelegramBotService) CheckChannel(ctx context.Context, bot *telego.Bot, channelID int64) (err error) {
	_, err = bot.GetChat(ctx, &telego.GetChatParams{
		ChatID: tu.ID(channelID),
	})
	if err != nil {
		return
	}
	return
}
