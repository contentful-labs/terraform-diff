FROM golang:1.24-alpine as build

RUN apk add --no-cache ca-certificates git
WORKDIR $GOPATH/src/github.com/contentful-labs/terraform-diff

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -tags netgo -ldflags '-w' .

FROM alpine:3.22
RUN apk add --no-cache git
COPY --from=build /go/src/github.com/contentful-labs/terraform-diff/terraform-diff terraform-diff
RUN git config --global --add safe.directory '*'
ENTRYPOINT ["/terraform-diff"]
