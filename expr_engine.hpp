//
// Created by jaydog on 12/27/25.
//

#ifndef ITX_EXPR_ENGINE_HPP
#define ITX_EXPR_ENGINE_HPP
#include <iostream>
#include <vector>

#include "lex_def.hpp"

using namespace std;

/*  Priority (precedence)
 * 1. Multiplication & Division (* /)
 * 2. Addition & Subtraction (+ -)
 * 3. Comparison (>=, >, <=)
 * 4. Equalities (==, !=)
 * 5. AND (&&)
 * 6. OR (||)
 */

using test = variant<bool>;

class Expression {
private:
    unordered_map<string, itx_types> variables{};
    int e_pos = 0;
    vector<Lex> expr;

    itx_types parseOr() {
        auto left = parseAnd();

        while (e_pos+1 < expr.size() && (expr[e_pos].type == OPERATOR && expr[e_pos].sub_type == BOOL)) {
            if (std::get<string>(expr[e_pos].value) != "||") return left;
            ++e_pos;
            auto right = parseAnd();
            left = std::get<bool>(left) || std::get<bool>(right);
        }

        return left;
    }

    itx_types parseAnd() {
        auto left = parseCompare();

        while (e_pos+1 < expr.size() && (expr[e_pos].type == OPERATOR && expr[e_pos].sub_type == BOOL)) {
            if (std::get<string>(expr[e_pos].value) != "&&") return left;
            ++e_pos;
            auto right = parseCompare();
            left = std::get<bool>(left) && std::get<bool>(right);
        }

        return left;
    }

    itx_types parseCompare() {
        auto left = parseAddnSub();

        if (e_pos+1 < expr.size() && (expr[e_pos+1].type == OPERATOR && expr[e_pos+1].sub_type == COMPARE)) {
            auto operator_ = std::get<string>(expr[e_pos+1].value);
            e_pos += 2;
            auto right = parseAddnSub();
            return compare(left, operator_, right);
        }

        return left;
    }

    itx_types parseAddnSub() {
        auto left = parseMultnDiv();

        while (e_pos+1 < expr.size()) {
            if (expr[e_pos+1].type != OPERATOR && expr[e_pos+1].sub_type != MATH) break;
            auto operator_ = e_pos+1 < expr.size() ? std::get<string>(expr[e_pos+1].value) : "";
            if (operator_ != "+" && operator_ != "-" || operator_ != "-" && operator_ != "+") break;

            e_pos += 2;
            auto right = parseMultnDiv();

            if (operator_ == "+" && left.index() == right.index()) {
                switch (left.index()) {
                case 0: left = std::get<string>(left) + std::get<string>(right); break;
                case 1: left = std::get<int>(left) + std::get<int>(right); break;
                case 2: left = std::get<float>(left) + std::get<float>(right); break;
                    default: break;
                }
            } else if (operator_ == "-" && left.index() == right.index() && (left.index() == 1 || left.index() == 2)) {
                left = left.index() == 1 ? std::get<int>(left) - std::get<int>(right) : std::get<float>(left) - std::get<float>(right);
            }
        }

        return left;
    }

    itx_types parseMultnDiv() {
        auto left = parsePrimary();

        while (e_pos+1 < expr.size()) {
            if (expr[e_pos+1].type != OPERATOR && expr[e_pos+1].sub_type != MATH) break;
            auto operator_ = std::get<string>(expr[e_pos+1].value);
            if (operator_ != "*" && operator_ != "/" || operator_ != "/" && operator_ != "*") break;

            e_pos += 2;
            auto right = parsePrimary();

            if (left.index() == right.index() && (left.index() == 1 || left.index() == 2)) {
                if (operator_ == "*") left.index() == 1 ? left = std::get<int>(left) * std::get<int>(right) : left = std::get<float>(left) * std::get<float>(right);
                if (operator_ == "/") left.index() == 1 ? left = std::get<int>(left) / std::get<int>(right) : left = std::get<float>(left) / std::get<float>(right);
            }
        }

        return left;
    }

    itx_types parsePrimary() {
        if (e_pos >= expr.size()) {
            cerr << "Unexpected end of expression on line " << expr[0].line << endl;
            exit(1);
        }

        return expr[e_pos].value;
    }

    bool compare(itx_types left, string op, itx_types right) {
        // Possibly change this in the future, as other types may be able to compare to each other (i.e. int and bool)
        if (left.index() != right.index()) return false;
        switch (left.index()) {
        case 0:
            if (op == "==") return std::get<string>(left) == std::get<string>(right);
            if (op == "!=") return std::get<string>(left) != std::get<string>(right);
            break;
        case 1:
            if (op == "==") return std::get<int>(left) == std::get<int>(right);
            if (op == "!=") return std::get<int>(left) != std::get<int>(right);
            break;
        case 2:
            if (op == "==") return std::get<float>(left) == std::get<float>(right);
            if (op == "!=") return std::get<float>(left) != std::get<float>(right);
            break;
        case 3:
            if (op == "==") return std::get<bool>(left) == std::get<bool>(right);
            if (op == "!=") return std::get<bool>(left) != std::get<bool>(right);
            break;
        }
    }
public:
    Expression(vector<Lex> expr, std::unordered_map<string, itx_types> vars = {}) : variables(std::move(vars)), expr(std::move(expr)) {}

    itx_types solve() {
        return parseOr();
    }
};
#endif //ITX_EXPR_ENGINE_HPP