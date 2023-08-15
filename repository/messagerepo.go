package repository

import (
	"github.com/Coreychen4444/Lite_TikTok/model"
)

// 获取聊天记录
func (r *DbRepository) GetMessages(user_id, to_user_id, pre_msg_time int64) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Where(
		"(create_time > ?) and "+
			"((from_user_id = ? and to_user_id = ?) "+
			"or (from_user_id = ? and to_user_id = ?))",
		pre_msg_time, user_id, to_user_id, to_user_id, user_id).Order("create_time asc").Limit(20).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// 创建消息记录
func (r *DbRepository) CreateMessage(message *model.Message) error {
	if err := r.db.Create(message).Error; err != nil {
		return err
	}
	return nil
}
