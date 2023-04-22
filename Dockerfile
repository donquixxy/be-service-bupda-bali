# builder image
FROM golang:1.20-alpine as builder
WORKDIR /build
COPY . .
RUN apk add git && CGO_ENABLED=0 GOOS=linux go build -o be-service-bupda-bali .
# RUN go build -o be-service-teman-bunda .

# generate clean, final image for end users
FROM alpine
RUN apk add --no-cache curl
RUN apk update && apk add ca-certificates && apk add tzdata && apk add git
COPY --from=builder /build .
ENV TZ="Asia/Makassar"
# Dev
# EXPOSE 9090 

# Prod
EXPOSE 9080

CMD ./be-service-bupda-bali

# Add the health check instruction
#Production
HEALTHCHECK --interval=30s --timeout=10s CMD curl -f https://eagle-go.bupdabali.com/ || kill 1

#Dev
# HEALTHCHECK --interval=30s --timeout=10s CMD curl -f http://117.53.44.216:9090/ || kill 1