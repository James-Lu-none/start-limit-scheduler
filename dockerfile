FROM golang:1.25rc2-alpine3.22 as builder

WORKDIR /app
COPY . .

RUN apk add git

RUN git clone -b release-1.32 --depth 1  --single-branch https://github.com/kubernetes/kubernetes ./kubernetes

RUN go mod init start-limit-scheduler && \
    go mod edit -replace=k8s.io/kubernetes=./kubernetes && \
    go get k8s.io/api@v0.32.0 && \
    go get k8s.io/component-base@v0.32.0 && \
    go get k8s.io/cloud-provider@v0.32.0 && \
    go get k8s.io/apiserver@v0.32.0 && \
    go get k8s.io/client-go@v0.32.0 && \
    go get k8s.io/apimachinery@v0.32.0 && \
    go get k8s.io/dynamic-resource-allocation@v0.32.0 && \
    go get k8s.io/kube-scheduler@v0.32.0 && \
    go get k8s.io/apiextensions-apiserver@v0.32.0 && \
    go get k8s.io/csi-translation-lib@v0.32.0 && \
    go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o start-limit-scheduler main.go

FROM gcr.io/distroless/static
COPY --from=builder /app/start-limit-scheduler /kube-scheduler
ENTRYPOINT ["/kube-scheduler"]
