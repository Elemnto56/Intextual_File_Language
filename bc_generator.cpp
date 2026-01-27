#include "bc_def.hpp"

using namespace std;
vector<Instr> bc_gen(vector<IR> ir) {
    vector<Instr> bc;
    vector<int> end_ifs;

    for (auto& iri : ir) {
        switch (iri.l_type) {
        case VARIABLE:
            if (!iri.metadata.contains("no-val")) bc.push_back({iri.rvalue,PUSH});
            bc.push_back({.value = iri.lvalue, .code = STORE, .jump = iri.metadata.contains("re-assign") ? -3 : -1});
            break;
        case KEYWORD:
            if (iri.lvalue == "output")
             // table used a metadata should mean to use the symbol table
             iri.metadata.contains("no-val") ?  bc.push_back({"table", PRINT}) : bc.push_back({iri.rvalue, PRINT});
            break;
        case L_KEYWORD:
            if (iri.lvalue == "if") {
                iri.r_type != BOOL ? bc.push_back({"table", IF_TRUE}) : bc.push_back({iri.rvalue, IF_TRUE});
                end_ifs.push_back(static_cast<int>(bc.size()-1));
            }
            break;
        case BODY:
            for (const auto& in : bc_gen(std::get<vector<IR>>(iri.metadata["body"]))) bc.push_back(in);
            bc[end_ifs.back()].jump = static_cast<int>(bc.size());
            end_ifs.pop_back();
            break;
        case LOGIC: {
            auto lex_expr = std::get<vector<Lex>>(iri.metadata["expression"]);
            for (int i{}; i < lex_expr.size(); i++) {
                if (auto lex = lex_expr[i]; (lex.type == TYPE && lex.sub_type != BUILT_IN) || lex.type == VARIABLE) {
                    if (auto op = !bc.empty() ? bc.back() : Instr{}; op.jump == -2) {
                        bc.pop_back();
                        lex.type == VARIABLE ? bc.push_back({lex.meta, PULL}) : bc.push_back({lex.value, PUSH});
                        bc.push_back(op);
                    }
                    else lex.type == VARIABLE ? bc.push_back({lex.meta, PULL}) : bc.push_back({lex.value, PUSH});
                }
                else if (lex.type == OPERATOR) {
                    if (lex.sub_type == BOOL) {
                        if (auto op = lex.meta; op == "&&") bc.push_back({.code = AND, .jump = -2});
                        else if (op == "||") bc.push_back({.code = OR, .jump = -2});
                    }
                    else if (lex.sub_type == COMPARE) {
                        if (auto op = lex.meta; op == "==") bc.push_back({.code = IS_EQUAL_TO, .jump = -2});
                        else if (op == ">" || op == ">=") bc.push_back({.code = GREATER_THAN, .sub_code = op == ">=" ? IS_EQUAL : ADD, .jump = -2});
                        else if (op == "<" || op == "<=") bc.push_back({.code = LESS_THAN, .sub_code = op == "<=" ? IS_EQUAL : ADD, .jump = -2});
                    }
                    else if (lex.sub_type == MATH)
                        if (auto op = lex.meta; op.size() == 1)
                            switch (op[0]) {
                            case '+': bc.push_back({op, MATH_, ADD, -2}); break;
                            case '-': bc.push_back({op, MATH_, SUB, -2}); break;
                            case '/': bc.push_back({op, MATH_, DIV, -2}); break;
                            default: bc.push_back({op, MATH_, MULT, -2});
                            }
                }
            }
        }
            break;
        default: callErr("Not handled: developer error", iri.line);
        }
    }

    return bc;
}