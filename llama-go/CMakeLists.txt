cmake_minimum_required(VERSION 3.15)
project(llama_go)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Configure llama.cpp build (minimal, embedding-only)
set(LLAMA_BUILD_TESTS OFF CACHE BOOL "" FORCE)
set(LLAMA_BUILD_EXAMPLES OFF CACHE BOOL "" FORCE)
set(LLAMA_BUILD_SERVER OFF CACHE BOOL "" FORCE)
set(LLAMA_EMBEDDINGS ON CACHE BOOL "" FORCE)

# Add llama.cpp submodule
add_subdirectory(llama.cpp)

# llama_go wrapper (shared library)
add_library(llama_go SHARED
    llama_go.cpp
)

target_include_directories(llama_go PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/llama.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/llama.cpp/include
)

target_link_libraries(llama_go PRIVATE
    llama
)

# Set output name
set_target_properties(llama_go PROPERTIES
    OUTPUT_NAME "llama_go"
    POSITION_INDEPENDENT_CODE ON
)

# Install rules
install(TARGETS llama_go
    LIBRARY DESTINATION lib
    ARCHIVE DESTINATION lib
)

install(FILES llama_go.h
    DESTINATION include
)
