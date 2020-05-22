FROM golang as base

RUN update-ca-certificates

WORKDIR /code

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .

ENV CGO_ENABLED=0
ENV GOARCH=amd64
RUN go test ./...

RUN GOOS=linux go build -ldflags "-X main.Version=$TAG" -o no-factor-auth.linux

FROM scratch

ARG AUTHSERVER
ARG TENANT_ID

COPY --from=base /code/no-factor-auth.linux no-factor-auth

EXPOSE 8089
CMD ["./no-factor-auth"]
