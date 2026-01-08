#include <vector>
#include <iostream>

#include "bc_def.hpp"

using namespace std;
void executor(vector<Instr> bc) {
    vector<itx_types> stack;
    unordered_map<string, itx_types> symbols;
    cout << boolalpha;

    int i = 0;
    while (i < bc.size()) {
        switch (auto b = bc[i]; b.code) {
        case PUSH: stack.push_back(b.value); break;
        case STORE:
            symbols.insert({std::get<string>(b.value), stack.back()});
            stack.clear();
            break;
        case ADD:

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