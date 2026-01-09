#ifndef ITX_IR_GEN_HPP
#define ITX_IR_GEN_HPP
#include <iostream>
#include <vector>
#include <unordered_map>

#include "lex_def.hpp"
enum ErrorType {
    MISSING_TOKEN,
    BAD_TYPE,
    LEX_ERR
};

inline void callErr(ErrorType et, std::string err_desc, int line, int status_code = 1) {
    std::cerr << "line: " << line << std::endl;
    std::cerr << err_desc << std::endl;

    switch (et) {
    case BAD_TYPE: std::cerr << "hint: The four built-in types are int, float, string, and bool" << std::endl;
    default: break;
    }

    exit(status_code);
}

using ir_type = std::variant<std::string, int, float, bool, std::vector<Lex>>;
struct IR {
    itx_types rvalue;
    std::unordered_map<std::string, ir_type> metadata;
    std::string lvalue;
    LexMainType l_type;
    LexSubType r_type;
    int line;
};

std::vector<IR> ir_gen(std::vector<Lex> tokens);
#endif //ITX_IR_GEN_HPP