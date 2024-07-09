package mail

import (
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailParam struct {
	user     string            // 发件人
	password string            // 授权码
	host     string            // 主机地址, QQ: smtp.qq.com
	port     int               // 端口, QQ: 587
	to       []string          // 接收人
	cc       []string          // 抄送
	bcc      []string          // 密送
	subject  string            // 主题
	body     string            // 内容
	mailType string            // 邮件类型
	attaches map[string]string // 附件
}

var mailParam *EmailParam

type EM func(*EmailParam) error

func WithUser(user string) EM {
	return func(e *EmailParam) error {
		if user == "" {
			return errors.New("user can not be null")
		}
		e.user = user
		return nil
	}
}

func WithPwd(pwd string) EM {
	return func(ep *EmailParam) error {
		if pwd == "" {
			return errors.New("pwd can not be null")
		}
		ep.password = pwd
		return nil
	}
}

func WithHost(host string) EM {
	return func(ep *EmailParam) error {
		if host == "" {
			return errors.New("host can not be null")
		}
		ep.host = host
		return nil
	}
}

func WithPort(port int) EM {
	return func(ep *EmailParam) error {
		if port == 0 {
			return errors.New("port can not be null")
		}
		ep.port = port
		return nil
	}
}

func WithMailType(mailType string) EM {
	return func(ep *EmailParam) error {
		if mailType == "" {
			return errors.New("mailType can not be null")
		}
		ep.mailType = mailType
		return nil
	}
}

func Init(options ...EM) (*EmailParam, error) {
	q := &EmailParam{
		mailType: "html",
	}
	for _, option := range options {
		err := option(q)
		if err != nil {
			return nil, err
		}
	}
	mailParam = q
	if q.host == "" {
		q.host = "smtp.qq.com"
		fmt.Println("未设置主机地址, 默认使用QQ邮箱: ", q.host)
	}
	if q.port == 0 {
		q.port = 587
		fmt.Println("未设置端口, 默认使用QQ邮箱: ", q.port)
	}
	if q.user == "" {
		return nil, fmt.Errorf("请设置发件人")
	}
	if q.password == "" {
		return nil, fmt.Errorf("请设置授权码")
	}
	return q, nil
}

func (ep *EmailParam) SetSubject(s string) *EmailParam {
	ep.subject = s
	return ep
}

func (ep *EmailParam) SetAttaches(a map[string]string) *EmailParam {
	ep.attaches = a
	return ep
}

func (ep *EmailParam) SetBody(b string) *EmailParam {
	ep.body = b
	return ep
}

func (ep *EmailParam) SetTo(to []string) *EmailParam {
	ep.to = to
	return ep
}

func (ep *EmailParam) SetCc(cc []string) *EmailParam {
	ep.cc = cc
	return ep
}

func (ep *EmailParam) SetBcc(bcc []string) *EmailParam {
	ep.bcc = bcc
	return ep
}

func (ep *EmailParam) Send() error {
	m := gomail.NewMessage()
	// 发送人
	m.SetHeader("From", ep.user)
	// 接收人
	m.SetHeader("To", ep.to...)
	// 抄送
	if len(ep.cc) > 0 {
		m.SetHeader("Cc", ep.cc...)
	}
	// 密送
	if len(ep.bcc) > 0 {
		m.SetHeader("Bcc", ep.bcc...)
	}
	// 主题
	m.SetHeader("Subject", ep.subject)
	// 内容
	m.SetBody("text/html", ep.body)
	// 附件
	for _, attaFile := range ep.attaches {
		m.Attach(attaFile)
	}

	// 拿到token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer(ep.host, ep.port, ep.user, ep.password)
	// 发送邮件
	return d.DialAndSend(m)
}

// Send 简易的邮件发送, 比如告警什么的, 直接发送
func Send(to []string, subject string, body string) error {
	user := string(mailParam.user)
	password := string(mailParam.password)
	host := string(mailParam.host)
	port, _ := strconv.Atoi(string(mailParam.port))

	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", "<html><body>"+body+"</body></html>")
	d := gomail.NewDialer(host, port, user, password)
	return d.DialAndSend(m)
}
