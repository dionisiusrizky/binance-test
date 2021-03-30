build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o binance ./cmd 
	docker build . -t binance:latest

run: build
	docker run -p 8080:8080 binance:latest

