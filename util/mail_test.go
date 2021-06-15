package util

import (
	"fmt"
	"testing"
)

func TestMail(t *testing.T) {
	mail := &Mail{
		userName:   "r",
		passWord:   "",
		smtpServer: "smtp.qq.com",
		mailType:   "html",
		smtpPort:   465,
		ssl:        true,
	}
	body, err := ReadText("../debug_report.html")
	err = mail.Send("956992888@qq.com", "[12324]测试报告", body)
	fmt.Println(err)
}
