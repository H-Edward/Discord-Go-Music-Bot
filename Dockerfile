FROM golang:latest

WORKDIR /app

USER root
RUN adduser --disabled-password --gecos '' botuser
RUN apt-get update && \
    apt-get install -y ffmpeg libopus0 libopus-dev libopusenc0 libopusfile-dev opus-tools --no-install-recommends && \
    apt-get remove -y yt-dlp && \
    curl -L --fail https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp && \
    rm -rf /var/lib/apt/lists/*

COPY . .

RUN go mod download

RUN make build

RUN chown -R botuser:botuser /app

USER botuser

CMD ["./music-bot"]