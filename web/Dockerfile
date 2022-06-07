FROM golang:alpine AS backend-builder
RUN apk add --no-cache make build-base
WORKDIR /go/src/app
COPY . .
RUN go mod download 
RUN go install ./... . 

FROM node:lts-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

FROM alpine
WORKDIR /app
RUN apk add --no-cache curl
COPY --from=frontend-builder /app/dist ./
COPY --from=backend-builder /go/bin/6tisch-cloud-web ./
RUN mkdir templates &&\
    mv ./index.html ./templates/index.html &&\
    sed -i 's/^/{{define "index"}}/' ./templates/index.html &&\
    sed -i 's/$/{{end}}/' ./templates/index.html 
EXPOSE 80
CMD ["./6tisch-cloud-web"]

# docker buildx build -t amyangxyz111/6tisch-cloud-web --platform linux/arm64,linux/amd64,linux/arm/v7 --push .