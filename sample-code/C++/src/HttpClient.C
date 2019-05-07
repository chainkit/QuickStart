#include <HttpClient.H>
#include <RuntimeError.H>
#include <cstring>
#include <iostream>

HttpClient::HttpClient()
{
  memset(errorBuffer,0,CURL_ERROR_SIZE + 1);
  curl = curl_easy_init();
  assertRuntime("Error when trying to curl_easy_init() a handle", curl == NULL);
  setOption(CURLOPT_ERRORBUFFER, (void *)errorBuffer);
  setOption(CURLOPT_WRITEFUNCTION, WriteCallback);
}

HttpClient::~HttpClient()
{
  curl_easy_cleanup(curl);
}

void HttpClient::assertRuntime(const char *errorMessage, CURLcode code)
{
  if (code != CURLE_OK) {
    throw RuntimeError(errorMessage);
  }
}

void HttpClient::assertRuntime(const char *errorMessage, bool shouldThrow)
{
  if (shouldThrow) {
    throw RuntimeError(errorMessage);
  }
}

size_t HttpClient::WriteCallback(char * buffer, size_t size, size_t nitems, std::string *str)
{
  size_t realwrote = size * nitems;
  str->append(buffer, realwrote);
  return realwrote;
}

std::string HttpClient::doGET(const char *url, KeyValueMap &headers)
{
  curl_slist *headerList = NULL;
  std::string postFields, responseString;

  for(KeyValueMap::iterator it = headers.begin();it != headers.end();++it)
  {
    std::string header = it->first + ": " + it->second;
    headerList = curl_slist_append(headerList, header.data());
  }
  setOption(CURLOPT_HTTPHEADER, headerList);
  setOption(CURLOPT_URL, url);
  setOption(CURLOPT_WRITEDATA, &responseString);
  CURLcode responseCode = curl_easy_perform(curl);
  assertRuntime("Request failed.", responseCode);
  return responseString;
}

std::string HttpClient::doPOST(const char *url, KeyValueMap &headers, std::string &postFieldJson)
{
  curl_slist *headerList = NULL;
  std::string postFields, responseString; 
 
  for(KeyValueMap::iterator it = headers.begin();it != headers.end();++it)
  {
    std::string header = it->first + ": " + it->second;
    headerList = curl_slist_append(headerList, header.data());
  }
  setOption(CURLOPT_HTTPHEADER, headerList);
  setOption(CURLOPT_POSTFIELDS, postFieldJson.data());
  setOption(CURLOPT_URL, url);
  setOption(CURLOPT_WRITEDATA, &responseString);
  CURLcode responseCode = curl_easy_perform(curl);
  assertRuntime("Request failed.", responseCode);
  return responseString;
}
