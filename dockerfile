FROM golang

WORKDIR /app

# Copy the source code
COPY . .

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o app .

# 환경 변수를 설정하여 데이터베이스 연결 정보를 제공
ENV DB_HOST rc_db
ENV DB_PORT 5432
ENV DB_USER racoondb
ENV DB_PASSWORD racoondb
ENV DB_NAME racoondb

# Run the executable
CMD ["./app"]