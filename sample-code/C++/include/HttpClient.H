#ifndef HttpClient_H
#define HttpClient_H

#include <curl/curl.h>
#include <string>
#include <sstream>
#include <map>


/**
 * This is http client which makes GET and POST call to a given url.
 * This is C++ wrapper over libcURL 
 */
class HttpClient
{
  public:
    typedef std::map<std::string, std::string> KeyValueMap;
    static size_t WriteCallback(char * buffer, size_t size, size_t nitems, std::string *str);

    HttpClient();
    virtual ~HttpClient();

    /**
     * This makes a GET call. User can pass header to the method.
     */
    std::string doGET(const char *url, KeyValueMap &headers);

    /**
     * This makes a POST call. User can pass headers and post fields to the method.
     */
    std::string doPOST(const char *url, KeyValueMap &headers, std::string &postFieldJson);

  private:
    /**
     * Calls curl_easy_setopt on the option.
     */
    template<typename OptionType>
    void setOption(CURLoption optionType, OptionType value)
    {
      CURLcode code;
      code = curl_easy_setopt(curl, optionType, value);
      assertRuntime(errorBuffer, code);
    }

    /**
     * If CURLcode is not equal to CURLE_OK, it throws a 
     * RunTimeException.
     */
    void assertRuntime(const char *errorMessage, CURLcode code);

    /**
     * It throws RunTimeException, if callers ask to throw.
     * It is genrally called if we fail to initialize libcurl object.
     */
    void assertRuntime(const char *errorMessage, bool shouldThrow);

  private:
    CURL *curl;
    char errorBuffer[CURL_ERROR_SIZE + 1];
};

#endif
