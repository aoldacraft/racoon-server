# Build the application from source
FROM golang:alpine AS build-stage

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

#COPY go.mod go.sum ./
COPY . .
RUN go mod download

#COPY *.go ./

RUN go build -o /rc_server

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /rc_server /rc_server

ENV DB_HOST rc_db
ENV DB_PORT 5432
ENV DB_USER racoondb
ENV DB_PASSWORD racoondb
ENV DB_NAME racoondb

EXPOSE 1323

ENTRYPOINT ["/rc_server"]

## 기본 이미지를 설정
#FROM golang:alpine AS build-stage
#
## 작업 디렉토리 설정
#WORKDIR /go/src/rc_app
#
## 소스 코드 복사
#COPY . .
#
## 환경 변수 설정
#ENV DB_HOST rc_db
#ENV DB_PORT 5432
#ENV DB_USER racoondb
#ENV DB_PASSWORD racoondb
#ENV DB_NAME racoondb
#
## 의존성 다운로드
#RUN go mod download
#
## 어플리케이션 빌드
#RUN go build -o rc_app
#
## 포트 노출
#EXPOSE 1323
#
## 어플리케이션 실행
#ENTRYPOINT ["./rc_app"]
