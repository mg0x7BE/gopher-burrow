# syntax=docker/dockerfile:1
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION}" -o /out/inventory ./cmd/inventory

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/inventory /inventory
EXPOSE 8080
USER nonroot:nonroot
HEALTHCHECK --interval=30s CMD ["/inventory", "health"]
ENTRYPOINT ["/inventory"]
