<?php
function requestToken($username, $password){

//Initiate cURL.
$ch = curl_init('https://api.chainkit.com/token');
 
 // First we need to get the JSON web token (JWT) -- so let's authenticate to the token api with our username and password
$jsonData = array(
    'userId' => $username,
    'password' => $password
);
 
//Encode the array into JSON.
$jsonDataEncoded = json_encode($jsonData);

//Tell cURL that we want to send a POST request.
curl_setopt($ch, CURLOPT_POST, 1);
 
//Attach our encoded JSON string to the POST fields.
curl_setopt($ch, CURLOPT_POSTFIELDS, $jsonDataEncoded);
 
//Set the content type to application/json
curl_setopt($ch, CURLOPT_HTTPHEADER, array('Content-Type: application/json')); 

curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
 
//Execute the request
$result = curl_exec($ch);

//Display Errors (if any)
if(curl_errno($ch)){
    echo 'Request Error in generating token:' . curl_error($ch);
}

//Close cURL
curl_close ($ch);

//Decode JSON into array
$json =  json_decode($result, true);

//The token api returns a number of items as json, and the one we care about is the access token, so let's pull that out.
return $json["data"]["accessToken"];

}

//This function calculates and generates the 256SHAsum of the file 
function generateHash($file_path) {
    $hash = hash_file("sha256", $file_path);
    return $hash;
     
}

// Now let's do the POST to the register API -- we pass in the JWT as an authorization bearer token header and then the post data is JSON in the body with the hash and any other settings you wish to send
function registerHash($hash, $storage){

global $token;

$result = array();
//Set other parameters as keys in the $postdata array
$postdata =  array('hash' => $hash , 'description' => "demo registration", 'storage' => $storage);
$url = "https://api.chainkit.com/register/";

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL,$url);
curl_setopt($ch, CURLOPT_POST, 1);
curl_setopt($ch, CURLOPT_POSTFIELDS,json_encode($postdata));  //Post Fields
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

$headers = [
  'Authorization: Bearer '.$token,
  'Content-Type: application/json',

];
curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);

$request = curl_exec($ch);

//Display Errors (if any)
if(curl_errno($ch)){
    echo 'Request Error in generating token:' . curl_error($ch);
}

curl_close ($ch);
// A successful registration returns the entity ID that was registered on the blockchain, along with some other metadata
return $request;

}




/*-------------------------------------------
This is the part that really concerns to you, 
just fill in your credentials and that's all!
--------------------------------------------*/
//Input username here.
echo "username: ";
$username = rtrim(fgets(STDIN));;
// echo "username= " .$username;
//input password here.
echo "passowrd: ";
system('stty -echo');
$password = trim(fgets(STDIN));
system('stty echo'); 
echo "\n";

// get storage type from users
echo "storage: ";
$storage = rtrim(fgets(STDIN));; 

//Input file path here.
$file = $argv[1];
echo "file: " .$file;
if($file == ''){
    $error = "Filename cannot be empty in \n";
    throw new Exception($error);
}
echo "\n";

$token = requestToken($username, $password); //Token is generated

//file hashed
$hash = generateHash($file); 

//Display generated hash value

//Registration 
$entityID = registerHash($hash, $storage); 

//The entity ID that was registered on the blockchain is displayed.
echo 'Registered Entity ID: ' .$entityID;
echo "\n";
