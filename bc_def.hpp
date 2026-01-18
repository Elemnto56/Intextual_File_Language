//
// Created by jaydog on 12/25/25.
//

#ifndef ITX_BC_DEF_HPP
#define ITX_BC_DEF_HPP
#include "ir_gen.hpp"

enum BC {
    PUSH,
    PULL,
    STORE,
    PRINT,
    EXIT,
    MATH_,
    OR,
    AND,
    GREATER_THAN,
    IS_EQUAL_TO,
    LESS_THAN,
    IF_TRUE,
    SCOPE_END
};

enum BC_SUB {
    ADD,
    SUB,
    MULT,
    DIV, // Division as in math
    IS_EQUAL
};

struct Instr {
    itx_types value;
    BC code;
    BC_SUB sub_code;
    int jump = -1;
};


std::vector<Instr> bc_gen(std::vector<IR> ir);
int executor(std::vector<Instr> bc);
#endif //ITX_BC_DEF_HPP