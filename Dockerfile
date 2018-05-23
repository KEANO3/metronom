FROM golang:1 AS builder
ADD . /go/src/metronom/password
WORKDIR /go/src/metronom/password
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -a -installsuffix cgo -o /thisbuild metronom/password/cmd/rest-server

FROM scratch
COPY --from=builder /thisbuild /rest-server
ENTRYPOINT ["/rest-server"]