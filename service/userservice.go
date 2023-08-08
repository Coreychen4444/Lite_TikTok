package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/Coreychen4444/Lite_TikTok/repository"
	"github.com/dgrijalva/jwt-go"
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
func (s *UserService) Register(username, password string) (*model.User, string, error) {
	//校验
	if len(username) == 0 || len(password) == 0 {
		return nil, "", errors.New("用户名或密码不能为空,请重新输入")
	}
	if len(password) < 6 {
		return nil, "", errors.New("密码长度不能小于6位,请重新输入")
	}
	if len(username) > 32 || len(password) > 32 {
		return nil, "", errors.New("用户名或密码长度不能超过32位,请重新输入")
	}
	user, err := s.r.GetUserByName(username)
	if err != nil {
		return nil, "", err
	}
	//判断用户名是否存在
	if user != nil {
		return nil, "", errors.New("用户名已存在,请重新输入")
	}
	var newUser model.User
	//加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	newUser.Username = username
	newUser.PasswordHash = string(hashedPassword)
	//创建用户
	user, err = s.r.CreateUsers(&newUser)
	if err != nil {
		return nil, "", err
	}
	token, tknerr := GenerateToken(user.ID)
	if tknerr != nil {
		return nil, "", fmt.Errorf("生成token时出错: %w", tknerr)
	}
	return user, token, nil
}

// 登录
func (s *UserService) Login(username, password string) (*model.User, string, error) {
	//校验输入
	if len(username) == 0 || len(password) == 0 {
		return nil, "", fmt.Errorf("用户名或密码不能为空,请重新输入")
	}
	if len(username) > 32 || len(password) > 32 {
		return nil, "", fmt.Errorf("用户名或密码长度不能超过32位,请重新输入")
	}
	//查找用户
	user, err := s.r.GetUserByName(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", fmt.Errorf("用户名不存在,请重新输入")
		}
		return nil, "", fmt.Errorf("查找用户时出错: %w", err)
	}
	//验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, "", fmt.Errorf("密码错误,请重新输入")
		}
		return nil, "", fmt.Errorf("验证密码时出错: %w", err)
	}
	token, tknerr := GenerateToken(user.ID)
	if tknerr != nil {
		return nil, "", fmt.Errorf("生成token时出错: %w", tknerr)
	}
	return user, token, nil
}

// 获取用户信息
func (s *UserService) GetUserInfo(id int64, token string) (*model.User, error) {
	_, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("验证token时出错: %w", err)
	}
	user, err := s.r.GetUserById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("该用户不存在")
		}
		return nil, fmt.Errorf("查找用户时出错: %w", err)
	}
	return user, nil

}

// 生成和验证token
type Claims struct {
	UserID int64
	jwt.StandardClaims
}

var jwtKey = []byte("tokenkey")

// 生成token
func GenerateToken(id int64) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// 验证token
func VerifyToken(tknStr string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("invalid JWT signature")
		}
		return nil, fmt.Errorf("could not parse JWT token")
	}

	if !tkn.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}
