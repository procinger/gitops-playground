####### demo-app builder
FROM docker.io/golang:1.25.1-alpine3.22 AS builder

WORKDIR /

ENV USER=app
ENV UID=1001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

COPY go.mod go.mod
COPY main.go  main.go

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="" -o backend

####### demo-app
FROM scratch

WORKDIR /

USER app:app

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /backend /backend

ENTRYPOINT ["/backend"]
