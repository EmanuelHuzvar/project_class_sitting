# Base image
FROM golang:latest

RUN mkdir /build
WORKDIR /build
RUN export GO111MODULE=on


# Build the application
RUN cd /go/api && go build

EXPOSE 9090

# Command to run the application
CMD ["/go/api/main/main"]