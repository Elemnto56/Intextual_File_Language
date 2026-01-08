// Ignore this file
#include <iostream>
#include <vector>

template <typename T>
std::vector<std::string> operator<<(const std::vector<T> all, T test) {
    all.push_back(test);
    return all;
}

int main() {
    std::vector<std::string> all;

    all << "Hello";
    for (int i{}; i < all.size(); i++) std::cout << all[i] << std::endl;
}
