cmake_minimum_required(VERSION 2.8)
project(pencildata)

# set C++11 compilation standard
set(CMAKE_CXX_STANDARD 11)

# project version
set(PROJECT_VERSION_MAJOR 1)
set(PROJECT_VERSION_MINOR 2)
set(PROJECT_VERSION_BUILD 0)

# This will provide verbose output from make
set(CMAKE_VERBOSE_MAKEFILE ON)


# Defining the name of the executable.
set(EXE_FILE_NAME "${PROJECT_NAME}_v${PROJECT_VERSION_MAJOR}.${PROJECT_VERSION_MINOR}.${PROJECT_VERSION_BUILD}${TARGET_POSTFIX}")

# Json library
include(ExternalProject)
ExternalProject_Add(jsonlib
    PREFIX ${CMAKE_CURRENT_BINARY_DIR}/jsonlibinstall
    URL "https://github.com/nlohmann/json/releases/download/v3.1.2/include.zip"
    CONFIGURE_COMMAND ""
    BUILD_COMMAND ""
    BUILD_IN_SOURCE 1
    INSTALL_COMMAND ""
)

# sha256 library
ExternalProject_Add(sha256lib
    PREFIX ${CMAKE_CURRENT_BINARY_DIR}/sha256libinstall
    GIT_REPOSITORY "https://github.com/JusticeRage/hash-library.git"
    CONFIGURE_COMMAND ""
    BUILD_COMMAND g++ -c -std=c++11 sha256.cpp && ar -cvq libhash-library.a sha256.o
    BUILD_IN_SOURCE 1
    INSTALL_COMMAND cp sha256.h hash.h libhash-library.a ${CMAKE_CURRENT_BINARY_DIR}/sha256libinstall
)


set(SHA256Lib_LIBRARIES ${CMAKE_CURRENT_BINARY_DIR}/sha256libinstall/libhash-library.a)

include_directories("${CMAKE_CURRENT_BINARY_DIR}/jsonlibinstall/src/jsonlib")
include_directories("${CMAKE_CURRENT_BINARY_DIR}/sha256libinstall")
include_directories(include)

# install directory for the project, where binary would be created
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY "${CMAKE_SOURCE_DIR}/install/")

FILE(GLOB SRC "src/*.C")
FILE(GLOB BIN "bin/*.C")

add_executable(${EXE_FILE_NAME} ${SRC} ${BIN})
set_target_properties("${EXE_FILE_NAME}" PROPERTIES LINKER_LANGUAGE CXX)
target_link_libraries("${EXE_FILE_NAME}" "curl" ${SHA256Lib_LIBRARIES})
add_dependencies("${EXE_FILE_NAME}" jsonlib)
add_dependencies("${EXE_FILE_NAME}" sha256lib)
