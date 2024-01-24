package main

import (
	"herald/env"
	"herald/mail"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	k8s "herald/kubernetes"
)

func main() {
	var namespace string
	var level string

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
						mail.SendMail(level, job.Labels, CompletionTime, job.Name, job.Namespace, job.Annotations, job.UID)
						pastJobs[jobUniqueHash] = true
					}
				}

				if level == "failed" {
					if job.Status.Failed > 0 {
						log.Println("An unsuccessful job was discovered. I'm sending email to people right now.")
						if job.Status.StartTime.Add(5*time.Hour).Unix() > time.Now().Unix() {
							FailureTime := job.Status.StartTime.Time
							mail.SendMail(level, job.Labels, FailureTime, job.Name, job.Namespace, job.Annotations, job.UID)
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
