#include <RuntimeError.H>

RuntimeError::RuntimeError(const char *reason) : std::runtime_error(std::string(reason)) 
{
}

RuntimeError::RuntimeError(const std::string &reason) : std::runtime_error(reason)
{
}

RuntimeError::~RuntimeError() throw()
{
}
