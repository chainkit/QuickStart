<?php
function requestToken($username, $password){

//Initiate cURL.
$ch = curl_init('https://api.pencildata.com/token');
 
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


function verifyEntityID($entityID, $hash, $storage){

global $token;

$ch = curl_init();
 
//Set the URL that you want to GET by using the CURLOPT_URL option.
curl_setopt($ch, CURLOPT_URL, 'https://api.pencildata.com/verify/'.$entityID.'?storage='.$storage.'&hash='.$hash);

//Set headers with token
$headers = [
  'Authorization: Bearer '.$token,
  'Content-Type: application/json',

];

curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
 
//Set CURLOPT_RETURNTRANSFER so that the content is returned as a variable.
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
 
//Set CURLOPT_FOLLOWLOCATION to true to follow redirects.
curl_setopt($ch, CURLOPT_FOLLOWLOCATION, true);
 
//Execute the request.
$data = curl_exec($ch);
 
//Close the cURL handle.
curl_close($ch);
 
//Print the data out onto the page.
return $data;
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

// get EntityId from users
echo "entityID: ";
$entityID = rtrim(fgets(STDIN));;
//Input file path here (Absolute file path or via post request($FILE["tmp_name"])).
$file = $argv[1];
echo "file: " .$file;
echo "\n";
if($file == ''){
    $error = "Filename cannot be empty in \n";
    throw new Exception($error);
}
$hash = generateHash($file); ; 

$token = requestToken($username, $password); //Token is generated.

$verify = verifyEntityID($entityID, $hash, $storage); // EntityID vs Hash Veification.

//Verification status is displayed. 
echo $verify;
echo "\n";

