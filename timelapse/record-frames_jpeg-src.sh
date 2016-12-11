#!/bin/zsh
setopt extendedglob histsubstpattern magicequalsubst cshnullglob braceccl

GST_LAUNCH=gst-launch-1.0

export leddev=/dev/ttyUSB0
export v4ldev=( /dev/video?([-1]) )
local intervall=20


local -i imgnum=0
local outdir=./timepics${1:-}
[[ ! -d $outdir ]] && mkdir $outdir
local tmpfile=$(mktemp).jpg

newestfiles=( ${outdir}/frame_<->.jpg(N[-1]) )
if [[ ${newestfiles} == (#b)${outdir}/frame_(<->).jpg ]] then
  imgnum=$((match+1))
fi

stty -F $leddev 57600 raw min 1 time 1 brkint
sleep 0.1
cat $leddev >| /dev/null < /dev/null &
local CATPID=$!
echo $CATPID
trap "kill -9 $CATPID" EXIT
sleep 0.5
echo -n "010" > $leddev

function takePicture {
  TRAPEXIT() {
      echo -n 0 > $leddev
  }
  local dir=$1
  local imgnum=$2
  local outfilename="${dir}/frame_$(print -f "%06d" imgnum).jpg"
  echo -n 1 > $leddev
  sleep 5 # give cam time to adjust
  $GST_LAUNCH  v4l2src device=$v4ldev num-buffers=40 !  textoverlay text="$((imgnum*intervall))s" line-alignment=0 halignment=2 ! jpegenc ! image/jpeg,width={ 1920, 1280, 1024, 864, 640 },height={ 1080, 768, 720, 480 },pixel-aspect-ratio=1/1 ! filesink location="$outfilename"
  echo -n 0 > $leddev
  ## Creative Live Camera, 640x480
  #mogrify -distort Perspective '40,36 0,0 680,50 640,0 141,401 0,480 556,409 640,480 ' "$outfilename"
  ## Microsoft Camera 1280x720
  mogrify -distort Perspective '330,0 0,0 1480,0 1080,0 0,1070 0,1080, 1950,1070 1080,1080' -crop 1080x1080+0+0 "$outfilename"
}


#$GST_LAUNCH  v4l2src device=/dev/video0 num-buffers=1 ! image/jpeg,width={ 1280, 864, 640 },height={ 720, 480 },framerate={ 25/1, 24/1, 30/1 },pixel-aspect-ratio=1/1 ! jpegdec ! jpegenc ! filesink location="$outfilename"

local LASTPID
while sleep $intervall; do
  [[ -n $LASTPID ]] && wait $LASTPID
  takePicture $outdir $imgnum &
  LASTPID=$!
  echo $((imgnum++))
done

echo 
echo EXITING....
echo
wait $LASTPID
exit 0
#convert "$tmpfile" -crop "1080x1080+0+0"  "$outfilename"
