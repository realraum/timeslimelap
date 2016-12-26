var $r3jq = jQuery.noConflict();

var timelapsedata={
  path:"/timepics/",
  imglist:["frame_000000.jpg","frame_000001.jpg","frame_000002.jpg","frame_000003.jpg"],
  interval_ms:1000,
}

var timelapsetimer=false;

function displayRecentImage()
{
  if (timelapsedata.imglist.length > 0)
  {
    $r3jq("#mostrecentimg").attr("src",timelapsedata.path+timelapsedata.imglist[timelapsedata.imglist.length-1]);
  }
}

function jumpInTimeLapse(distance)
{
  if (timelapsedata.imglist == 0)
  {
    return;
  }
  var alen=timelapsedata.imglist.length;
  var tlimgelem = $r3jq("#timelapseimg");
  var tlarrayidx=parseInt(tlimgelem.attr("tlarrayidx"),10);
  var nextidx = tlarrayidx+distance;
  //does negative modulo work in js ?? fear not...
  if (nextidx >= alen) {
    nextidx = 0;
  } else if (nextidx < 0) {
    nextidx = alen - 1;
  }
  tlimgelem.attr("src",timelapsedata.path+timelapsedata.imglist[nextidx]);
  tlimgelem.attr("tlarrayidx",nextidx);
}

function displayNextTimeLapse()
{
  jumpInTimeLapse(1);
}

function playTimeLapse()
{
  if (timelapsetimer == false)
  {
    timelapsetimer = setInterval(displayNextTimeLapse, timelapsedata.interval_ms);
    $r3jq("#button-timelapse-play").css("background-color","green");
    $r3jq("#button-timelapse-stop").css("background-color","");
  }
}

function stopTimeLapse()
{
  if (timelapsetimer) {
    clearInterval(timelapsetimer);
    timelapsetimer=false;
    $r3jq("#button-timelapse-play").css("background-color","");
    $r3jq("#button-timelapse-stop").css("background-color","red");
  }
}

function nextTimeLapse()
{
  stopTimeLapse();
  jumpInTimeLapse(1);
}


function prevTimeLapse()
{
  stopTimeLapse();
  jumpInTimeLapse(-1);
}

$r3jq(document).ready(function()
{
  $r3jq("#button-timelapse-stop").click(stopTimeLapse);
  $r3jq("#button-timelapse-play").click(playTimeLapse);
  $r3jq("#button-timelapse-next").click(nextTimeLapse);
  $r3jq("#button-timelapse-prev").click(prevTimeLapse);
  displayRecentImage();
  displayNextTimeLapse();
  playTimeLapse();

});