FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /proxyd ./cmd/proxyd

FROM gcr.io/distroless/static-debian12
COPY --from=build /proxyd /proxyd
EXPOSE 8232
ENTRYPOINT ["/proxyd"]
