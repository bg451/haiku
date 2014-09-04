var current_match
var ida = 0
var idb = 0
function get_new_match() {
  $.get("/matches/new", {"idA": ida, "idB": idb},
    function(data) {
      current_match = data;
      update_videos()
      console.log(current_match)
    });
}
function update_ids() {
  if (window.current_match != undefined) {
    ida = current_match.video_a.id;
    idb = current_match.video_b.id;
    }
}
function update_videos() {
  $('#videoA').attr('src', current_match.video_a.url);
  $('#videoB').attr('src', current_match.video_b.url);
}

function post_match(winnerA) {
  data = current_match;
  data.winnerA = winnerA;
  data = JSON.stringify(data)
  $.post("/matches/result", data, "json");
}

