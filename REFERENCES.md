# Online References

Online references (in no particular order, just like bookmarks) which I needed to create this project.
- https://luminarys.com/posts/writing-a-bittorrent-client.html
- https://gist.github.com/CharlesHolbrow/8adfcf4915a9a6dd20b485228e16ead0
- https://docs.peer5.com/guides/production-ready-hls-vod/
- https://www.martin-riedl.de/2018/08/24/using-ffmpeg-as-a-hls-streaming-server-part-1/
- https://github.com/sampotts/plyr
- https://codepen.io/pen?template=oyLKQb
- https://stackoverflow.com/questions/31069002/ffmpeg-generate-m3u8-and-segments-manually
- https://mayur-solanki.medium.com/how-to-create-mpd-or-m3u8-video-file-from-server-using-ffmpeg-97e9e1fbf6a3
- https://superuser.com/questions/1150872/grab-vtt-subtitles-from-m3u8-stream
- https://stackoverflow.com/questions/25407474/how-to-create-dynamic-m3u8-by-pasting-the-url-in-browser
- https://superuser.com/questions/583393/how-to-extract-subtitle-from-video-using-ffmpeg
- http://seemer-dorfet.ch/cgi/manpages/ffmpeg-formats.pdf#page=31&zoom=100,0,736
- https://github.com/mifi/hls-vod
- https://stackoverflow.com/questions/56891221/generate-single-mpeg-dash-segment-with-ffmpeg
- https://superuser.com/questions/1255858/ffmpeg-hls-create-playlist-only-or-predict-segment-times-before-encoding
- https://github.com/UnicornTranscoder/UnicornTranscoder
- https://github.com/jellyfin/jellyfin/blob/master/Jellyfin.Api/Controllers/VideoHlsController.cs#L398
- https://gist.github.com/mrbar42/ae111731906f958b396f30906004b3fa
- https://kipalog.com/posts/FFMPEG-HLS-STREAM-MULTIPLE-AUDIO-SUBTITLES
- https://developer.apple.com/documentation/http_live_streaming/example_playlists_for_http_live_streaming/adding_alternate_media_to_a_playlist
- https://stackoverflow.com/questions/55649679/is-it-real-connect-subtitles-when-streaming-video-hls-m3u8
- https://johnvansickle.com/ffmpeg/
- https://github.com/divijbindlish/parse-torrent-name
- https://www.martin-riedl.de/2018/09/13/using-ffmpeg-as-a-hls-streaming-server-part-6-independent-segments/
- https://developer.apple.com/documentation/http_live_streaming/example_playlists_for_http_live_streaming/adding_alternate_media_to_a_playlist
- https://stackoverflow.com/questions/32252337/how-to-style-text-tracks-in-html5-video-via-css
- https://videojs-http-streaming.netlify.app/
- https://stackoverflow.com/questions/37091316/how-to-get-the-realtime-output-for-a-shell-command-in-golang
- https://stackoverflow.com/questions/34855343/how-to-stream-an-io-readcloser-into-a-http-responsewriter
- https://pkg.go.dev/io#Copy
- https://stackoverflow.com/questions/47216029/how-to-stream-http-content-to-a-file-and-a-buffer
- https://developer.apple.com/documentation/http_live_streaming/hls_authoring_specification_for_apple_devices
- https://stackoverflow.com/questions/61913288/on-the-fly-transcoding-and-hls-streaming-with-ffmpeg
- https://developer.apple.com/documentation/http_live_streaming/example_playlists_for_http_live_streaming/creating_a_primary_playlist
- https://developer.apple.com/documentation/http_live_streaming/example_playlists_for_http_live_streaming/video_on_demand_playlist_construction
- https://datatracker.ietf.org/doc/html/rfc8216
- https://www.nginx.com/wp-content/uploads/2018/12/NGINX-Conf-2018-slides_Choi-streaming.pdf
- https://i.imgur.com/r4NtWy2.png
- https://blog.zazu.berlin/internet-programmierung/mpeg-dash-and-hls-adaptive-bitrate-streaming-with-ffmpeg.html
- https://superuser.com/questions/1308746/splitting-re-encoding-and-joining-video-files-results-in-clicking-audio
- https://stackoverflow.com/questions/29527882/ffmpeg-copyts-to-preserve-timestamp
- https://gitlab.com/olaris/olaris-server/-/commit/6de5c2f80e3e7da42c53630ace7e88adf55c6199
- https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/QTFFAppenG/QTFFAppenG.html
- https://ffmpeg.org/ffmpeg-formats.html#segment
- https://gist.github.com/steven2358/ba153c642fe2bb1e47485962df07c730
- http://underpop.online.fr/f/ffmpeg/help/hls-2.htm.gz
# Thoughts and ideas
This is how jellyfin transcodes media.
```
/usr/lib/jellyfin-ffmpeg/ffmpeg -fflags +genpts -f matroska,webm -i file:"/XXXXXX.mkv" -map_metadata -1 -map_chapters -1 -threads 0 -map 0:0 -map 0:1 -map -0:s -codec:v:0 copy -bsf:v h264_mp4toannexb -vsync -1 -codec:a:0 aac -ac 2 -ab 384000 -af "volume=2" -copyts -avoid_negative_ts disabled -f hls -max_delay 5000000 -hls_time 6 -individual_header_trailer 0 -hls_segment_type mpegts -start_number 0 -hls_segment_filename "/config/data/transcodes/92f62516d14f6390b6ae3c93d0837690%d.ts" -hls_playlist_type vod -hls_list_size 0 -y "/config/data/transcodes/92f62516d14f6390b6ae3c93d0837690.m3u8"
```
```
startInfo.Arguments = string.Format("-report -loglevel warning -hide_banner -progress pipe:2 -vsync passthrough -noaccurate_seek -copyts -ss 300.5 -i \"{0}\" -c:v libx264 -preset veryfast -b:v 80M -segment_time 10 -map 0 -segment_list pipe:1 -segment_list_type csv -f ssegment \"{1}\"", _sourceFile, _outputFilenameFormat);

```
```
https://app.element.io/#/room/#jellyfin-dev:matrix.org/$gcdIn-SqsYS4O8qaLjGZrn7MRY4F5ir7cZRyIdhKmZw
```
https://gist.github.com/cvium/c5a65e92e426b29f9278ab404e8421a2