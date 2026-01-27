#ifndef ITX_IR_GEN_HPP
#define ITX_IR_GEN_HPP
#include <iostream>
#include <vector>
#include <unordered_map>

#include "lex_def.hpp"

struct IR;
using ir_type = std::variant<std::vector<Lex>, std::vector<IR>>;
struct IR {
    itx_types rvalue;
    std::unordered_map<std::string, ir_type> metadata;
    std::string lvalue;
    LexMainType l_type;
    LexSubType r_type;
    int line;
};

std::vector<IR> ir_check(const std::vector<IR>& ir);
std::vector<IR> ir_gen(std::vector<Lex> tokens);
#endif //ITX_IR_GEN_HPP