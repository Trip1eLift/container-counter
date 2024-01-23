FROM golang:1.21-alpine3.19

RUN apk add curl
# RUN apk add --update nodejs npm
# RUN npm install -g nodemon

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

EXPOSE 7001

HEALTHCHECK CMD curl -sf http://localhost:7001/health || exit 1 --health-retries=3 --health-timeout=3s --health-interval=5s --health-start-period=5s