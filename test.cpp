// Ignore this file
#include <iostream>
#include <string>

int main() {
    std::string first = "I like // dogs";
    std::string second = "//";

    std::cout << (first.contains(second)) << std::endl;
}
