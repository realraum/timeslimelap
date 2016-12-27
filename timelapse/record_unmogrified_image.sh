#!/bin/zsh
setopt extendedglob histsubstpattern magicequalsubst cshnullglob braceccl


export leddev=/dev/ttyUSB0
export v4ldev=( /dev/video?([-1]) )

GST_LAUNCH=gst-launch-1.0

function takePicture {
  local outfilename="./test.jpg"
  $GST_LAUNCH  v4l2src device=$v4ldev num-buffers=40 brightness=$((2147483647)) ! jpegenc ! image/jpeg,width={ 1920, 1280, 1024, 864, 640 },height={ 1080, 768, 720, 480 },pixel-aspect-ratio=1/1 ! filesink location="$outfilename"
}


#$GST_LAUNCH  v4l2src device=/dev/video0 num-buffers=1 ! image/jpeg,width={ 1280, 864, 640 },height={ 720, 480 },framerate={ 25/1, 24/1, 30/1 },pixel-aspect-ratio=1/1 ! jpegdec ! jpegenc ! filesink location="$outfilename"

takePicture


