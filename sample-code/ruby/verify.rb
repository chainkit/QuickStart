require 'rest-client'
require 'json'
require 'base64'
require 'digest'
require 'io/console'
require 'net/http'
require 'uri'


# Authentication
# url      = "https://api.pencildata.com/token"
puts "username: "
username = gets.chomp
puts "password: "
password = STDIN.noecho(&:gets).chomp
# password = Base64.encode64('XXX')
puts "storage: "
storage = gets.chomp
puts "filename: "
filename = gets.chomp
puts "entity_id: "
entity_id = gets.chomp

params   = { 'userId': username, 'password': password }
headers  = { 'Content-type': 'application/json' }

uri = URI.parse("https://api.pencildata.com/token")
https = Net::HTTP.new(uri.host, uri.port)
https.use_ssl = true

request = Net::HTTP::Post.new(uri.path, initheader = {'Content-Type' =>'application/json'})
request.body = params.to_json

response = https.request(request)
token = JSON.parse(response.body)['data']['accessToken']

file = File.new(filename, "r")
hash_content = Digest::SHA256.hexdigest(file.read)

# Verify
url          = "https://api.pencildata.com/verify/#{entity_id}?hash=#{hash_content}&storage=#{storage}"
headers      = { Authorization: "Bearer #{token}", content_type: 'application/json' }

response = RestClient.get(url, headers)

puts response
