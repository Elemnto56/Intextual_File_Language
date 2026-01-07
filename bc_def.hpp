//
// Created by jaydog on 12/25/25.
//

#ifndef ITX_BC_DEF_HPP
#define ITX_BC_DEF_HPP
#include <any>
#include <vector>

#include "ir_gen.hpp"

enum BC {
    PUSH,
    PULL,
    STORE,
    PRINT,
    EXIT
};

struct Instr {
    LexSubType pull_type;
    itx_types value;
    BC code;
};


std::vector<Instr> bc_gen(std::vector<IR> ir);
void executor(std::vector<Instr> bc);
#endif //ITX_BC_DEF_HPP