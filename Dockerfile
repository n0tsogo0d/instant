FROM golang:alpine3.14
WORKDIR /download
RUN apk add wget
RUN wget -nv \
    https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz
RUN tar xf ffmpeg-release-amd64-static.tar.xz
WORKDIR /build
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o instant

FROM scratch
COPY --from=0 /download/ffmpeg-*-amd64-static/ffmpeg /usr/bin/ffmpeg
COPY --from=0 /download/ffmpeg-*-amd64-static/ffprobe /usr/bin/ffprobe
COPY --from=0 /build/instant /instant
COPY --from=0 /build/web /web
EXPOSE 8000
CMD ["/instant"]