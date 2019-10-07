FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
RUN git config --global user.name "IntouchOpec"
RUN git config --global user.password "we.learn01"
COPY . .
RUN go get ./...
RUN GOOS=linux go build -o ./bin/web-app ./main.go

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 80
ENTRYPOINT /go/bin/web-app --port 80