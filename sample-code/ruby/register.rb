#!/usr/bin/env ruby
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
puts "storage: "
storage = gets.chomp
puts "filename: "
filename = gets.chomp

params   = { 'userId': username, 'password': password }
headers  = { 'Content-type': 'application/json' }

uri = URI.parse("https://api.pencildata.com/token")
https = Net::HTTP.new(uri.host, uri.port)
https.use_ssl = true

request = Net::HTTP::Post.new(uri.path, initheader = {'Content-Type' =>'application/json'})
request.body = params.to_json

response = https.request(request)
token = JSON.parse(response.body)['data']['accessToken']

# Register
url          = 'https://api.pencildata.com/register/'
file = File.new(filename, "r")
hash_content = Digest::SHA256.hexdigest(file.read)
description  = 'demo registration'
storage      = storage

params       = { hash: hash_content, description: description, storage: storage }.to_json
headers      = { Authorization: "Bearer #{token}", content_type: 'application/json' }

response  = RestClient.post(url, params, headers)    
entity_id = response
puts "Use this entity_id for verification."
puts entity_id
