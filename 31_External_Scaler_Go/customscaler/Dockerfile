FROM golang:1.25 as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o external-scaler .

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /src/external-scaler .
ENTRYPOINT ["/external-scaler"]
