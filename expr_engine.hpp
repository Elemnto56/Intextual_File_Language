//
// Created by jaydog on 12/27/25.
//

#ifndef ITX_EXPR_ENGINE_HPP
#define ITX_EXPR_ENGINE_HPP
#include <cmath>
#include <functional>
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
inline std::unordered_map<std::string, itx_types> compile_time_constants;

class Expression {
private:
    std::function<LexSubType(size_t)> valToType = [](size_t index) {
        switch(index) {
            case 0: return STRING;
            case 1: return INT;
            case 2: return FLOAT;
            default: return  BOOL;
        }
    };
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
            if (expr[e_pos].meta != "&&") return left;
            ++e_pos;
            auto right = parseCompare();
            left = std::get<bool>(left) && std::get<bool>(right);
        }

        return left;
    }

    itx_types parseCompare() {
        auto left = parseAddnSub();

        if (e_pos+1 < expr.size() && (expr[e_pos+1].type == OPERATOR && expr[e_pos+1].sub_type == COMPARE)) {
            auto operator_ = expr[e_pos+1].meta;
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
            auto operator_ = e_pos+1 < expr.size() ? expr[e_pos+1].meta : "";
            if (operator_ != "+" && operator_ != "-") break;

            e_pos += 2;
            auto right = parseMultnDiv();

            if (std::holds_alternative<int>(left)&&std::holds_alternative<float>(right)) left = static_cast<float>(std::get<int>(left));
            if (std::holds_alternative<int>(right)&&std::holds_alternative<float>(left)) right = static_cast<float>(std::get<int>(right));
            if (left.index() == right.index()) {
                if (operator_ == "-") {
                    if (!std::holds_alternative<string>(left)&&!std::holds_alternative<string>(right)) switch (left.index()) {
                    case 1: left = std::get<int>(left) + std::get<int>(right); break;
                    case 2: left = std::get<float>(left) + std::get<float>(right); break;
                    default: callErr("Cannot subtract two booleans or with another type", expr[e_pos].line);
                    } else callErr("Cannot subtract two strings or with another type", expr[e_pos].line);
                } else {
                    if (!std::holds_alternative<bool>(left)&&!std::holds_alternative<bool>(right))
                        switch (left.index()) {
                        case 0: left = std::get<string>(left) + std::get<string>(right); break;
                        case 1: left = std::get<int>(left) + std::get<int>(right); break;
                        case 2: left = std::get<float>(left) + std::get<float>(right); break;
                        }
                    else callErr("Cannot add two booleans or with another type", expr[e_pos].line);
                }
            } else callErr("Cannot perform math on two different types (unless it's an int and float)", expr[e_pos].line);
        }

        return left;
    }

    itx_types parseMultnDiv() {
        auto left = parsePrimary();

        while (e_pos+1 < expr.size()) {
            if (expr[e_pos+1].type != OPERATOR && expr[e_pos+1].sub_type != MATH) break;
            auto operator_ = expr[e_pos+1].meta;
            if (operator_ != "*" && operator_ != "/" && operator_ != "%") break;

            e_pos += 2;
            auto right = parsePrimary();

            if (left.index() == right.index() && (left.index() == 1 || left.index() == 2)) {
                if (operator_ == "*") left.index() == 1 ? left = std::get<int>(left) * std::get<int>(right) : left = std::get<float>(left) * std::get<float>(right);
                if (operator_ == "/") left.index() == 1 ? left = std::get<int>(left) / std::get<int>(right) : left = std::get<float>(left) / std::get<float>(right);
                if (operator_ == "/") left.index() == 1 ? left = std::get<int>(left) % std::get<int>(right) : left = std::fmod(std::get<float>(left), std::get<float>(right));
            }
        }

        return left;
    }

    itx_types parsePrimary() {
        if (e_pos >= expr.size()) {
            callErr("Empty expression", !expr.empty() ? expr[0].line : -1);
            exit(1);
        }

        if (expr[e_pos].type == VARIABLE) {
            expr[e_pos].value = compile_time_constants[expr[e_pos].meta];
            expr[e_pos].sub_type = valToType(expr[e_pos].value.index());
            expr[e_pos].type = TYPE;
        }

        if (expr[e_pos].meta == "expression") {

            expr[e_pos].value = Expression(expr[e_pos].scope).solve();
            expr[e_pos].sub_type = valToType(expr[e_pos].value.index());

            if (e_pos+1 < expr.size() && expr[e_pos+1].type != OPERATOR && expr[e_pos+1].type == TYPE && expr[e_pos+1].sub_type != BUILT_IN) expr.insert(expr.begin()+(e_pos+1), {.sub_type = MATH, .value = "*", .type = OPERATOR, .line = expr[e_pos].line});
            if (e_pos-1 >= 0 && expr[e_pos-1].type != OPERATOR && expr[e_pos-1].type == TYPE && expr[e_pos-1].sub_type != BUILT_IN) expr.insert(expr.begin()+(e_pos-1), {.sub_type = MATH, .value = "*", .type = OPERATOR, .line = expr[e_pos].line});
        }

        return expr[e_pos].value;
    }

    bool compare(itx_types left, string op, itx_types right) {
        return std::visit([op](auto&& l, auto&& r) {
            using lt = std::decay_t<decltype(l)>;
            using rt = std::decay_t<decltype(r)>;

            if constexpr (is_same_v<lt, int> && is_same_v<rt, float>) l = static_cast<float>(l);
            if constexpr (is_same_v<rt, int> && is_same_v<lt, float>) r = static_cast<float>(r);

            if constexpr (std::is_same_v<lt, rt>) {
                if (op == "==") return l == r;
                if (op == "!=") return l != r;
                if (op == ">=") return l >= r;
                if (op == "<=") return l <= r;
                if (op == ">") return l > r;
                if (op == "<") return l < r;
            }

            return false;
        }, left, right);
    }
public:
    Expression(vector<Lex> expr) : expr(std::move(expr)) {}

    itx_types solve() {
        return parseOr();
    }
};
#endif //ITX_EXPR_ENGINE_HPP