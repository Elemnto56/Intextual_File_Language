cmake_minimum_required(VERSION 4.0)
project(itx)

set(CMAKE_CXX_STANDARD 23)

set(Boost_USE_STATIC_LIBS OFF)
set(Boost_USE_MULTITHREADED ON)
set(Boost_USE_STATIC_RUNTIME OFF)
find_package(Boost 1.89.0 COMPONENTS regex)

include_directories(${Boost_INCLUDE_DIRS})
include_directories(.)

add_executable(itx
        lex_def.hpp
        lexer.cpp
        bc_generator.cpp
        main.cpp
        bc_def.hpp
        ir_gen.hpp
        ir_generator.cpp
        vm.cpp
        evaluator.hpp

)
target_link_libraries(itx ${Boost_LIBRARIES})