package env

import (
	"os"
	"strings"
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

func GetSender() (sender string) {
	sender = os.Getenv("sender")
	return sender
}

func GetSenderPassword() (senderPassword string) {
	senderPassword = os.Getenv("sender_password")
	return senderPassword
}

func GetReceivers() (receivers []string) {
	receiversStr := os.Getenv("receivers")
	trimmedReceiver := strings.Replace(receiversStr, " ", "", -1)

	receivers = strings.Split(trimmedReceiver, ",")
	return receivers
}

func GetSmtpHost() (smtpHost string) {
	smtpHost = os.Getenv("smtp_host")
	return smtpHost
}

func GetSmtpPort() (smtpPort string) {
	smtpPort = os.Getenv("smtp_port")
	return smtpPort
}

//func GetKubeConfigPath() (kubeConfigPath string) {
//	kubeConfigPath = os.Getenv("kube_config_path")
//	return kubeConfigPath
//}
//
//func GetKubeConfig() (kubeConfig string) {
//	kubeConfig = os.Getenv("kube_config")
//	return kubeConfig
//}
//
//func GetKubeConfigContext() (kubeConfigContext string) {
//	kubeConfigContext = os.Getenv("kube_config_context")
//	return kubeConfigContext
//}
//
//func GetKubeConfigNamespace() (kubeConfigNamespace string) {
//	kubeConfigNamespace = os.Getenv("kube_config_namespace")
//	return kubeConfigNamespace
//}
//
//func GetKubeConfigInCluster() (kubeConfigInCluster string) {
//	kubeConfigInCluster = os.Getenv("kube_config_in_cluster")
//	return kubeConfigInCluster
//}
