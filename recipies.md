# working example to run hardware-encoded streaming on buster:

```
./test-launch "rpicamsrc bitrate=800000 preview=false ! video/x-h264, width=640, height=480, framerate=30/1 ! h264parse ! rtph264pay name=pay0 pt=96"
```
