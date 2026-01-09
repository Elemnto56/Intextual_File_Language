#include <vector>
#include <iostream>

#include "bc_def.hpp"

using namespace std;
void executor(vector<Instr> bc) {
    vector<itx_types> stack;
    unordered_map<string, itx_types> symbols;
    cout << boolalpha;

    // TODO: Finalize expressions that use variables
    int i{};
    itx_types operand1, operand2;
    while (i < bc.size()) {
        operand1 = stack.size() == 2 ? stack[0] : "";
        operand2 = stack.size() == 2 ? stack[1] : "";

        switch (auto b = bc[i]; b.code) {
        case PULL: stack.push_back(symbols[std::get<string>(b.value)]); break;
        case PUSH: stack.push_back(b.value); break;
        case STORE:
            symbols.insert({std::get<string>(b.value), stack.back()});
            stack.clear();
            break;
        case MATH_:
            switch (b.sub_code) {
            case ADD:
                    switch (operand1.index()) {
                    case 0:
                        if (!std::holds_alternative<string>(operand2)) {cerr<< "Un handled" << endl; exit(1);}
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
                        exit(1);
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
                    exit(1);
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
                    exit(1);
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
                    exit(1);
                }
            }
            break;
        case PRINT:
                if (std::holds_alternative<string>(b.value) && std::get<string>(b.value) == "table") {
                    std::visit([](auto&& val) {cout << val << endl;}, stack.back());
                    stack.clear();
                } else std::visit([](auto&& val) {cout << val << endl;}, b.value);
                break;
        default:
            goto loop_end;
        }

        i++;
    }
    loop_end:;
}