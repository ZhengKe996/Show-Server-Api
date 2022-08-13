package forms

type PassWordLoginForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"`
	PassWord string `json:"password" form:"password" binding:"required,min=3,max=10"`
}
