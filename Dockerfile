FROM golang:1.18-alpine as build-stage

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/swayne275/joke-web-server

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /GenesisTask .

#
# final build stage
#
FROM scratch

# Copy ca-certs for app web access
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-stage /GenesisTask /GenesisTask

# app uses port 3000
EXPOSE 3000

ENTRYPOINT ["/GenesisTask"]