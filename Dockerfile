FROM alpine:3.18
ENTRYPOINT ["/go-vxbot"]
COPY go-vxbot /
