package gitclient

const DEPLOYMENT_MANIFEST = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ide-%s
  namespace: ide
  labels:
    app: ide-%s
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ide-%s
  template:
    metadata:
      labels:
        app: ide-%s
    spec:
      containers:
        - name: ide-%s
          image: milkymilky0116/ide:%s
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
            periodSeconds: 20
`

const SERVICE_MANIFEST = `---
apiVersion: v1
kind: Service
metadata:
  name: ide-%s
  namespace: ide
spec:
  type: NodePort
  selector:
    app: ide-%s
  ports:
    - port: 8080
      targetPort: 8080
`

const INGRESS_ROUTE_MANIFEST = `---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: ide-%s
  namespace: ide
spec:
  entryPoints:
    - web
    - websecure
  routes:
    - match: Host(%s)
      kind: Rule
      services:
        - name: ide-%s
          port: 8080
  tls:
    certResolver: myresolver
`

const APPLICATION_MANIFEST = `---
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: ide-user%scourse%s
  namespace: argocd
  annotations:
    notifications.argoproj.io/subscribe.on-sync-succeeded: "true"
    notifications.argoproj.io/subscribe.on-sync-failed: "true"
spec:
  project: default
  source:
    repoURL: "https://github.com/9oormthon-univ/2024_DANPOON_TEAM_25_MANIFEST"
    path: "ide-user%scourse%s"
    targetRevision: main
  destination:
    server: "https://kubernetes.default.svc"
    namespace: ide
  syncPolicy:
    automated:
      prune: true # 더 이상 필요없는 리소스 자동 제거
      selfHeal: true # 직접 변경된 리소스 자동 복구
      allowEmpty: false # 빈 디렉토리 배포 방지
    syncOptions:
      - CreateNamespace=true # 네임스페이스 자동 생성
      - PrunePropagationPolicy=foreground # 리소스 삭제 시 종속성 고려
      - PruneLast=true # 새 리소스 생성 후 이전 리소스 제거
    retry:
      limit: 5 # 동기화 실패 시 재시도 횟수
      backoff:
        duration: 5s # 초기 대기 시간
        factor: 2 # 재시도 간격 증가 비율
`
