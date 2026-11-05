FROM golang:1.27.3-alpine AS build
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/raftkv ./cmd/raftkv

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/raftkv /raftkv
USER nonroot:nonroot
ENTRYPOINT ["/raftkv"]
