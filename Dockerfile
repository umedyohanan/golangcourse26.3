FROM bitnami/golang:latest AS builder
RUN mkdir -p /go/src/project
WORKDIR /go/src/project
ADD main.go .
ADD pipelinehelper ./pipelinehelper
ADD ringbuffer ./ringbuffer
ADD go.mod .
RUN go build -o pipeline

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/project/pipeline .
ENTRYPOINT ./pipeline