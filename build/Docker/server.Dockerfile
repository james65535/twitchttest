FROM golang:1.14 AS build-env
COPY /cmd/server /go/src/app
COPY /vendor /go/src/app/vendor
WORKDIR /go/src/app
#RUN go install -i
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main /go/src/app/main.go

FROM scratch
# WORKDIR /app
COPY --from=build-env /go/src/app/main main
CMD ["/main"]
EXPOSE 8080