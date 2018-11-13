FROM golang
WORKDIR /go/src/github.com/the4thamigo-uk/ingestor
ADD . .
RUN mkdir bin
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ./bin/ingestor ./cmd/ingestor
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ./bin/cli ./cmd/cli
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ./bin/reader ./cmd/reader

FROM alpine
COPY --from=0 /go/src/github.com/the4thamigo-uk/ingestor/bin/* /
CMD ["/ingestor"]
