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
        case PRINT:
            visit([](auto&& value) {cout << value << endl;}, b.pull_type == VAR ? symbols[std::get<string>(b.value)] : b.value);
            break;
        default:
            goto loop_end;
        }

        i++;
    }
    loop_end:;
}