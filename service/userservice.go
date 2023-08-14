package service

import (
	"errors"
	"fmt"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/Coreychen4444/Lite_TikTok/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	r *repository.DbRepository
}

func NewUserService(r *repository.DbRepository) *UserService {
	return &UserService{r: r}
}

// 注册
func (s *UserService) Register(username, password string) (int64, string, error) {
	//校验
	if len(username) == 0 || len(password) == 0 {
		return -1, "", errors.New("用户名或密码不能为空,请重新输入")
	}
	if len(password) < 6 {
		return -1, "", errors.New("密码长度不能小于6位,请重新输入")
	}
	if len(username) > 32 || len(password) > 32 {
		return -1, "", errors.New("用户名或密码长度不能超过32位,请重新输入")
	}
	user, err := s.r.GetUserByName(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return -1, "", fmt.Errorf("查找用户时出错")
	}
	//判断用户名是否存在
	if user != nil {
		return -1, "", errors.New("用户名已存在,请重新输入")
	}
	var newUser model.User
	//加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, "", fmt.Errorf("设置的的密码格式有误")
	}
	newUser.Username = username
	newUser.PasswordHash = string(hashedPassword)
	//创建用户
	id, err := s.r.CreateUsers(&newUser)
	if err != nil {
		return -1, "", fmt.Errorf("创建用户时出错")
	}
	token, tknerr := GenerateToken(user.ID)
	if tknerr != nil {
		return -1, "", fmt.Errorf("生成token时出错")
	}
	return id, token, nil
}

// 登录
func (s *UserService) Login(username, password string) (int64, string, error) {
	//校验输入
	if len(username) == 0 || len(password) == 0 {
		return -1, "", fmt.Errorf("用户名或密码不能为空,请重新输入")
	}
	if len(username) > 32 || len(password) > 32 {
		return -1, "", fmt.Errorf("用户名或密码长度不能超过32位,请重新输入")
	}
	//查找用户
	user, err := s.r.GetUserByName(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return -1, "", fmt.Errorf("用户名不存在,请重新输入")
		}
		return -1, "", fmt.Errorf("查找用户时出错")
	}
	//验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return -1, "", fmt.Errorf("密码错误,请重新输入")
		}
		return -1, "", fmt.Errorf("验证密码时出错")
	}
	token, tknerr := GenerateToken(user.ID)
	if tknerr != nil {
		return -1, "", fmt.Errorf("生成token时出错")
	}
	return user.ID, token, nil
}

// 获取用户信息
func (s *UserService) GetUserInfo(id int64, token string) (*model.User, error) {
	clamis, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token无效,请重新登录")
	}
	user, err := s.r.GetUserById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("该用户不存在")
		}
		return nil, fmt.Errorf("查找用户时出错")
	}
	// 判断是否关注该用户
	isFollow, err := s.r.IsFollow(id, clamis.UserID)
	if err != nil {
		return nil, fmt.Errorf("查找用户时出错")
	}
	user.IsFollow = isFollow
	return user, nil
}

// 获取用户视频列表
func (s *UserService) GetUserVideoList(id int64, token string) ([]model.Video, error) {
	_, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token无效,请重新登录")
	}
	video, err := s.r.GetVideoListByUserId(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("该用户不存在")
		}
		return nil, fmt.Errorf("获取视频失败")
	}
	return video, nil
}
