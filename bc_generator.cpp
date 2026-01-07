#include "bc_def.hpp"

using namespace std;
vector<Instr> bc_gen(vector<IR> ir) {
    vector<Instr> bc;
    // Add a symbol table

    for (auto& iri : ir) {
        if (std::get<string>(iri.metadata["wgo"]) == "variable_declare") {
            bc.push_back({.value = iri.rvalue, .code = PUSH});
            bc.push_back({iri.r_type, iri.lvalue, STORE});
        }
        if (iri.l_type == KEYWORD && iri.lvalue == "output") {
            iri.r_type == VAR ? bc.push_back({VAR,iri.rvalue, PRINT}) : bc.push_back({.value = iri.rvalue, .code = PRINT});
        }
    }

    return bc;
}