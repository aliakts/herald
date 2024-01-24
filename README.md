# Herald

This tool sends an e-mail alert when there is a Kubernetes Job failure/success.

No extra setup required to deploy this tool on to your cluster, just apply below Kubernetes deployment manifest.

Uses `InClusterConfig` to access the Kubernetes API.

# Limitations

Namespace scoped i.e., each namespace should have this deploy separately

# Development

**Notification Levels**

* succeeded
* failed

```bash
export namespace="foo" && export notification_level="failed" && export in_cluster="0" && go build -o herald && ./herald
```

# Installation

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: foo
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: herald
  namespace: foo
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: foo
  name: job-reader
rules:
  - apiGroups: ['batch']
    resources:
      - jobs
    verbs:
      - get
      - list
      - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: herald
  namespace: foo
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: job-reader
subjects:
  - kind: ServiceAccount
    name: herald
    namespace: foo
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: herald
  namespace: foo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: herald
  template:
    metadata:
      labels:
        app: herald
    spec:
      serviceAccountName: herald
      containers:
        - name: herald
          image: aliaktas/herald:amd64
          imagePullPolicy: Always
          env:
            - name: namespace
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: notification_level
              value: "failed"
            - name: in_cluster
              value: "1"
          resources:
            limits:
              cpu: 500m
              memory: 256Mi
            requests:
              cpu: 500m
              memory: 128Mi
```