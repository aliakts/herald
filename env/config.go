package env

import (
	"os"
)

func GetNamespace() (namespace string) {
	namespace = os.Getenv("namespace")
	if namespace == "" {
		namespace = "default"
	}
	return namespace
}

func GetNotificationLevel() (level string) {
	level = os.Getenv("notification_level")
	if level == "" {
		level = "failed"
	}
	return level
}

func IsInCluster() bool {
	inCluster := os.Getenv("in_cluster")
	return inCluster == "1"
}
