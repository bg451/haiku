var current_match;
function get_new_match() {
  $.get("/matches/new", function(data) {
    console.log(data);
    current_match = data;
    update_videos()
  });

}
function update_videos() {
  $('#videoA').attr('src', current_match.video_a.url);
  $('#videoB').attr('src', current_match.video_b.url);
}

function post_match(winnerA) {
  data = current_match;
  data.winnerA = winnerA;
  data = JSON.stringify(data)
  console.log(current_match)
  $.post("/matches/result", data, "json");
}

