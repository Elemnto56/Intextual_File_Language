#include "bc_def.hpp"

using namespace std;
vector<Instr> bc_gen(vector<IR> ir) {
    vector<Instr> bc;
    vector<int> end_ifs;

    for (auto& iri : ir) {
        switch (iri.l_type) {
        case VARIABLE:
            if (!iri.metadata.contains("no-val")) bc.push_back({iri.rvalue,PUSH});
            bc.push_back({iri.lvalue, STORE});
            break;
        case KEYWORD:
            if (iri.lvalue == "output")
             // table used a metadata should mean to use the symbol table
             iri.metadata.contains("no-val") ?  bc.push_back({"table", PRINT}) : bc.push_back({iri.rvalue, PRINT});
            break;
        case L_KEYWORD:
            iri.r_type != BOOL ? bc.push_back({"table", IF_TRUE}) : bc.push_back({iri.rvalue, IF_TRUE});
            end_ifs.push_back(static_cast<int>(bc.size()-1));
            break;
        case BODY:
            for (const auto& in : bc_gen(std::get<vector<IR>>(iri.metadata["body"]))) bc.push_back(in);
            bc[end_ifs.back()].jump = static_cast<int>(bc.size());
            end_ifs.pop_back();
            break;
        case LOGIC: {
            auto lex_expr = std::get<vector<Lex>>(iri.metadata["expression"]);
            for (int i{}; i < lex_expr.size(); i++) {
                if (lex_expr[i].type == TYPE && lex_expr[i].sub_type != BUILT_IN) bc.push_back({lex_expr[i].value, PUSH});
                switch (lex_expr[i].type) {
                case VARIABLE: bc.push_back({lex_expr[i].value, PULL}); break;
                case TYPE:
                    if (lex_expr[i].sub_type != BUILT_IN)
                        bc.push_back({lex_expr[i].value, PUSH});
                    else { cerr << "Unexpected value in this statement on line " << lex_expr[i].line << endl; exit(1); }
                    break;
                case OPERATOR:
                    if (lex_expr[i].sub_type == BOOL) {
                        if (auto op = std::get<string>(lex_expr[i].value); op == "&&") bc.push_back({.code = AND});
                        else if (op == "||") bc.push_back({.code = OR});
                    }
                    else if (lex_expr[i].sub_type == COMPARE) {
                        if (auto op = std::get<string>(lex_expr[i].value); op == "==") bc.push_back({.code = IS_EQUAL_TO});
                        else if (op == ">" || op == ">=") bc.push_back({.code = GREATER_THAN, .sub_code = op == ">=" ? IS_EQUAL : ADD});
                        else if (op == "<" || op == "<=") bc.push_back({.code = LESS_THAN, .sub_code = op == "<=" ? IS_EQUAL : ADD});
                    }
                    else if (lex_expr[i].sub_type == MATH)
                        if (std::get<string>(lex_expr[i].value).size() == 1)
                            switch (auto op = std::get<string>(lex_expr[i].value); op[0]) {
                            case '+': bc.push_back({op, MATH_, ADD}); break;
                            case '-': bc.push_back({op, MATH_, SUB}); break;
                            case '/': bc.push_back({op, MATH_, DIV}); break;
                            default: bc.push_back({op, MATH_, MULT});
                            }

                    break;
                default: cerr << "Not handled: developer error" << endl; exit(1);
                }
            }
        }
            break;
        default: cerr << "Not handled: developer error" << endl; exit(1);
        }
    }

    return bc;
}