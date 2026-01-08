#include <iostream>

#include "expr_engine.hpp"
#include "ir_gen.hpp"

using namespace std;

vector<IR> ir_gen(vector<Lex> tokens) {
    vector<IR> all_ir;

    int i=0;
    while (i < tokens.size()) {
        Lex token  = tokens[i];

        if (token.type == TYPE && token.sub_type == BUILT_IN) {
            IR var_declare{};

            if (auto type = get<string>(token.value); type == "int") var_declare.r_type = INT;
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

            vector<Lex> expr;
            while (true) {  i++; if (std::holds_alternative<string>(tokens[i].value) && std::get<string>(tokens[i].value) == ";") break; expr.push_back(tokens[i]); }
            var_declare.rvalue = Expression(expr).solve();

            all_ir.push_back(var_declare);
            if (tokens[i].meta != ";") callErr(MISSING_TOKEN, "Expected an ';' in this variable declaration", tokens[i].line);
        }
        else if (token.type == KEYWORD && std::get<string>(token.value) == "output") {
            i++;
            vector<Lex> expr;
            bool var_include = false;
            while (tokens[i].meta != ";") { if (tokens[i].type == VARIABLE) var_include = true; expr.push_back(tokens[i]); i++;}

            if (!var_include) {
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
                    .metadata = {{"var_solve", expr}},
                    .lvalue = "output",
                    .l_type = KEYWORD,
                    .r_type = MATH,
                    .line = tokens[i].line
                });
            }

            //if (tokens[++i].meta != ";") callErr(MISSING_TOKEN, "Expected an ';' in this print statement", tokens[i].line);
        }

        i++;
    }

    return all_ir;
}