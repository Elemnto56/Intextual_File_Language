#include "bc_def.hpp"

using namespace std;
vector<Instr> bc_gen(vector<IR> ir) {
    vector<Instr> bc;

    for (auto& iri : ir) {
        switch (iri.l_type) {
        case VARIABLE:
            bc.push_back({iri.rvalue,PUSH});
            bc.push_back({iri.lvalue, STORE});
            break;
        case KEYWORD:
            if (iri.lvalue == "output") {
                // In this context, MATH should be treated as an expression that involves a variable
                if (iri.r_type != MATH) {bc.push_back({iri.rvalue, PRINT}); break;}

                int i{};
                auto lex_expr = std::get<vector<Lex>>(iri.metadata["var_solve"]);
                while (i < lex_expr.size()) {
                    switch (lex_expr[i].type) {
                    case VARIABLE: bc.push_back({lex_expr[i].value, PULL}); break;
                    case TYPE:
                        if (lex_expr[i].sub_type != BUILT_IN)
                            bc.push_back({lex_expr[i].value, PUSH});
                        else { cerr << "Unexpected value in this statement on line " << lex_expr[i].line << endl; exit(1); }
                        break;
                    case OPERATOR:
                        if ((lex_expr[++i].type == TYPE && lex_expr[i].sub_type != BUILT_IN) || lex_expr[i].type == VARIABLE) bc.push_back({lex_expr[i].value, PUSH});
                        switch (auto op = std::get<string>(lex_expr[i-1].value); op[0]) {
                            case '+': bc.push_back({op, MATH_, ADD}); break;
                            case '-': bc.push_back({op, MATH_, SUB}); break;
                            case '/': bc.push_back({op, MATH_, DIV}); break;
                            default: bc.push_back({op, MATH_, MULT});
                        }
                        break;
                    default: cerr << "Not handled: developer error" << endl; exit(1);
                    }
                    i++;
                }
                bc.push_back({"table", PRINT});
            }
            break;
        default: cerr << "Not handled: developer error" << endl; exit(1);
        }
    }

    return bc;
}