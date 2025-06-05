
FROM golang:1.24.4 AS build

ADD go.mod /go/
ADD go.sum /go/

RUN go mod download

ADD builders/ /go/builders/
ADD cmd/ /go/cmd
ADD graphql/ /go/graphql/
ADD internal/ /go/internal/
ADD misc/ /go/misc/

RUN go build ./cmd/main.go

FROM golang:1.24.4

COPY --from=build /go/main /main

ENTRYPOINT ["/main"]

