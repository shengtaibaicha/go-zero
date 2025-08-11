package myCaptcha

import (
	"bytes"
	"encoding/base64"
	"github.com/dchest/captcha"
	"image/png"
)

// GetCaptcha 获取验证码
func GetCaptcha() (resp string, code string) {
	// 第一种验证码
	//captchaId := captcha.New()
	//var buf bytes.Buffer
	//err := captcha.WriteImage(&buf, captchaId, 240, 80)
	//if err != nil {
	//	resp = result.Err()
	//}
	//
	//// 将图片数据编码为Base64
	//v := base64.StdEncoding.EncodeToString(buf.Bytes())
	//
	//// 构建带格式的Base64字符串，前端可直接用于img标签的src属性
	//base64Data := "data:image/png;base64," + imgBase64

	// 第二种验证码
	// 获取六位数验证码
	digits := captcha.RandomDigits(6)
	// 创建buffer缓冲
	buf := new(bytes.Buffer)
	// 用验证码生成图片
	image := captcha.NewImage("", digits, 150, 60)
	// 将图片写入缓冲区
	png.Encode(buf, image)
	// 对缓冲区进行base64编码
	imageStr := base64.StdEncoding.EncodeToString(buf.Bytes())
	resp = "data:image/png;base64," + imageStr

	// 将[]byte类型的数据转换为string
	buff := make([]byte, len(digits))
	for i, d := range digits {
		buff[i] = d + 48
	}
	code = string(buff)
	return
}

// VerifyCaptcha 检查验证码是否正确
func VerifyCaptcha(captchaId string, code string) (r bool) {
	r = captcha.VerifyString(captchaId, code)
	return
}
