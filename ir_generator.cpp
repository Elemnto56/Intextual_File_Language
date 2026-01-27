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
            else callErr("Invalid type: \""+type+"\"", token.line);

            if (tokens[++i].type == VARIABLE) var_declare.lvalue = tokens[i].meta; else {
                callErr("Not a valid variable name", tokens[i].line);
            }
            var_declare.l_type = VARIABLE;

            if (tokens[++i].type != OPERATOR && tokens[i].meta != "=") callErr("Expected an '=' in this variable declaration", tokens[i].line);

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
                        callErr("The expression in line " + std::to_string(expr.line) + " did not return a boolean\n", tokens[i].line);
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
            else if (token.meta == "while") {
                auto expr = tokens[++i];
                if (!find_var_expr(expr.scope)) {
                    bool pos_bool;
                    try {pos_bool = std::get<bool>(Expression(expr.scope).solve());} catch (...) {
                        callErr("The expression in line " + std::to_string(expr.line) + " did not return a boolean\n", tokens[i].line);
                    }

                    all_ir.push_back({
                        .rvalue = pos_bool,
                        .lvalue = "while",
                        .l_type = L_KEYWORD,
                        .r_type = BOOL,
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
        else if (token.type == VARIABLE ||(token.type == OPERATOR && token.sub_type == MATH)) {
            if (tokens[++i].meta == "=") {
                vector<Lex> expr;
                i++;
                while (tokens[i].meta != ";") {expr.push_back(tokens[i]); i++;}
                if (!find_var_expr(expr))
                    all_ir.push_back({
                        .rvalue = Expression(expr).solve(),
                        .metadata = {{"re-assign", {}}},
                        .lvalue = token.meta,
                        .l_type = VARIABLE,
                        .line = tokens[i].line
                    });
                else {
                    all_ir.push_back({
                        .metadata = {{"expression",  expr}},
                        .l_type = LOGIC,
                        .line = tokens[i].line
                    });
                    all_ir.push_back({
                        .metadata = {{"no-val", {}}, {"re-assign", {}}},
                        .lvalue = token.meta,
                        .l_type = VARIABLE,
                        .line = tokens[i].line
                    });
                }
            }
        }

        i++;
    }

    return all_ir;
}

vector<IR> ir_check(const vector<IR>& ir) {
    int i{};
    std::unordered_map<string, LexSubType> type_table;

    while (i < ir.size()) {
        if (ir[i].l_type == VARIABLE) {
            if (ir[i].metadata.contains("re-assign")) {if (!type_table.contains(ir[i].lvalue)) callErr("Could not find variable " + ir[i].lvalue, ir[i].line); if (type_table[ir[i].lvalue] != ir[i].r_type) callErr("Variable " + ir[i].lvalue + " re-assign did not match the original " + type_look[type_table[ir[i].lvalue]], ir[i].line);}
            else type_table.insert({ir[i].lvalue, ir[i].r_type});
        }
        i++;
    }

    return ir;
}