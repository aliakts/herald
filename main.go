package main

import (
	"herald/env"
	"herald/mail"
	"log"
	"os"
	"time"

	k8s "herald/kubernetes"

	"go.uber.org/zap"
)

func main() {
	var namespace string
	var level string
	var sender string
	var senderPassword string
	var receivers []string
	var smtpHost string
	var smtpPort string

	//flag.StringVar(&namespace, "namespace", "default", "set namespace (optional)")
	//flag.StringVar(&level, "notification_level", "", "set notification level (required) (options: succeeded, failed)")
	//flag.Parse()

	pastJobs := make(map[string]bool)

	client, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client %v", zap.Error(err))
		os.Exit(1)
	}

	namespace = env.GetNamespace()
	log.Printf("Fetching jobs from %s namespace", namespace)

	level = env.GetNotificationLevel()
	log.Printf("Notification level set at %s", level)

	sender = env.GetSender()
	log.Printf("Sender email is %s", sender)

	senderPassword = env.GetSenderPassword()
	log.Printf("Sender password set")

	receivers = env.GetReceivers()
	log.Printf("Receivers are %v", receivers)

	smtpHost = env.GetSmtpHost()
	log.Printf("SMTP host is %s", smtpHost)

	smtpPort = env.GetSmtpPort()
	log.Printf("SMTP port is %s", smtpPort)

	for {
		jobs, err := client.ListJobs(namespace)

		if err != nil {
			log.Fatalf("Failed to list jobs in the namespace %v", zap.Error(err))
			continue
		}

		for _, job := range jobs.Items {
			jobUniqueHash := job.Name + job.CreationTimestamp.String()

			if pastJobs[jobUniqueHash] == false {
				if level == "succeeded" {
					if job.Status.Succeeded > 0 && (job.Status.CompletionTime.Add(20*time.Minute).Unix() > time.Now().Unix()) {
						log.Println("A successful job was discovered. I'm sending email to people right now.")
						CompletionTime := job.Status.CompletionTime.Time
						mail.SendMail(level, job.Labels, CompletionTime, job.Name, job.Namespace, job.Annotations, job.UID, sender, senderPassword, receivers, smtpHost, smtpPort)
						pastJobs[jobUniqueHash] = true
					}
				}

				if level == "failed" {
					if job.Status.Failed > 0 {
						log.Println("An unsuccessful job was discovered. I'm sending email to people right now.")
						if job.Status.StartTime.Add(5*time.Hour).Unix() > time.Now().Unix() {
							FailureTime := job.Status.StartTime.Time
							mail.SendMail(level, job.Labels, FailureTime, job.Name, job.Namespace, job.Annotations, job.UID, sender, senderPassword, receivers, smtpHost, smtpPort)
							pastJobs[jobUniqueHash] = true
						}
					}
				}
			}
		}

		log.Println("Fetching jobs...")
		time.Sleep(time.Minute * 1)
	}
}
