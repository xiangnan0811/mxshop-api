package forms

type PassWordLoginForm struct {
    Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"`
    Password  string `form:"password" json:"password" binding:"required,min=6,max=20"`
    Captcha   string `form:"captcha" json:"captcha" binding:"required,len=5"`
    CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required,len=20"`
}

type RegisterForm struct {
    Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"`
    Password  string `form:"password" json:"password" binding:"required,min=6,max=20"`
    Code      string `form:"code" json:"code" binding:"required,len=6"`
    Captcha   string `form:"captcha" json:"captcha" binding:"required,len=5"`
    CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required,len=20"`
}
