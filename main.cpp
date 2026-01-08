#include <complex>
#include <fstream>
#include <iostream>

#include "bc_def.hpp"
#include "expr_engine.hpp"
#include "ir_gen.hpp"
#include "lex_def.hpp"

#define DEBUG 0

int main(int argc, char* argv[]) {
    if (argc > 2 || argc < 2) {std::cerr << "Too many or too little args provided\n"; return 1;}
    if (!std::filesystem::exists(argv[1])) {std::cerr << argv[1] << " does not exist\n"; return 1;}

    std::vector<std::string> lines{};
    std::ifstream ifile(argv[1]);
    std::string token;
    while (getline(ifile, token)) lines.push_back(token);

    auto lex = lexer(lines);
    auto ir = ir_gen(lex);
    auto b = bc_gen(ir);

#if DEBUG
    std::cout << boolalpha;
    /*
    visit([](auto&& val){std::cout << val << std::endl;}, Expression(vector<Lex>{
        Expression testing
    }).solve());
    */

    for (auto l : lex) {
        std::cout << "====================\n";
        std::cout << "LINE: " << l.line << std::endl;
        std::cout << "TYPE: " << l.type << std::endl;
        std::cout << "SUB-TYPE: " << l.sub_type << std::endl;
        std::visit([](auto&& value){std::cout << value << std::endl;}, l.value);
        std::cout << "====================\n";
    }

    std::cout << "=====IR=====" << std::endl;

    for (auto i : ir) {
        std::cout << "line: " << i.line << std::endl;
        std::cout << "l-type: " << i.l_type;
        std::cout << " | l-value: " << i.lvalue << std::endl;
        std::cout << "r-type: " << i.r_type;
        std::visit([](auto&& rval) {std::cout << " | r-value: " << rval << std::endl;}, i.rvalue);
        std::cout << "====================\n";
    }

    std::cout << "====Bytecode====" << std::endl;

    for (auto bc : b) {
        std::cout << bc.code << " " << bc.pull_type << " ";
        std::visit([](auto&& val) {std::cout << val << std::endl;}, bc.value);
    }
#else
    executor(b);
#endif
}