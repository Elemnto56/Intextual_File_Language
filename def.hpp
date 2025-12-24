//
// Created by jaydog on 12/23/25.
//
// If you're wondering why this header is called "def", it's because I wanted a name to represent definitions for stuff

#ifndef INTEXTUAL_FILE_LANGUAGE_DEF_HPP
#define INTEXTUAL_FILE_LANGUAGE_DEF_HPP
#include <algorithm>
#include <string>
// Checks if string s contains the substr
std::string contains(std::string s, std::string substr) {
    if (s.find(substr) != std::string::npos) {

    }
}

enum LexType {
    OPERATOR,
    SYMBOL,
    ARROW,
    MACRO,
    FUNC,
    LOGIC
};

template <typename T>
struct Lex {
    std::string sub_type;
    std::string meta;
    T value;
    LexType type;
    int line;
};
#endif //INTEXTUAL_FILE_LANGUAGE_DEF_HPP