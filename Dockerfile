ARG GOLANG_VERSION=1.14.4-alpine3.12

FROM golang:${GOLANG_VERSION}
# Force go compiler to use modules with vendor mode
ENV GOFLAGS -mod=vendor
ENV GO111MODULE=on
RUN mkdir -m 777 /.cache
RUN apk --no-cache add build-base
WORKDIR /app
COPY ./ /app
RUN go install ./cmd/s3url-presigner
CMD ["s3url-presigner", "--help"]