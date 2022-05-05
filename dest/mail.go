package dest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

type Mail struct {
	Auth smtp.Auth
	Addr string
	From mail.Address
}

func NewMail(host, password, smtpserver string, port int) *Mail {
	return &Mail{
		Auth: smtp.PlainAuth("", host, password, smtpserver),
		Addr: fmt.Sprintf("%s:%d", smtpserver, port),
		From: mail.Address{Name: "", Address: host},
	}
}

func (m *Mail) writeString(b *bytes.Buffer, s string) *bytes.Buffer {
	_, err := b.WriteString(s)
	if err != nil {
		fmt.Print(err.Error())
	}
	return b
}

func (m *Mail) encodeSubject(subject string) string {
	b := bytes.NewBuffer([]byte(""))
	strs := []string{}
	length := 13
	for k, c := range strings.Split(subject, "") {
		b.WriteString(c)
		if k%length == length-1 {
			strs = append(strs, b.String())
			b.Reset()
		}
	}
	if b.Len() > 0 {
		strs = append(strs, b.String())
	}
	// MIME エンコードする
	b2 := bytes.NewBuffer([]byte(""))
	b2.WriteString("Subject:")
	for _, line := range strs {
		b2.WriteString(" =?utf-8?B?")
		b2.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
		b2.WriteString("?=\r\n")
	}
	return b2.String()
}

func (m *Mail) encodeBody(body string) string {
	b := bytes.NewBufferString(body)
	s := base64.StdEncoding.EncodeToString(b.Bytes())
	b2 := bytes.NewBuffer([]byte(""))
	for k, c := range strings.Split(s, "") {
		b2.WriteString(c)
		if k%76 == 75 {
			b2.WriteString("\r\n")
		}
	}
	return b2.String()

}

func (m *Mail) Send(to string, subject string, body string) (err error) {

	// msg := bytes.NewBuffer([]byte(""))
	// msg = m.writeString(msg, "From: "+m.From.String()+"\r\n")
	// msg = m.writeString(msg, "To: "+to+"\r\n")
	// msg = m.writeString(msg, "Bcc: "+m.From.String()+"\r\n")
	// msg = m.writeString(msg, m.encodeSubject(subject))
	// msg = m.writeString(msg, "MIME-Version: 1.0\r\n")
	// msg = m.writeString(msg, "Content-Type: text/plain; charset=\"utf-8\"\r\n")
	// msg = m.writeString(msg, "Content-Transfer-Encoding: base64\r\n")
	// msg = m.writeString(msg, "\r\n")

	// msg = m.writeString(msg, m.encodeBody(body))
	b := []byte(fmt.Sprintf("subject: %s\r\n\r\n%s\r\n", m.encodeSubject(subject), body))

	er := smtp.SendMail(m.Addr, m.Auth, m.From.Address, []string{to}, b)
	if er != nil {
		panic(er)
	}
	// tlsconfig := &tls.Config{
	// 	InsecureSkipVerify: true,
	// 	ServerName:         m.From.Address,
	// }
	// fmt.Println(m.From.Address)
	// fmt.Println(m.Addr)
	// fmt.Println(m.Auth)
	// conn, err := tls.Dial("tcp", m.Addr, tlsconfig)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// c, err := smtp.NewClient(conn, m.From.Address)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// if err = c.Auth(m.Auth); err != nil {
	// 	fmt.Println(m.Auth)
	// 	log.Panic(err)
	// }
	// if err = c.Mail(m.From.Address); err != nil {
	// 	log.Panic(err)
	// }
	// if err = c.Rcpt(to); err != nil {
	// 	log.Panic(err)
	// }

	// w, err := c.Data()
	// if err != nil {
	// 	log.Panic(err)
	// }
	// _, err = w.Write(msg.Bytes())
	// if err != nil {
	// 	log.Panic(err)
	// }
	// err = w.Close()
	// if err != nil {
	// 	log.Panic(err)
	// }

	// c.Quit()

	// fmt.Print(msg, "\n")
	// fmt.Print(body, "\n")

	return nil
}
