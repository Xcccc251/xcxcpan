package define

var TS_NAME = "index.ts"

var CMD_TRANSFER_TO_TS_WITH_CUT = "ffmpeg -i %s  -c:v libx264 -c:a aac -hls_time 10 -hls_segment_filename \"output_%03d.ts\" \"output.m3u8\""

var M3U8_NAME = "index.m3u8"
