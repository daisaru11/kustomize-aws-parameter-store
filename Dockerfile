FROM golang:1.13 as builder

WORKDIR /workspace

ARG KUSTOMIZE_VERSION=3.5.4
RUN curl -L -O "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_linux_amd64.tar.gz" \
    && tar xzf kustomize_v${KUSTOMIZE_VERSION}_linux_amd64.tar.gz \
    && chmod +x kustomize

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o AWSParameterStoreSecret main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/kustomize /usr/local/bin/kustomize
COPY --from=builder /workspace/AWSParameterStoreSecret /home/nonroot/.config/kustomize/plugin/kustomize.daisaru11.dev/v1/awsparameterstoresecret/AWSParameterStoreSecret

USER nonroot:nonroot
