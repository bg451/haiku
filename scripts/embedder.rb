require 'json'
require 'rest_client'
# https://www.youtube.com/watch?v=QKnVWr2PkIU
# //www.youtube.com/embed/QKnVWr2PkIU
f = File.open("add_links.txt")
link_list = []
f.each_line do |line|
  id = line.split("=").last
  id = id.chomp
  link_list << "//youtube.com/embed/#{id}"
end
link_list.each do |l|
  RestClient.post "http://haitube.herokuapp.com/videos/new", {'url' => l}.to_json
end
