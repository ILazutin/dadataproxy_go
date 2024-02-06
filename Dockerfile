FROM golang:1.21 as builder

WORKDIR /usr/src/app

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./dadataproxy ./cmd/dadataproxy


FROM alpine:3.18 as runner

WORKDIR /app

COPY --from=builder /usr/src/app/dadataproxy .

COPY ./config/prod.yaml ./
ENV CONFIG_PATH=./prod.yaml

EXPOSE 3001
CMD [ "/app/dadataproxy" ]
