package service

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/Coreychen4444/Lite_TikTok/repository"
)

type VideoService struct {
	r *repository.DbRepository
}

func NewVideoService(r *repository.DbRepository) *VideoService {
	return &VideoService{r: r}
}

// 获取视频列表
func (s *VideoService) GetVideoFlow(latest_time, token string) ([]model.Video, error) {
	latestTime, err := strconv.ParseInt(latest_time, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("时间戳格式错误")
	}
	video, err := s.r.GetVideoList(latestTime)
	if err != nil {
		return nil, fmt.Errorf("获取视频失败")
	}
	claims, err := VerifyToken(token)
	if err != nil {
		return video, nil
	}
	// 返回点赞个性化信息
	videoId, err := s.r.GetUserLikeId(claims.UserID)
	if err != nil {
		return video, nil
	}
	isLike := make(map[int64]bool)
	for _, id := range videoId {
		isLike[id] = true
	}
	for i := range video {
		video[i].IsFavorite = isLike[video[i].ID]
	}
	return video, nil
}

// 发布视频
func (s *VideoService) PublishVideo(fileHeader *multipart.FileHeader, title string, token string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("获取文件失败")
	}
	defer file.Close()
	// 校验token
	claims, err := VerifyToken(token)
	if err != nil {
		return fmt.Errorf("token无效, 请重新登录")
	}
	// 创建视频文件存放路径
	if _, err := os.Stat("../public/videofile"); os.IsNotExist(err) {
		err = os.Mkdir("../public/videofile", 0755)
		if err != nil {
			return err
		}
	}
	fileName := filepath.Base(fileHeader.Filename)
	// 生成唯一文件名
	uniqueFileName := generateRandomString(10) + fileName
	// 保存视频文件
	path := "../public/videofile/" + uniqueFileName
	destFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败")
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, file)
	if err != nil {
		return fmt.Errorf("保存文件失败")
	}
	// 制定视频封面
	uniqueNameWithoutExt := uniqueFileName[0 : len(uniqueFileName)-len(filepath.Ext(uniqueFileName))]
	// 创建封面图片存放路径
	if _, err := os.Stat("../pubilc/cover"); os.IsNotExist(err) {
		err = os.Mkdir("../public/cover", 0755)
		if err != nil {
			return err
		}
	}
	coverPath := "../public/cover/" + uniqueNameWithoutExt + ".jpg"
	// 生成封面图片
	err = generateCoverFromVideo(path, coverPath)
	if err != nil {
		return fmt.Errorf("生成封面失败")
	}
	// 创建视频
	video := &model.Video{
		Title:       title,
		UserID:      claims.UserID,
		PlayURL:     fmt.Sprintf("/public/videofile/%s", uniqueFileName),
		CoverURL:    fmt.Sprintf("/public/cover/%s.jpg", uniqueNameWithoutExt),
		PublishedAt: time.Now().UTC().Unix(),
	}
	err = s.r.CreateVideo(video)
	if err != nil {
		return fmt.Errorf("发布视频失败")
	}
	return nil
}

// 生成随机文件名
func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// 生成视频封面
func generateCoverFromVideo(videoPath string, coverPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", videoPath, // 输入文件路径
		"-ss", "00:00:01", // 开始时间，这里设置为视频的第1秒
		"-vframes", "1", // 只输出1帧图片
		"-f", "image2", // 输出格式
		coverPath, // 输出文件路径
	)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate cover: %w", err)
	}
	return nil
}

// 点赞和取消点赞视频
func (s *VideoService) LikeVideo(token, video_id, action_type string) error {
	// 校验token
	claims, err := VerifyToken(token)
	if err != nil {
		return fmt.Errorf("token无效, 请重新登录")
	}
	// 获取视频id
	id, err := strconv.Atoi(video_id)
	if err != nil {
		return fmt.Errorf("视频出错")
	}
	// 根据action_type执行点赞或取消点赞操作
	if action_type == "1" {
		err = s.r.LikeVideo(claims.UserID, int64(id))
		if err != nil {
			return fmt.Errorf("点赞视频失败")
		}
	} else if action_type == "2" {
		err = s.r.DislikeVideo(claims.UserID, int64(id))
		if err != nil {
			return fmt.Errorf("取消点赞失败")
		}
	} else {
		return fmt.Errorf("请求参数错误")
	}
	return nil
}

// 获取用户点赞的视频列表
func (s *VideoService) GetUserLike(token, user_id string) ([]model.Video, error) {
	// 校验token
	_, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token无效, 请重新登录")
	}
	// 获取用户id
	id, err := strconv.Atoi(user_id)
	if err != nil {
		return nil, fmt.Errorf("用户id格式错误")
	}
	// 获取用户点赞的视频列表
	videos, err := s.r.GetUserLike(int64(id))
	if err != nil {
		return nil, fmt.Errorf("获取列表失败")
	}
	return videos, nil
}

// 评论视频
func (s *VideoService) CommentVideo(token, video_id string, content *string) (*model.Comment, error) {
	// 校验token
	claims, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token无效, 请重新登录")
	}
	// 获取视频id
	id, err := strconv.Atoi(video_id)
	if err != nil {
		return nil, fmt.Errorf("视频出错")
	}
	// 创建评论
	comment := &model.Comment{
		Video_id:   int64(id),
		UserID:     claims.UserID,
		Content:    *content,
		CreateDate: time.Now().Format("01-02"),
	}
	comment, err = s.r.CreateComment(comment)
	if err != nil {
		if err.Error() == "评论发表成功, 但返回评论信息失败" {
			return nil, err
		}
		return nil, fmt.Errorf("评论失败")
	}
	return comment, nil
}

// 删除评论
func (s *VideoService) DeleteComment(token, comment_id string) error {
	// 校验token
	_, err := VerifyToken(token)
	if err != nil {
		return fmt.Errorf("token无效, 请重新登录")
	}
	// 获取评论id
	id, err := strconv.Atoi(comment_id)
	if err != nil {
		return fmt.Errorf("删除评论出错")
	}
	// 删除评论
	err = s.r.DeleteComment(int64(id))
	if err != nil {
		return fmt.Errorf("删除评论失败")
	}
	return nil
}

// 获取视频评论列表
func (s *VideoService) GetCommentList(token, video_id string) ([]model.Comment, error) {
	// 校验token
	_, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token无效, 请重新登录")
	}
	// 获取视频id
	id, err := strconv.Atoi(video_id)
	if err != nil {
		return nil, fmt.Errorf("视频出错")
	}
	// 获取视频评论列表
	comments, err := s.r.GetCommentList(int64(id))
	if err != nil {
		return nil, fmt.Errorf("获取评论列表失败")
	}
	return comments, nil
}
