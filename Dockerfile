FROM golang:1.15-alpine as build
ENV GO111MODULE=on

WORKDIR $GOPATH/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd cmd

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app cmd/main.go

FROM alpine:3.9
ENV CS_API=https://api
EXPOSE 8811
COPY --from=build /app /bin
CMD ["/bin/app"]
