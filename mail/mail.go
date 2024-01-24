package mail

import (
	"bytes"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"net/smtp"
	"text/template"
	"time"
)

func SendMail(level string, labels map[string]string, date time.Time, jobName string, namespace string, annotations map[string]string, uid types.UID) {
	from := "sender@example.com"
	password := "password"

	to := []string{
		"receiver@example.com",
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)

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
