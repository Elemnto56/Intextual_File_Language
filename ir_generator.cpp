#include <iostream>
#include <ranges>

#include "expr_engine.hpp"
#include "ir_gen.hpp"

using namespace std;

vector<IR> ir_gen(vector<Lex> tokens) {
    vector<IR> all_ir;
    int i=0;
    auto find_var_expr = [](vector<Lex>& expr) {
        int j=0;
        while (j < expr.size() && expr[j].meta != ";") { if (expr[j].type == VARIABLE) return true; j++;}
        return false;
    };
    while (i < tokens.size()) {
        if (Lex token = tokens[i]; token.type == TYPE && token.sub_type == BUILT_IN) {
            IR var_declare{};

            if (auto type = token.meta; type == "int") var_declare.r_type = INT;
            else if (type == "string") var_declare.r_type = STRING;
            else if (type == "float") var_declare.r_type = FLOAT;
            else if (type == "bool") var_declare.r_type = BOOL;
            else callErr(BAD_TYPE, "Invalid type: \""+type+"\"", token.line);

            if (tokens[++i].type == VARIABLE) var_declare.lvalue = get<string>(tokens[i].value); else {
                std::cerr << "Not a valid variable name\n line " << tokens[i].line;
                exit(1);
            }
            var_declare.l_type = VARIABLE;

            if (tokens[++i].type != OPERATOR && std::get<string>(tokens[i].value) != "=") callErr(MISSING_TOKEN, "Expected an '=' in this variable declaration", tokens[i].line);

            i++;
            vector<Lex> expr;
            while (tokens[i].meta != ";") { expr.push_back(tokens[i]); i++;}

            if (!find_var_expr(expr)) var_declare.rvalue = Expression(expr).solve();
            else {
                all_ir.push_back({
                    .metadata = {{"expression", expr}},
                    .l_type = LOGIC,
                    .line = tokens[i].line
                });
                var_declare.metadata = {{"no-val", {}}};
            }

            all_ir.push_back(var_declare);
           // if (tokens[i].meta != ";") callErr(MISSING_TOKEN, "Expected an ';' in this variable declaration", tokens[i].line);
        }
        else if (token.type == KEYWORD && token.meta == "output") {
            i++;
            vector<Lex> expr;
            while (tokens[i].meta != ";") { expr.push_back(tokens[i]); i++;}

            if (!find_var_expr(expr)) {
                auto rval = Expression(expr).solve();
                all_ir.push_back({
                    .rvalue = rval,
                    .lvalue = "output",
                    .l_type = KEYWORD,
                    .r_type = static_cast<LexSubType>(rval.index()+1),
                    .line = tokens[i].line
                });
            }
            else {
                all_ir.push_back({
                    .metadata = {{"expression", expr}},
                    .l_type = LOGIC,
                    .line = tokens[i].line
                });
                all_ir.push_back({
                    .metadata = {{"no-val", {}}},
                    .lvalue = "output",
                    .l_type = KEYWORD,
                    .line = tokens[i].line
                });
            }
            //if (tokens[++i].meta != ";") callErr(MISSING_TOKEN, "Expected an ';' in this print statement", tokens[i].line);
        }
        else if (token.type == L_KEYWORD) {
            if (token.meta == "if") {
                auto expr = tokens[++i];
                if (!find_var_expr(expr.scope)) {
                    bool pos_bool;
                    try {pos_bool = std::get<bool>(Expression(expr.scope).solve());} catch (...) {
                        std::cerr << "The expression in line " << expr.line << " did not return a boolean\n"; exit(1);
                    }

                    all_ir.push_back({
                        .rvalue = pos_bool,
                        .lvalue = "if",
                        .l_type = L_KEYWORD,
                        .r_type = BOOL,
                        .line = expr.line
                    });
                }
                else {
                    all_ir.push_back({
                        .metadata = {{"expression", expr.scope}},
                        .l_type = LOGIC,
                        .line = expr.line
                    });
                    all_ir.push_back({
                        .lvalue = "if",
                        .l_type = L_KEYWORD,
                        .line = expr.line
                    });
                }
            }
        }
        else if (token.type == BODY) {
            all_ir.push_back({
                .metadata = {{"body", ir_gen(token.scope)}},
                .l_type = BODY,
                .line = token.line,
            });
        }

        i++;
    }

    return all_ir;
}