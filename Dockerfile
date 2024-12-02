# Builder image
FROM golang:1.23.3-alpine3.19 as builder
RUN echo "--- Start build builder image ---"
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN echo "Executing tests"
RUN CGO_ENABLED=0 GOOS=linux go test ./...
RUN echo "Building binary"
RUN CGO_ENABLED=0 GOOS=linux go mod tidy && go build -a -o main cmd/challenge/main.go
RUN echo "--- End build builder image ---"
# Runtime image
RUN echo "--- Start build running image ---"
FROM alpine:3.16.2
COPY --from=builder /build/main .
ADD template /template 

ENTRYPOINT [ "./main" ]
RUN echo "--- End build running image ---"