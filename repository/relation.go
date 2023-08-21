package repository

import (
	"sync"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"gorm.io/gorm"
)

// 关注
func (r *DbRepository) CreateFollow(authorID, fansID int64) error {
	relation := model.Relation{
		AuthorID: authorID,
		FansID:   fansID,
	}
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var txErr error
	defer func() {
		if txErr != nil || tx.Commit().Error != nil {
			tx.Rollback()
		}
	}()
	// 添加关注记录
	if err := r.db.Create(&relation).Error; err != nil {
		txErr = err
		return err
	}
	// 更新作者的粉丝数
	if err := r.db.Model(&model.User{}).Where("id = ?", authorID).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	// 更新用户的关注数
	if err := r.db.Model(&model.User{}).Where("id = ?", fansID).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	return tx.Commit().Error
}

// 取消关注
func (r *DbRepository) DeleteFollow(authorID, fansID int64) error {
	var relation model.Relation
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var txErr error
	defer func() {
		if txErr != nil || tx.Commit().Error != nil {
			tx.Rollback()
		}
	}()

	// 删除关注记录
	if err := r.db.Where("author_id = ? and fans_id = ?", authorID, fansID).Delete(&relation).Error; err != nil {
		txErr = err
		return err
	}
	// 更新作者的粉丝数
	if err := r.db.Model(&model.User{}).Where("id = ?", authorID).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	// 更新用户的关注数
	if err := r.db.Model(&model.User{}).Where("id = ?", fansID).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	return tx.Commit().Error
}

// 获取用户关注列表
func (r *DbRepository) GetFollowList(userID int64) ([]model.User, error) {
	var following_id []int64
	err := r.db.Model(&model.Relation{}).Where("fans_id = ?", userID).Pluck("author_id", &following_id).Error
	if err != nil {
		return nil, err
	}
	var followings []model.User
	err = r.db.Where("id in (?)", following_id).Find(&followings).Error
	if err != nil {
		return nil, err
	}
	return followings, nil
}

// 判断是否关注
func (r *DbRepository) IsFollow(authorID, fansID int64) (bool, error) {
	var relation model.Relation
	err := r.db.Where("author_id = ? and fans_id = ?", authorID, fansID).First(&relation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 获取用户粉丝列表
func (r *DbRepository) GetFansList(userID int64) ([]model.User, error) {
	var follower_id []int64
	err := r.db.Model(&model.Relation{}).Where("author_id = ?", userID).Pluck("fans_id", &follower_id).Error
	if err != nil {
		return nil, err
	}
	var followers []model.User
	err = r.db.Where("id in (?)", follower_id).Find(&followers).Error
	if err != nil {
		return nil, err
	}
	return followers, nil
}

// 获取用户好友列表
// 由于涉及多次数据库查询，使用协程并发查询，并且使用redis缓存好友id,减少数据库IO
func (r *DbRepository) GetFriendList(userID int64) ([]model.User, error) {
	var following, follower []int64
	errChan := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		err := r.db.Model(&model.Relation{}).Where("fans_id = ?", userID).Pluck("author_id", &following).Error
		if err != nil {
			errChan <- err
		}
		wg.Done()
	}()
	go func() {
		err := r.db.Model(&model.Relation{}).Where("author_id = ?", userID).Pluck("fans_id", &follower).Error
		if err != nil {
			errChan <- err
		}
		wg.Done()
	}()
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	friendID := GetFriID(following, follower)
	friends, err := r.GetUserListByIds(friendID)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

// 获取好友ID
func GetFriID(following, follower []int64) []int64 {
	var friendID []int64
	friendmap := make(map[int64]bool)
	for _, v := range following {
		friendmap[v] = true
	}
	for _, v := range follower {
		if friendmap[v] {
			friendID = append(friendID, v)
		}
	}
	return friendID
}
