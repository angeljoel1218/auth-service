FROM golang:1.16-alpine AS build
LABEL MAINTAINER = 'Auth (angeljoel)'
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /auth-service -ldflags="-w -s"

FROM scratch
WORKDIR /
COPY --from=build /auth-service /auth-service
CMD ["/auth-service"]
