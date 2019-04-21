# BUILDER
FROM golang:1.12-alpine as builder
RUN mkdir /build
# The value of the Docker context must have a project root directory.
ADD . /build
WORKDIR /build
RUN cd cmd && go build -o bodis .

# App Image
FROM alpine
COPY --from=builder /build/cmd/bodis /usr/bin/bodis
CMD ["bodis"]