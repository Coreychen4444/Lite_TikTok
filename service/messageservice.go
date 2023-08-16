package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/Coreychen4444/Lite_TikTok/repository"
)

// MessageService
type MessageService struct {
	r *repository.DbRepository
}

// NewMessageService
func NewMessageService(r *repository.DbRepository) *MessageService {
	return &MessageService{r: r}
}

// 获取聊天记录
func (s *MessageService) GetChatMessages(token, to_user_id, pre_msg_time string) ([]model.Message, error) {
	claims, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token无效,请重新登录")
	}
	toUserID, err := strconv.ParseInt(to_user_id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("用户id无效")
	}
	preMsgTime, err := strconv.ParseInt(pre_msg_time, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("时间戳格式错误")
	}
	messages, err := s.r.GetMessages(claims.UserID, toUserID, preMsgTime)
	if err != nil {
		return nil, fmt.Errorf("获取聊天记录失败")
	}
	return messages, nil
}

// 发送消息
func (s *MessageService) SendMessage(token, to_user_id, content string) error {
	claims, err := VerifyToken(token)
	if err != nil {
		return fmt.Errorf("token无效,请重新登录")
	}
	toUserID, err := strconv.Atoi(to_user_id)
	if err != nil {
		return fmt.Errorf("用户id无效")
	}
	message := model.Message{
		FromUserID: claims.UserID,
		ToUserID:   int64(toUserID),
		Content:    content,
		CreateTime: time.Now().Unix(),
	}
	err = s.r.CreateMessage(&message)
	if err != nil {
		return fmt.Errorf("发送消息失败")
	}
	return nil
}
