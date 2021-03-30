FROM alpine:3.13

COPY binance /bin/binance

WORKDIR /home/binance

COPY views /home/binance

ENTRYPOINT ["binance"]