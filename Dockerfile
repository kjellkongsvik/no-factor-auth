FROM golang as base

RUN update-ca-certificates

WORKDIR /code

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .

ENV CGO_ENABLED=0
RUN go test ./...

RUN go build

FROM scratch

COPY --from=base /code/no-factor-auth .

EXPOSE 8089
CMD ["./no-factor-auth"]
