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
        case PULL: stack.push_back(scope[c_scope][get<string>(b.value)]); break;
        case PUSH: stack.push_back(b.value); break;
        case STORE:
            if (scope[c_scope].contains(std::get<string>(b.value))) scope[c_scope][std::get<string>(b.value)] = stack.back();
            else scope[c_scope].insert({get<string>(b.value), stack.back()});
            stack.clear();
            break;
        case IF_TRUE:
            if (holds_alternative<string>(b.value) && get<string>(b.value) == "table") {
                if (!holds_alternative<bool>(stack.back())) {cerr << "Expected a boolean but didn't receive one\n"; return  1;}
                if (!get<bool>(stack.back())) {i = b.jump; continue;}
            }
            else { if (!get<bool>(b.value)) {i = b.jump; continue;} }
            bc.insert(bc.begin()+(b.jump-1), {.code = SCOPE_END});
            scope.emplace_back(scope[c_scope]);
            c_scope++;
            break;
        case MATH_:
            if (holds_alternative<string>(operand1) && scope[c_scope].contains(get<string>(operand1))) operand1 = scope[c_scope][get<string>(operand1)];
            if (holds_alternative<string>(operand2) && scope[c_scope].contains(get<string>(operand2))) operand2 = scope[c_scope][get<string>(operand2)];
            
            visit([b, &stack](auto&& op1, auto&& op2) {
                using L = decay_t<decltype(op1)>;
                using R = decay_t<decltype(op2)>;
                if constexpr (is_same_v<L, int> && is_same_v<R, float>) op1 = static_cast<float>(op1);
                if constexpr (is_same_v<R, int> && is_same_v<L, float>) op2 = static_cast<float>(op2);
                
                stack.clear();
                if constexpr (is_same_v<L, R>)
                switch (b.sub_code) {
                case ADD: stack.emplace_back(op1 + op2); break;
                case SUB: 
                    if constexpr (is_same_v<L, string>) callErr("Could not perform a subtraction between two strings", -1);
                    else stack.emplace_back(op1 - op2); 
                    break;
                case DIV:
                    if constexpr (is_same_v<L, string>) callErr("Could not perform division between two strings", -1);
                    else stack.emplace_back(op1 / op2);
                    break;
                default:
                    if constexpr (is_same_v<L, string>) callErr("Could not perform multiplication between two strings", -1);
                    else stack.emplace_back(op1 * op2);
                    break;
                }
            }, operand1, operand2);
                break;
        case OR:
        case AND:
            if (!holds_alternative<bool>(stack.front()) || !holds_alternative<bool>(stack.back())) callErr("Didn't find a boolean in a boolean expression\n", -1);
            stack.clear();
            stack.emplace_back(b.code == OR ? get<bool>(stack.front()) || get<bool>(stack.back()) : get<bool>(stack.front()) && get<bool>(stack.back()));
            break;
        case IS_EQUAL_TO: {
            auto res = visit([](auto&& front, auto&& back) {
                using L = decay_t<decltype(&front)>;
                using R = decay_t<decltype(&back)>;
                if constexpr (is_same_v<L, int> && is_same_v<R, float>) front = static_cast<float>(front);
                if constexpr (is_same_v<R, int> && is_same_v<L, float>) back = static_cast<float>(back);
                if constexpr (is_same_v<L, R>) return front == back;
                else return false;
            }, stack.front(), stack.back());

            stack.clear();
            stack.emplace_back(res);
        }
            break;
        case GREATER_THAN: {
            bool res = visit([b](auto&& front, auto&& back) {
                using L = decay_t<decltype(&front)>;
                using R = decay_t<decltype(&back)>;
                if constexpr (is_same_v<L, int> && is_same_v<R, float>) front = static_cast<float>(front);
                if constexpr (is_same_v<R, int> && is_same_v<L, float>) back = static_cast<float>(back);
                
                if constexpr (!is_same_v<L, float> || !is_same_v<R, float>) callErr("Could not do greater than (>/>=) between non-floats or non-ints", -1); else
                    if (b.sub_code != IS_EQUAL) return front > back; else return front >= back;
                        return false; // Yes, I know it's unreachable, but it prevents an error
            }, stack.front(), stack.back());
            stack.clear();
            stack.emplace_back(res);
        }
            break;
        case LESS_THAN: {
            bool res = visit([b](auto&& front, auto&& back) {
                using L = decay_t<decltype(&front)>;
                using R = decay_t<decltype(&back)>;
                if constexpr (is_same_v<L, int> && is_same_v<R, float>) front = static_cast<float>(front);
                if constexpr (is_same_v<R, int> && is_same_v<L, float>) back = static_cast<float>(back);
                
                if constexpr (!is_same_v<L, float> || !is_same_v<R, float>) callErr("Could not do less than (</<=) between non-floats or non-ints", -1); else
                if (b.sub_code != IS_EQUAL) return front < back; else return front <= back;
                    return false;
            }, stack.front(), stack.back());
            stack.clear();
            stack.emplace_back(res);
        }
            break;
        case PRINT:
                if (holds_alternative<string>(b.value) && get<string>(b.value) == "table") {
                    visit([](auto&& val) {cout << val << endl;}, stack.back());
                    stack.clear();
                } else visit([](auto&& val) {cout << val << endl;}, b.value);
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