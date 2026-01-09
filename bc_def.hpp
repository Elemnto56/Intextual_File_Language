//
// Created by jaydog on 12/25/25.
//

#ifndef ITX_BC_DEF_HPP
#define ITX_BC_DEF_HPP
#include <vector>

#include "ir_gen.hpp"

enum BC {
    PUSH,
    PULL,
    STORE,
    PRINT,
    EXIT,
    MATH_
};

enum BC_SUB {
    ADD,
    SUB,
    MULT,
    DIV, // Division as in math
};

struct Instr {
    itx_types value;
    BC code;
    BC_SUB sub_code;
};


std::vector<Instr> bc_gen(std::vector<IR> ir);
void executor(std::vector<Instr> bc);
#endif //ITX_BC_DEF_HPP