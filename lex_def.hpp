//
// Created by jaydog on 12/23/25.
//
// If you're wondering why this header is called "def", it's because I wanted a name to represent definitions for stuff

#ifndef INTEXTUAL_FILE_LANGUAGE_DEF_HPP
#define INTEXTUAL_FILE_LANGUAGE_DEF_HPP
#include <filesystem>
#include <vector>
#include <string>
#include <variant>
#include <unordered_map>

using itx_types = std::variant<std::string, int, float, bool>;

bool contains(std::string test, std::vector<std::string> finds);
std::vector<std::string> split_by_ws(std::string test);

enum LexSubType {
    BUILT_IN,
    STRING,
    INT,
    FLOAT,
    BOOL,
    ARRAY,
    L_PARA,
    R_PARA,
    VAR,
    MATH,
    COMPARE,
    R_CURL
};

enum LexMainType {
    OPERATOR,
    SYMBOL,
    MACRO,
    FUNC,
    LOGIC,
    KEYWORD,
    L_KEYWORD, // Logic keywords such as if, while, for
    TYPE,
    VARIABLE,
    BODY
};

struct Lex {
    LexSubType sub_type;
    std::string meta;
    itx_types value;
    LexMainType type;
    std::vector<Lex> scope;
    int line;
};

inline std::string resolve(const itx_types& t) {
    switch (t.index()) {
    case 0: return "string";
    case 1: return "int";
    case 2: return "float";
    default: return "bool";
    }
}

inline std::unordered_map<LexSubType, std::string> type_look {
    {STRING, "string"},
    {INT, "int"},
    {FLOAT, "float"},
    {BOOL, "bool"}
};

inline void callErr(std::string_view err_desc, const int line, const int status_code = 1) {
    std::cerr << "line: " << line << std::endl
    << err_desc << std::endl;
    exit(status_code);
}
std::vector<Lex> lexer(std::vector<std::string> lines, char stop_at = '~');
#endif //INTEXTUAL_FILE_LANGUAGE_DEF_HPP