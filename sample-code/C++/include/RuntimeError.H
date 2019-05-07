#ifndef RuntimeError_H
#define RuntimeError_H

#include <stdexcept>

/**
 * This class is a parent to all curl's RuntimeErrors.
 * This class takes a const std::string & as argument for it's parent: std::runtime_errors. 
 * This class is thrown when curlpp is encountering an error.
 */
class RuntimeError : public std::runtime_error
{
  public:
    RuntimeError(const char *reason);
    RuntimeError(const std::string &reason);
    virtual ~RuntimeError() throw();
};

#endif
