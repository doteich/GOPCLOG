
#Build Stage
FROM golang:1.19-alpine3.16 as build
WORKDIR /app
COPY . .
RUN go build -o /gopclogs

#Copy files from build, to slim down the overall image size
FROM alpine:3.16.2
WORKDIR /root/
COPY --from=build ./app ./
COPY --from=build /gopclogs ./
CMD ["./gopclogs"] 

