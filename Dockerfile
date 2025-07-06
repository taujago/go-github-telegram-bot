FROM golang:1.22.5 AS builder 


# Create a new directory for the application and change its owner to 'nonroot'
RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy -v
COPY . /app 
RUN CGO_ENABLED=0 GOOS=linux go build -o letsgo cmd/server/main.go


FROM alpine
ARG USER=nonroot
ARG GROUP=nonroot
ENV HOME /app
RUN apk add --update sudo
RUN adduser -D $USER -u 1001    \
    && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
    && chmod 0440 /etc/sudoers.d/$USER

RUN apk add --no-cache tzdata

ENV TZ=Asia/Singapore \
    IMAGE=true
WORKDIR /app

COPY --chown=${USER}:${GROUP} --from=builder /app/letsgo .

USER nonroot
CMD ["./letsgo"]
