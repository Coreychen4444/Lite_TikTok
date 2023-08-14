package service

import (
	"fmt"
	"strconv"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/Coreychen4444/Lite_TikTok/repository"
)

type RelationService struct {
	r *repository.DbRepository
}

func NewRelationService(r *repository.DbRepository) *RelationService {
	return &RelationService{r: r}
}

// 关注或取消关注
func (s *RelationService) FollowOrCancel(token, author_id, action_type string) error {
	clamis, err := VerifyToken(token)
	if err != nil {
		return fmt.Errorf("token无效,请重新登录")
	}
	authorID, err := strconv.Atoi(author_id)
	if err != nil {
		return fmt.Errorf("作者id无效")
	}
	if action_type == "1" {
		err = s.r.CreateFollow(int64(authorID), clamis.UserID)
		if err != nil {
			return fmt.Errorf("关注失败")
		}
	}
	if action_type == "2" {
		err = s.r.DeleteFollow(int64(authorID), clamis.UserID)
		if err != nil {
			return fmt.Errorf("取消关注失败")
		}
	}
	return nil
}

// 获取用户关注列表
func (s *RelationService) GetFollowings(token, user_id string) ([]model.User, error) {
	_, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token无效,请重新登录")
	}
	userID, err := strconv.Atoi(user_id)
	if err != nil {
		return nil, fmt.Errorf("用户id无效")
	}
	followings, err := s.r.GetFollowList(int64(userID))
	if err != nil {
		return nil, fmt.Errorf("获取关注列表失败")
	}
	return followings, nil
}

// 获取用户粉丝列表
func (s *RelationService) GetFollowers(token, user_id string) ([]model.User, error) {
	_, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token无效,请重新登录")
	}
	userID, err := strconv.Atoi(user_id)
	if err != nil {
		return nil, fmt.Errorf("用户id无效")
	}
	followers, err := s.r.GetFansList(int64(userID))
	if err != nil {
		return nil, fmt.Errorf("获取粉丝列表失败")
	}
	return followers, nil
}

// 获取用户好友列表
func (s *RelationService) GetFriends(token, user_id string) ([]model.User, error) {
	_, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token无效,请重新登录")
	}
	userID, err := strconv.Atoi(user_id)
	if err != nil {
		return nil, fmt.Errorf("用户id无效")
	}
	friends, err := s.r.GetFriendList(int64(userID))
	if err != nil {
		return nil, fmt.Errorf("获取好友列表失败")
	}
	return friends, nil
}
