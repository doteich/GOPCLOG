
#Build Stage
FROM golang:1.19-alpine3.16 as build
WORKDIR /app
COPY . .
RUN go build -o /gopc_operator

#Copy files from build, to slim down the overall image size
FROM alpine:3.16.2
WORKDIR /root/
COPY --from=build ./app ./
COPY --from=build /gopc_operator ./
CMD ["./gopc_operator"] 