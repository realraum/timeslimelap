#!/bin/zsh

[[ -n $1 ]] && cd $1

GST_LAUNCH=gst-launch-1.0
FRAMES_D=.

rm -v frame*.jpg(L0)

$GST_LAUNCH multifilesrc location="$FRAMES_D/frame_%06d.jpg" loop=true ! image/jpeg,width=1080,height=1080,framerate=25/1 ! jpegdec ! videoconvert ! xvimagesink
#$GST_LAUNCH multifilesrc location="$FRAMES_D/frame_%06d.jpg" loop=true ! image/jpeg,width=1080,height=1080,framerate=25/1 ! jpegdec ! videoconvert ! ximagesink
