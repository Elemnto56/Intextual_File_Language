#include <vector>
#include <iostream>

#include "bc_def.hpp"

using namespace std;
int executor(vector<Instr> bc) {
    vector<itx_types> stack;
    vector<unordered_map<string, itx_types>> scope{{}};
    int c_scope{};

    cout << boolalpha;

    int i{};
    itx_types operand1, operand2;
    while (i < bc.size()) {
        operand1 = stack.size() == 2 ? stack[0] : "";
        operand2 = stack.size() == 2 ? stack[1] : "";

        switch (auto b = bc[i]; b.code) {
        case PULL: stack.push_back(scope[c_scope][std::get<string>(b.value)]); break;
        case PUSH: stack.push_back(b.value); break;
        case STORE:
            scope[c_scope].insert({std::get<string>(b.value), stack.back()});
            stack.clear();
            break;
        case IF_TRUE:
            if (std::holds_alternative<string>(b.value) && std::get<string>(b.value) == "table") {
                if (!std::holds_alternative<bool>(stack.back())) {cerr << "Expected a boolean but didn't receive one\n"; return  1;}
                if (!std::get<bool>(stack.back())) {i = b.jump; continue;}
            }
            else { if (!std::get<bool>(b.value)) {i = b.jump; continue;} }
            bc.insert(bc.begin()+(b.jump-1), {.code = SCOPE_END});
            scope.emplace_back();
            c_scope++;
            break;
        case MATH_:
            if (std::holds_alternative<string>(operand1) && scope[c_scope].contains(std::get<string>(operand1))) operand1 = scope[c_scope][std::get<string>(operand1)];
            if (std::holds_alternative<string>(operand2) && scope[c_scope].contains(std::get<string>(operand2))) operand2 = scope[c_scope][std::get<string>(operand2)];
            switch (b.sub_code) {
            case ADD:
                    switch (operand1.index()) {
                    case 0:
                        if (!std::holds_alternative<string>(operand2)) {cerr<< "Un handled" << endl; return 1;}
                        stack.clear();
                        stack.emplace_back(std::get<string>(operand1)+std::get<string>(operand2));
                        break;
                    case 1:
                    case 2:
                        stack.clear();
                        if (operand1.index() == operand2.index()) stack.emplace_back((operand1.index() == 1 ? std::get<int>(operand1) : std::get<float>(operand1)) + (operand2.index() == 1 ? std::get<int>(operand2) : std::get<float>(operand2)));
                        else stack.emplace_back(operand1.index() == 1 ? static_cast<float>(std::get<int>(operand1)) + std::get<float>(operand2) : std::get<float>(operand1) + static_cast<float>(std::get<int>(operand2)));
                        break;
                    default:
                        cerr << "Cannot add two bools" << endl;
                        return 1;
                    }
                break;
            case SUB:
                switch (operand1.index()) {
                case 1:
                case 2:
                    stack.clear();
                    if (operand1.index() == operand2.index()) stack.emplace_back((operand1.index() == 1 ? std::get<int>(operand1) : std::get<float>(operand1)) - (operand2.index() == 1 ? std::get<int>(operand2) : std::get<float>(operand2)));
                    else stack.emplace_back(operand1.index() == 1 ? static_cast<float>(std::get<int>(operand1)) - std::get<float>(operand2) : std::get<float>(operand1) - static_cast<float>(std::get<int>(operand2)));
                    break;
                default:
                    cerr << "Either a string or bool was attempted to be subtracted by each other" << endl;
                    return 1;
                }
                break;
            case DIV:
                switch (operand1.index()) {
                case 1:
                case 2:
                    stack.clear();
                    if (operand1.index() == operand2.index()) stack.emplace_back((operand1.index() == 1 ? std::get<int>(operand1) : std::get<float>(operand1)) / (operand2.index() == 1 ? std::get<int>(operand2) : std::get<float>(operand2)));
                    else stack.emplace_back(operand1.index() == 1 ? static_cast<float>(std::get<int>(operand1)) / std::get<float>(operand2) : std::get<float>(operand1) / static_cast<float>(std::get<int>(operand2)));
                    break;
                default:
                    cerr << "Either a string or bool was attempted to be divided by each other" << endl;
                    return 1;
                }
                break;
            default:
                switch (operand1.index()) {
                case 1:
                case 2:
                    stack.clear();
                    if (operand1.index() == operand2.index()) stack.emplace_back((operand1.index() == 1 ? std::get<int>(operand1) : std::get<float>(operand1)) * (operand2.index() == 1 ? std::get<int>(operand2) : std::get<float>(operand2)));
                    else stack.emplace_back(operand1.index() == 1 ? static_cast<float>(std::get<int>(operand1)) * std::get<float>(operand2) : std::get<float>(operand1) * static_cast<float>(std::get<int>(operand2)));
                    break;
                default:
                    cerr << "Either a string or bool was attempted to be multiplied by each other" << endl;
                    return 1;
                }
            }
            break;
        case OR:
        case AND: {
            if (!std::holds_alternative<bool>(stack.front()) || !std::holds_alternative<bool>(stack.back())) {cerr << "Didn't find a boolean in a boolean expression\n"; return 1;}
            bool res = b.code == OR ? std::get<bool>(stack.front()) || std::get<bool>(stack.back()) : std::get<bool>(stack.front()) && std::get<bool>(stack.back());
            stack.clear();
            stack.emplace_back(res);
        }
            break;
        case IS_EQUAL_TO: {
            bool res;
            auto resolve = [stack]() {
                auto front = stack.front();
                auto back = stack.back();
              switch (front.index()) {
              case 0: return std::get<string>(front) == std::get<string>(back);
              case 1: return std::get<int>(front) == std::get<int>(back);
              case 2: return std::get<float>(front) == std::get<float>(back);
              default: return std::get<bool>(front) == std::get<bool>(back);
              }
            };
            if (operand1.index() == 1 && operand2.index() == 2 || operand1.index() == 2 && operand2.index() == 1) res = operand1.index() == 1 && operand2.index() == 2 ? std::get<int>(operand1) == std::get<float>(operand2) : std::get<float>(operand1) == std::get<int>(operand2);
            else if (operand1.index() == operand2.index()) res = resolve();
            else {cerr << "Was not booleans, could not compare\n"; return 1;}

            stack.clear();
            stack.emplace_back(res);
        }
            break;
        case GREATER_THAN: {
            bool res;
            auto resolve = [b, stack]() {
                auto front = stack.front();
                auto back = stack.back();
                switch (front.index()) {
                case 0: return b.sub_code == IS_EQUAL ? std::get<string>(front) >= std::get<string>(back) : std::get<string>(front) > std::get<string>(back);
                case 1: return b.sub_code == IS_EQUAL ? std::get<int>(front) >= std::get<int>(back) : std::get<int>(front) > std::get<int>(back);
                case 2: return b.sub_code == IS_EQUAL ? std::get<float>(front) >= std::get<float>(back) : std::get<float>(front) > std::get<float>(back);
                default: return b.sub_code == IS_EQUAL ? std::get<bool>(front) >= std::get<bool>(back) : std::get<bool>(front) > std::get<bool>(back);
                }
            };
            if (stack.front().index() == 1 || stack.front().index() == 2 &&stack.back().index() == 1 || stack.back().index() == 2)
                if (b.sub_code == IS_EQUAL) res = stack.front().index() == 1 && stack.back().index() == 2 ? std::get<int>(stack.front()) >= std::get<float>(stack.back()) : std::get<float>(stack.front()) >= std::get<int>(stack.back());
                else res = stack.front().index() == 1 && stack.back().index() == 2 ? std::get<int>(stack.front()) > std::get<float>(stack.back()) : std::get<float>(stack.front()) > std::get<int>(stack.back());
            else res = resolve();

            stack.clear();
            stack.emplace_back(res);
        }
            break;
        case LESS_THAN: {
            bool res;
            auto resolve = [b, stack]() {
                auto front = stack.front();
                auto back = stack.back();
                switch (front.index()) {
                case 1: return b.sub_code == IS_EQUAL ? std::get<int>(front) <= std::get<int>(back) : std::get<int>(front) < std::get<int>(back);
                case 2: return b.sub_code == IS_EQUAL ? std::get<float>(front) <= std::get<float>(back) : std::get<float>(front) < std::get<float>(back);
                default: //TODO: call error
                }
            };
            if (stack.front().index() == 1 || stack.front().index() == 2 &&stack.back().index() == 1 || stack.back().index() == 2)
                if (b.sub_code == IS_EQUAL) res = stack.front().index() == 1 && stack.back().index() == 2 ? std::get<int>(stack.front()) <= std::get<float>(stack.back()) : std::get<float>(stack.front()) <= std::get<int>(stack.back());
                else res = stack.front().index() == 1 && stack.back().index() == 2 ? std::get<int>(stack.front()) < std::get<float>(stack.back()) : std::get<float>(stack.front()) < std::get<int>(stack.back());
            else res = resolve();

            stack.clear();
            stack.emplace_back(res);
        }
            break;
        case PRINT:
                if (std::holds_alternative<string>(b.value) && std::get<string>(b.value) == "table") {
                    std::visit([](auto&& val) {cout << val << endl;}, stack.back());
                    stack.clear();
                } else std::visit([](auto&& val) {cout << val << endl;}, b.value);
                break;
        case SCOPE_END:
            scope.pop_back();
            c_scope--;
            break;
        default:
            goto loop_end;
        }

        i++;
    }
    loop_end:;
    return 0;
}