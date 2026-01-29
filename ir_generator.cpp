#include <iostream>
#include <ranges>
#include <boost/predef/make.h>

#include "expr_engine.hpp"
#include "ir_gen.hpp"

using namespace std;

vector<IR> ir_gen(vector<Lex> tokens) {
    vector<IR> all_ir;
    int i=0;
    auto findRuntimeVar_expr = [](vector<Lex>& expr) {
        int j=0;
        while (j < expr.size() && expr[j].meta != ";") { if (expr[j].type == VARIABLE && !compile_time_constants.contains(expr[j].meta)) return true; j++;}
        return false;
    };

    while (i < tokens.size()) {
        if (Lex token = tokens[i]; token.type == TYPE && (token.sub_type == BUILT_IN || token.sub_type == COMPILE_CONST)) {
            IR var_declare{};

            if (auto type = token.meta; type == "int") var_declare.r_type = INT;
            else if (type == "string") var_declare.r_type = STRING;
            else if (type == "float") var_declare.r_type = FLOAT;
            else if (type == "bool") var_declare.r_type = BOOL;
            else callErr("Invalid type: \""+type+"\"", token.line);

            if (tokens[++i].type == VARIABLE) var_declare.lvalue = tokens[i].meta; else callErr("Not a valid variable name", tokens[i].line);
            var_declare.l_type = VARIABLE;

            if (tokens[++i].type != OPERATOR && tokens[i].meta != "=") callErr("Expected an '=' in this variable declaration", tokens[i].line);

            i++;
            vector<Lex> expr;
            while (tokens[i].meta != ";") { expr.push_back(tokens[i]); i++;}

            if (!findRuntimeVar_expr(expr)) {
                var_declare.rvalue = Expression(expr).solve();
                if (token.sub_type == COMPILE_CONST) {compile_time_constants.insert({var_declare.lvalue, var_declare.rvalue});}
            }
            else {
                all_ir.push_back({
                    .metadata = {{"expression", expr}},
                    .l_type = LOGIC,
                    .line = tokens[i].line
                });
                var_declare.metadata = {{"no-val", {}}};
            }

            if (token.sub_type != COMPILE_CONST) all_ir.push_back(var_declare);
           // if (tokens[i].meta != ";") callErr(MISSING_TOKEN, "Expected an ';' in this variable declaration", tokens[i].line);
        }
        else if (token.type == KEYWORD && token.meta == "output") {
            i++;
            vector<Lex> expr;
            while (tokens[i].meta != ";") { expr.push_back(tokens[i]); i++;}

            if (!findRuntimeVar_expr(expr)) {
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
            if (bool elif = token.meta == "elseif"; token.meta == "if" || elif) {
                if (auto expr = tokens[++i]; !findRuntimeVar_expr(expr.scope)) {
                    bool pos_bool;
                    try {pos_bool = std::get<bool>(Expression(expr.scope).solve());} catch (...) {callErr("The expression in line " + std::to_string(expr.line) + " did not return a boolean\n", tokens[i].line);}

                    all_ir.push_back({
                        .rvalue = pos_bool,
                        .lvalue = elif ? "else_if" : "if",
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
                        .lvalue = elif ? "else_if" : "if",
                        .l_type = L_KEYWORD,
                        .line = expr.line
                    });
                }
            }
            else if (token.meta == "else")
                all_ir.push_back({
                    .rvalue = true,
                    .lvalue = "else",
                    .l_type = L_KEYWORD,
                    .r_type = BOOL,
                    .line = tokens[i].line
                });
            else if (token.meta == "while") {
                auto expr = tokens[++i];
                if (!findRuntimeVar_expr(expr.scope)) {
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
                else {
                    all_ir.push_back({
                        .metadata = {{"expression", expr.scope}, {"include_expr_start", {}}},
                        .l_type = LOGIC,
                        .line = expr.line
                    });
                    all_ir.push_back({
                        .lvalue = "while",
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
        else if (token.type == VARIABLE ||(token.type == OPERATOR && token.sub_type == MATH)) {
            if (tokens[++i].meta == "=") {
                vector<Lex> expr;
                i++;
                while (tokens[i].meta != ";") {expr.push_back(tokens[i]); i++;}
                if (!findRuntimeVar_expr(expr))
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
            else if (contains(tokens[i].meta, {"+=", "-=", "*=", "/=", "%="})) {
                Lex mathOp = {.sub_type = MATH, .meta = "+", .type = OPERATOR, .line = token.line};
                if (tokens[i].meta == "-=") mathOp.meta = "-";
                else if (tokens[i].meta == "*=") mathOp.meta = "*";
                else if (tokens[i].meta == "/=") mathOp.meta = "/";
                else if (tokens[i].meta == "%=") mathOp.meta = "%";
                all_ir.push_back({
                    .metadata = {{"expression",  vector{
                        token,
                        mathOp,
                        tokens[++i]
                    }}},
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

        i++;
    }

    return all_ir;
}
//TODO
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