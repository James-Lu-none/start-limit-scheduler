FROM golang:1.25rc2-alpine3.22 as builder

WORKDIR /app
COPY . .

RUN apk add --no-cache git

RUN go mod init kube-scheduler-custom && \
    go get k8s.io/client-go@v0.32.0 && \
    go get k8s.io/component-base@v0.32.0 && \
    go get k8s.io/api@v0.32.0 && \
    go get k8s.io/apimachinery@v0.32.0 && \
    go get k8s.io/kube-scheduler@v0.32.0 && \
    go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o kube-scheduler-custom main.go

FROM gcr.io/distroless/static
COPY --from=builder /app/kube-scheduler-custom /kube-scheduler
ENTRYPOINT ["/kube-scheduler"]
