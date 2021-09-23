
FROM golang:1.16.8 
COPY . /app/.
RUN apt-get update && apt-get install -y
RUN apt-get install libdlib-dev libblas-dev libatlas-base-dev liblapack-dev -y
#RUN apt-get install libjpeg-turbo8-dev -y
#RUN apt-get install libjpeg62-turbo-dev -y
RUN go get -d -v ./...

#ENV DB_URI_CLOUD=
ENV DB_NAME=db_goface_api_echo

CMD [ "go","run","/app/main.go" ]

















#build stage
# FROM golang:alpine AS builder
# RUN apk add --no-cache git
# WORKDIR /go/src/app
# COPY . .
# RUN go get -d -v ./...
# RUN go build -o /go/bin/app -v ./...

# #final stage
# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# COPY --from=builder /go/bin/app /app
# ENTRYPOINT /app
# LABEL Name=gofaceapi Version=0.0.1
# EXPOSE 3000