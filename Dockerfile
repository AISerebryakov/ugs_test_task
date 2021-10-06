FROM golang:1.17.1-alpine as build
RUN apk update && apk add --no-cache
COPY . /src
WORKDIR /src
RUN go build -mod=vendor -o ugc_test_task_service main.go && chmod u+x ugc_test_task_service

FROM alpine:3.14.0
COPY --from=build /src/ugc_test_task_service .
CMD ["/ugc_test_task_service"]
