cmake_minimum_required(VERSION 3.20)
project(inventory C CXX)

set(CMAKE_C_STANDARD 17)
set(CMAKE_CXX_STANDARD 20)

add_executable(inventory_c main.c)
add_executable(inventory_cpp main.cpp)
