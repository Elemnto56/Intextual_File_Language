// Ignore this file
#include <iostream>
#include <string>

int main() {
    std::string first = "I like / dogs";
    std::string second = "//";

    std::cout << (static_cast<int>(first.find(second)) > static_cast<int>(first.size()) ? -1 : first.find(second)) << std::endl;
}
