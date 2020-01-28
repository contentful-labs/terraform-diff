FROM golang:1.13-alpine as build

RUN apk add --no-cache ca-certificates git
WORKDIR $GOPATH/src/github.com/contentful-labs/terraform-config-deps

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .

FROM scratch
COPY --from=build /go/src/github.com/contentful-labs/terraform-config-deps/terraform-config-deps /terraform-config-deps
ENTRYPOINT ["/terraform-config-deps"]
