// Ignore this file
#include <bits/stdc++.h>
#include <iostream>
#include <vector>

int main() {
    std::vector x = {1, 2, 3};
    std::rotate(x.begin()+1, x.begin()+2, x.end());
    for (const int i : x) std::cout << i << std::endl;
}
