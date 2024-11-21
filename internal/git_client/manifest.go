package gitclient

const DEPLOYMENT_MANIFEST = `---
apiVersion: apps/v1
kind: Deployment
metadata:
	name: ide:%s
  namespace: ide
  labels:
		app: ide:%s
spec:
  replicas: 1
  selector:
    matchLabels:
			app: ide:%s
  template:
    metadata:
      labels:
				app: ide:%s
    spec:
      containers:
				- name: ide:%s
          image: %s
          ports:
            - containerPort: 8080
          resources:
            requests:
              cpu: 100m
              memory: 1Gi
            limits:
              cpu: 250m
              memory: 2Gi
          readinessProbe:
            httpGet:
              path: /
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /
              port: 8080
            initialDelaySeconds: 15
            periodSeconds: 20`

const SERVICE_MANIFEST = `
apiVersion: v1
kind: Service
metadata:
	name: ide:%s
  namespace: ide
spec:
  type: NodePort
  selector:
    app: ide%s
  ports:
    - port: 8080
      targetPort: 8080
	`
