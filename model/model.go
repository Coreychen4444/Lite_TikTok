package model

import (
	"strconv"

	"gorm.io/gorm"
)

// User
type User struct {
	ID              int64  `json:"id"`               // 用户id
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Name            string `json:"name"`             // 用户名称
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  string `json:"total_favorited"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`       // 作品数
	Username        string `json:"-" gorm:"unique"`  // 注册用户名，最长32个字符
	PasswordHash    string `json:"-"`                // 密码，最长32个字符   service层完成对应的逻辑操作
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	u.Name = "用户" + strconv.FormatInt(u.ID, 10)
	u.Signature = "谢谢你的关注"
	return tx.Model(u).Updates(User{Name: u.Name, Signature: u.Signature}).Error
}

// Video
type Video struct {
	ID            int64  `json:"id" gorm:"primaryKey"`            // 视频唯一标识
	Author        User   `json:"author" gorm:"foreignKey:UserID"` // 视频作者信息
	Title         string `json:"title"`                           // 视频标题
	PlayURL       string `json:"play_url"`                        // 视频播放地址
	CoverURL      string `json:"cover_url"`                       // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"`                  // 视频的点赞总数
	CommentCount  int64  `json:"comment_count"`                   // 视频的评论总数
	IsFavorite    bool   `json:"is_favorite"`                     // true-已点赞，false-未点赞
	PublishedAt   int64  `json:"published_at" gorm:"index"`       // 视频发布时间
	UserID        int64  `json:"-"`                               // 视频作者id
}

// VideoLike
type VideoLike struct {
	ID      int64 `json:"id" gorm:"primaryKey"` // 视频点赞记录唯一标识
	UserID  int64 `json:"-" gorm:"index"`       // 点赞用户id
	VideoID int64 `json:"-" gorm:"index"`       // 被点赞视频id
}

// Comment
type Comment struct {
	ID         int64  `json:"id"`                            // 评论id
	Content    string `json:"content"`                       // 评论内容
	CreateDate string `json:"create_date" gorm:"index"`      // 评论发布日期，格式 mm-dd
	Video_id   int64  `json:"-" gorm:"index"`                // 评论视频id
	UserID     int64  `json:"-"`                             // 评论用户id,外键用于关联User表
	User       User   `json:"user" gorm:"foreignKey:UserID"` // 评论用户信息
}

// relation
type Relation struct {
	ID       int64 `json:"id" gorm:"primaryKey"`      // 关注记录唯一标识
	AuthorID int64 `json:"following_id" gorm:"index"` // 作者ID
	FansID   int64 `json:"follower_id" gorm:"index"`  // 粉丝ID
}

// message
type Message struct {
	ID         int64  `json:"id"`           // 消息id
	FromUserID int64  `json:"from_user_id"` // 消息发送者id
	ToUserID   int64  `json:"to_user_id"`   // 消息接收者id
	Content    string `json:"content"`      // 消息内容
	CreateTime int64  `json:"create_time"`  // 消息发送时间 yyyy-MM-dd HH:MM:ss
}
