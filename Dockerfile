FROM golang:1.17

RUN go install github.com/jstemmer/go-junit-report@latest

ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/cyverse-de/QMS
COPY . .
RUN make

FROM scratch

WORKDIR /app

COPY --from=0 /go/src/github.com/cyverse-de/QMS/QMS /bin/QMS
COPY --from=0 /go/src/github.com/cyverse-de/QMS/swagger.json swagger.json

ENTRYPOINT ["QMS"]

EXPOSE 8080
