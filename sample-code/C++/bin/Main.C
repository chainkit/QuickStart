#include <HttpClient.H>
#include <iostream>
#include <sstream>
#include <fstream>
#include <nlohmann/json.hpp>
#include <sha256.h>


std::string getAuthorizationToken(const std::string &username, const std::string &password);
std::string prepareHashForAFile(const std::string &filename);
std::string registerFile(const std::string &token, const std::string &hash);
void verifyToken(const std::string &token, const std::string &hash, const std::string &entityId);

int main(int argc, char *argv[])
{
  if (argc < 4)
  {
    std::cout << "[USAGE] " << argv[0] << " <username> <password> <filepath>" << std::endl;
    return -1;
  }
  std::cout << "***Start of pencildata api testing***" << std::endl;
  std::cout << "***Registration process in pencildata***" << std::endl;
  std::string accessToken = getAuthorizationToken(argv[1], argv[2];
  std::string hash = prepareHashForAFile(argv[3]);
  std::string entityId = registerFile(accessToken, hash);

  std::cout << std::endl << std::endl;
  std::cout << "***Verification process in pencildata***" << std::endl;
  accessToken = getAuthorizationToken(argv[1], argv[2]);
  hash = prepareHashForAFile(argv[3]);
  verifyToken(accessToken, hash, entityId);
  std::cout << "**Done with pencildata api testing***" << std::endl;
  std::cout << "***Exiting now***" << std::endl;
  return 0;
}


/**
 * Find authorization token using a given username and password
 */
std::string getAuthorizationToken(const std::string &username, const std::string &password)
{
  std::cout << "Attempting to authenticate user with credentials(" << username << "," << password << ")." << std::endl;
  HttpClient httpClient;
  const char *url = "https://api.pencildata.com/token";

  std::string postFieldJsonString;
  std::map<std::string, std::string> headers, postFields;
  nlohmann::json postFieldJson;

  headers["Content-Type"] = "application/json";
  postFieldJson["username"]  = username;
  postFieldJson["password"]  = password;


  postFieldJsonString = postFieldJson.dump();
  std::string responseString = httpClient.doPOST(url, headers, postFieldJsonString);
  nlohmann::json tokenResponse = nlohmann::json::parse(responseString);
  
  if (responseString != "" && tokenResponse["statusCode"] == 200)
  {
    std::cout << "Successfully authenticated user." << std::endl;
  }
  else
  {
    std::cerr << "Failed to authenticate user." << std::endl;
    exit(1);
  }
  std::string accessToken = tokenResponse["result"]["idToken"].get<std::string>();
  return accessToken;
}

std::string prepareHashForAFile(const std::string &filename)
{
  std::cout << "Attempting to hash file: " << filename << std::endl;
  SHA256 sha256;
  const size_t bufferSize = 144*7*1024;
  char* buffer = new char[bufferSize];

  std::ifstream file;
  std::istream* input = NULL;
  
  // open file
  file.open(filename.c_str(), std::ios::in | std::ios::binary);
  if (!file)
  {
    std::cerr << "Can't open '" << filename << "'." << std::endl;
    exit(1);
  }
  while(file)
  {
     file.read(buffer, bufferSize);
     std::size_t numBytesRead = size_t(file.gcount());

     sha256.add(buffer, numBytesRead);
  }
  // clean up
  file.close();
  delete[] buffer;
  std::cout << "Successfully prepared hash: " << sha256.getHash() << std::endl;
  return sha256.getHash();
}

std::string registerFile(const std::string &token, const std::string &hash)
{
  std::cout << "Attempting to register hash: " << hash << " with token: " << token << std::endl;
  HttpClient httpClient;
  const char *url = "https://api.pencildata.com/register";
  std::string postFieldJsonString;
  std::map<std::string, std::string> headers, postFields;
  nlohmann::json postFieldJson;

  headers["Authorization"] = "Bearer " + token;
  headers["Content-type"] = "application/json";

  postFieldJson["hash"]  = hash;
  postFieldJson["storage"]  = "private";


  postFieldJsonString = postFieldJson.dump();
  std::string entityId = httpClient.doPOST(url, headers, postFieldJsonString);
  if (entityId != "")
  {
    std::cout << "Successfully regsitered the hash. EntityId: " << entityId << std::endl;
  }
  else 
  {
    std::cerr << "Failed to regsitered the hash." << std::endl;
    exit(1);
  }
  return entityId;
}

void verifyToken(const std::string &token, const std::string &hash, const std::string &entityId)
{
  std::cout << "Attempting to verify entityId: " << entityId << " with hash: " << hash << " with token: " << token << std::endl;
  HttpClient httpClient;
  std::string url = "https://api.pencildata.com/verify/" + entityId + "?hash=" + hash + "&storage=private";
  std::cout << url << std::endl;
  
  std::map<std::string, std::string> headers;
  headers["Authorization"] = "Bearer " + token;

  std::string response = httpClient.doGET(url.data(), headers);
  if (response == "true")
  {
    std::cout << "EntityId: " << entityId << " has been verified." << std::endl;
  }
  else
  {
    std::cerr << "EntityId: " << entityId << " could not be verified." << std::endl;
  }
}






