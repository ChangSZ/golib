package mail

import (
	"testing"
)

func TestSend(t *testing.T) {
	t.Skip("skip")
	client, err := Init(
		WithUser("xxx@163.com"),
		WithPwd("123456"),
		WithHost("smtp.163.com"),
		WithPort(465))
	client.SetTo([]string{"yyy@qq.com"}).
		SetSubject("测试").
		SetBody(`<html><body>
	<p><img src="https://golang.org/doc/gopher/doc.png"></p><br/>
	<h1>测试玩儿</h1>
	</body></html>`).
		SetAttaches(map[string]string{"文件1": "/a/b/c.log"}).
		Send()
	if err != nil {
		t.Error("Mail Send error", err)
		return
	}
	t.Log("success")
}
