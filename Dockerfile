FROM golang:1.18-alpine as build-stage

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/Hlebushkek/genesis_task

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /genesis_task .

#
# final build stage
#
FROM scratch

# Copy ca-certs for app web access
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-stage /genesis_task /genesis_task

# app uses port 3000
EXPOSE 3000

ENTRYPOINT ["/genesis_task"]
