# The first stage: compile watcher
FROM golang:1.13-alpine as dependency-builder
ENV APP_NAME watcher
WORKDIR /opt/${APP_NAME}
COPY . .
RUN go mod download

# The second stage:
FROM dependency-builder as app-builder
ENV APP_NAME watcher
WORKDIR /opt/${APP_NAME}
COPY --from=dependency-builder /opt/watcher .
RUN CGO_ENABLED=0 go build -o ./bin/watcher .

# The third stage: copy the watcher binary to another container
FROM alpine:3.9
LABEL name="watcher" maintainer="o.kaya" version="0.1"
WORKDIR /opt/watcher
COPY --from=app-builder /opt/watcher/bin/watcher ./bin/
COPY --from=app-builder /opt/watcher/configs/config.json ./configs/
CMD ["./bin/watcher", "-c", "./configs/config.json"]
