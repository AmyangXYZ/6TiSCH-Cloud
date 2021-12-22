FROM golang:alpine AS backend-builder
RUN apk add --no-cache make build-base
WORKDIR /go/src/app
COPY . .
RUN go mod download 
RUN go install ./... . 


FROM alpine
WORKDIR /app
COPY --from=backend-builder /go/bin/6tisch-cloud-db ./
EXPOSE 54321
CMD ["./6tisch-cloud-db"]