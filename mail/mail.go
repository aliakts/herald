package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"text/template"
	"time"

	"k8s.io/apimachinery/pkg/types"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func SendMail(level string, labels map[string]string, date time.Time, jobName string, namespace string, annotations map[string]string, uid types.UID, sender string, senderPassword string, receivers []string, smtpHost string, smtpPort string) {
	from := sender
	password := senderPassword

	to := receivers

	//auth := smtp.PlainAuth("", from, password, smtpHost)
	conn, err := net.Dial("tcp", smtpHost+":"+smtpPort)
	if err != nil {
		println(err)
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		println(err)
	}

	tlsconfig := &tls.Config{
		ServerName: smtpHost,
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		println(err)
	}

	auth := LoginAuth(from, password)

	if err = c.Auth(auth); err != nil {
		println(err)
	}

	if level == "failed" {
		t, err := template.ParseFiles("template/fail.html")
		if err != nil {
			log.Fatalf("Failed to parse template %v", err)
		}

		var body bytes.Buffer
		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf("Subject: [ALERT - Warning] Kubernetes Job Failure \n%s\n\n", mimeHeaders)))

		err = t.Execute(&body, struct {
			Labels      map[string]string
			Date        time.Time
			JobName     string
			Namespace   string
			Annotations map[string]string
			UID         types.UID
		}{
			Labels:      labels,
			Date:        date,
			JobName:     jobName,
			Namespace:   namespace,
			Annotations: annotations,
			UID:         uid,
		})
		if err != nil {
			log.Fatalf("Failed to execute template %v", err)
		}

		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
		if err != nil {
			log.Fatalf("Failed to send mail %v", err)
		}

		log.Println("Email sent.")
	} else {
		t, err := template.ParseFiles("template/success.html")
		if err != nil {
			log.Fatalf("Failed to parse template %v", err)
		}

		var body bytes.Buffer
		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf("Subject: [ALERT - Warning] Kubernetes Job Success \n%s\n\n", mimeHeaders)))

		err = t.Execute(&body, struct {
			Labels      map[string]string
			Date        time.Time
			JobName     string
			Namespace   string
			Annotations map[string]string
			UID         types.UID
		}{
			Labels:      labels,
			Date:        date,
			JobName:     jobName,
			Namespace:   namespace,
			Annotations: annotations,
			UID:         uid,
		})
		if err != nil {
			log.Fatalf("Failed to execute template %v", err)
		}

		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
		if err != nil {
			log.Fatalf("Failed to send mail %v", err)
		}

		log.Println("Email sent.")
	}
}
