#include <filesystem>
#include <fstream>
#include <iostream>
#include <vector>

#include "def.hpp"

using namespace std;
int main(int argc, char* argv[]) {
    if (argc > 2 || argc < 2) {cerr << "Too many or too little args provided\n"; return 1;}
    if (!filesystem::exists(argv[1])) {cerr << argv[1] << " does not exist\n"; return 1;}

    vector<string> lines;
    ifstream ifile(argv[1]);
    string line;
    while (getline(ifile, line)) lines.push_back(line);

    int i = 0;
    while (i < lines.size()) {
        if (lines[i].contains("//")) {

        }
        i++;
    }
}