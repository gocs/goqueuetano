# Stage 1
FROM golang:alpine as stage

COPY . /app
WORKDIR /app

RUN go build -o server ./cmd/http/main.go

# Stage 2
FROM alpine

RUN apk update && apk upgrade

WORKDIR /usr/bin/goqueuetano/
COPY --from=stage /app/public ./public/
COPY --from=stage /app/server ./

EXPOSE 3000

ENTRYPOINT [ "./server" ]
