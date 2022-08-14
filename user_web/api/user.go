package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"server_api/user_web/forms"
	"server_api/user_web/global"
	"server_api/user_web/global/response"
	"server_api/user_web/middlewares"
	"server_api/user_web/models"
	"server_api/user_web/proto"
	"strconv"
	"strings"
	"time"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// HandleGrpcError2Http 将GRPC的Code转换成Http的状态码
func HandleGrpcError2Http(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"message": e.Message() + "❎"})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误❎"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误❎"})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "用户服务不可用❎"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "其他错误❎"})
			}
			return
		}
	}
}

// HandleValidatorError 处理表单验证的错误信息
func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
}

// GetUserList 获取用户列表
func GetUserList(ctx *gin.Context) {
	// 拨号连接用户 GRPC 服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
	}
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)
	// 生成grpc的client 调用接口
	userClient := proto.NewUserClient(userConn)

	page := ctx.DefaultQuery("page", "0")
	pSize := ctx.DefaultQuery("size", "10")
	fmt.Println("page", page, "size", pSize)
	pageInt, _ := strconv.Atoi(page)
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pageInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Error("[GetUserList] 【查询用户列表失败】")
		HandleGrpcError2Http(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		userResponse := response.UserResponse{
			ID:       value.Id,
			NikeName: value.NickName,
			Mobile:   value.Mobile,
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
		}

		result = append(result, userResponse)
	}

	ctx.JSON(http.StatusOK, result)
}

// PassWordLogin 用户登录
func PassWordLogin(ctx *gin.Context) {
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "验证码错误❎"})
		return
	}
	// 拨号连接用户 GRPC 服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[PassWordLogin] 连接 【用户服务失败】", "msg", err.Error())
	}

	// 生成grpc的client 调用接口
	userClient := proto.NewUserClient(userConn)

	// 登录逻辑
	if rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: passwordLoginForm.Mobile}); err != nil {
		//HandleGrpcError2Http(err, ctx)
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		if checkPasswordResult, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "密码错误❎"})
		} else {
			if checkPasswordResult.Success {
				// 生成 token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),
						ExpiresAt: time.Now().Unix() + 60*60*24*7, // 7天过期
						Issuer:    "timu.fun",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": "生成token失败❎"})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nike_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*7) * 1000,
				})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "密码错误❎"})
			}
		}
	}
}

// Register 用户注册
func Register(ctx *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 验证码校验
	redisDB := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	if value, err := redisDB.Get(context.Background(), registerForm.Mobile).Result(); err == redis.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "验证码错误❎"})
		return
	} else {
		if value != registerForm.Code {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "验证码错误❎"})
			return
		}
	}

	// 拨号连接用户 GRPC 服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
	}

	// 生成grpc的client 调用接口
	userClient := proto.NewUserClient(userConn)

	if user, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	}); err != nil {
		zap.S().Errorf("[Register] 【新建用户失败】原因: %s", err.Error())
		HandleGrpcError2Http(err, ctx)
		return
	} else {
		// 生成 token
		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			ID:          uint(user.Id),
			NickName:    user.NickName,
			AuthorityId: uint(user.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),
				ExpiresAt: time.Now().Unix() + 60*60*24*7, // 7天过期
				Issuer:    "timu.fun",
			},
		}
		token, err := j.CreateToken(claims)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "生成token失败❎"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"id":         user.Id,
			"nike_name":  user.NickName,
			"token":      token,
			"expired_at": (time.Now().Unix() + 60*60*24*7) * 1000,
		})
	}
}
