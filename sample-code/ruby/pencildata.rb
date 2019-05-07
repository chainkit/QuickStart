require 'rest-client'
require 'json'
require 'base64'
require 'digest'

# Authentication
url      = "https://api.pencildata.com/token"
username = 'XXX'
password = Base64.encode64('XXX')

params   = { username: username, password: password }.to_json
headers  = { content_type: 'application/json' }

json     = RestClient.post(url, params, headers)
response = JSON.parse(json)	
token    = response['result']['accessToken']

# Register
url          = 'https://api.pencildata.com/register/'
hash_content = Digest::SHA256.hexdigest('message')
description  = 'demo registration'
storage      = 'private'

params       = { hash: hash_content, description: description, storage: storage }.to_json
headers      = { Authorization: "Bearer #{token}", content_type: 'application/json' }

response  = RestClient.post(url, params, headers)    
entity_id = response

# Verify
url          = "https://api.pencildata.com/verify/#{entity_id}?hash=#{hash_content}&storage=#{storage}"
headers      = { Authorization: "Bearer #{token}", content_type: 'application/json' }

response = RestClient.get(url, headers)
