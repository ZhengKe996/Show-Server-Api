package forms

type SendSmsForm struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required,mobile"`
	Type   uint   `json:"type" form:"type" binding:"required,oneof=1 2"` //1. 注册发送短信验证码和动态验证码登录发送验证码
}
