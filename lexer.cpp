#include <iostream>
#include <boost/regex.hpp>
#include "ir_gen.hpp"
#include "lex_def.hpp"

using namespace std;
bool contains(string test, vector<string> finds) {
    for (string str : finds) if (test == str) return true;
    return false;
}

vector<string> split_by_ws(string test) {
    vector<string> res;
    auto b_itr = boost::sregex_iterator(test.begin(), test.end(),  boost::regex(R"([^\s|\(]+)"));
    while (b_itr != boost::sregex_iterator()) {res.push_back(b_itr->str()); ++b_itr;}
    return res;
}

vector<Lex> lexer(vector<string> lines, char stop_at) {
    std::unordered_map<string, string> alias{{"print", "output"}};
    vector<Lex> allTokens{};
    int j;
    int i = j = 0;

    while (i < lines.size()) {
        string line= lines[i];
        if (line.contains("//")) line = boost::regex_replace(line, boost::regex(R"(\/\/.+)"), "");
        if (line.empty() || ranges::all_of(line, [](unsigned char c){return isspace(c);})) { i++; continue; }
        line.erase(line.begin(), ranges::find_if_not(line, [](unsigned char ch) {return isspace(ch);}));

        string pos_keyword;
        vector keywords = split_by_ws(line);
        if (keywords[0] == "!alias") { alias.insert({keywords[2], keywords[1]}); i++; continue;}

        pos_keyword = alias.contains(keywords[0]) ? pos_keyword = alias[keywords[0]] : pos_keyword = keywords[0];
        if (contains(pos_keyword, {"while", "if", "else if", "else"})) {
            allTokens.push_back({
                .meta = pos_keyword,
                .type = L_KEYWORD,
                .line = i+1
            });

            if (alias.contains(keywords[0])) pos_keyword = keywords[0];
            j += static_cast<int>(pos_keyword.size());
        }

        // Future note when implementing arrays: for each element, run it through the lexer and put the result in that
        // array's value. The parser should be aware of that and parse accordingly.
        else if (contains(pos_keyword, {"int", "string", "bool", "arr", "float"})) {
            allTokens.push_back({
                .meta = pos_keyword,
                .type = TYPE,
                .line = i+1
            });
            j += static_cast<int>(pos_keyword.size());
        }

        else if (contains(pos_keyword, {"output"})) {
            allTokens.push_back({
                .meta = pos_keyword,
                .type = KEYWORD,
                .line = i+1
            });
            j += static_cast<int>(pos_keyword.size());
        }

        while (j < line.size()) {
            char ch = line[j];

            if (ch == ' ') { j++; continue;}

            if (j+1 < line.size() && line.substr(j, 2) == "->") {
                allTokens.push_back({
                    .meta = line.substr(j, 2),
                    .type = SYMBOL,
                    .line = i+1
                });
                j += 1;
                continue;
            }

            if (j+1 < line.size() && contains(line.substr(j, 2), {">=", "<=", "==", "||", "&&", "!="})) {
                line.substr(j, 2) == "&&" || line.substr(j, 2) == "||" ?
                allTokens.push_back({
                    .sub_type = BOOL,
                    .meta = line.substr(j, 2),
                    .type = OPERATOR,
                    .line = i+1,
                })
                :
                allTokens.push_back({
                    .sub_type = COMPARE,
                    .meta = line.substr(j, 2),
                    .type = OPERATOR,
                    .line = i+1
                });

                j += 2;
                continue;
            }

            if (j+1 < line.size() && contains(line.substr(j, 2), {"+=", "-=", "*=", "/=", "%=", "++", "--"})) {
                allTokens.push_back({
                    .sub_type = MATH,
                    .meta = line.substr(j, 2),
                    .type = OPERATOR,
                    .line = i+1
                });
                j += 1;
                continue;
            }


            // Multi-line comment support /*
            if (j+1 < line.size() && line.substr(j, 2) == "/*") {

                while (i < lines.size()) {
                    line = lines[i];

                    j = 0;
                    while (j < line.size()) {
                        if (j+1 < line.size() && line.substr(j, 2) == "*/") {
                            j += 2;
                            line = line.substr(j, j-line.size());
                            j = 0; // Reset the iteration since a new line is initialized
                            goto capture_end;;
                        }
                        j++;
                    }
                    i++;
                }
                capture_end:;

                continue;
            }

            if (isdigit(ch)) {
                string num_capture;
                num_capture += num_capture;

                bool is_float = false;

                while (j < line.size()) {
                    if (!isdigit(line[j]) && line[j] != '.') break;
                    if (line[j] == '.') is_float = true;
                    num_capture += line[j];
                    j++;
                }

                if (is_float) {
                    allTokens.push_back({
                        .sub_type = FLOAT,
                        .value = stof(num_capture),
                        .type = TYPE,
                        .line = i+1
                    });
                    continue;
                }

                allTokens.push_back({
                    .sub_type = INT,
                    .value = stoi(num_capture),
                    .type = TYPE,
                    .line = i+1
                });
                continue;
            }

            if (j+4 < line.size() && line.substr(j, 5) == "false" || line.substr(j, 4) == "true") {
                if (line.substr(j, 5) == "false") {
                    allTokens.push_back({
                        .sub_type = BOOL,
                        .value = false,
                        .type = TYPE,
                        .line = i+1
                    });
                    j += 5;
                }
                else if (line.substr(j, 4) == "true") {
                    allTokens.push_back({
                        .sub_type = BOOL,
                        .value = true,
                        .type = TYPE,
                        .line = i+1
                    });
                    j += 4;
                }

                continue;
            }

            // Variable capturing -- Stuff like booleans and keywords should be captured before
            if (boost::regex_search(std::string{ch}, boost::regex("\\w"))) {
                string var_name;

                while (boost::regex_search(std::string{line[j]}, boost::regex("\\w"))) {
                    var_name += line[j];
                    j++;
                }
                allTokens.push_back({
                    .meta = var_name,
                    .type = VARIABLE,
                    .line = i+1
                });

                continue;
            }

            switch (ch) {
            case '+':
            case '-':
            case '/':
            case '*':
                allTokens.push_back({
                    .sub_type = MATH,
                    .meta = std::string{ch},
                    .type = OPERATOR,
                    .line = i+1
                });
                j++;
                continue;
            case '>':
            case '<':
                allTokens.push_back({
                    .sub_type = COMPARE,
                    .meta = std::string{ch},
                    .type = OPERATOR,
                    .line = i+1
                });
            case '=':
                allTokens.push_back({
                    .meta = "=",
                    .type = OPERATOR,
                    .line = i+1
                });
                j++;
                continue;
            case '"': {
                j++; // Skip the "
                string res;

                while (j < line.size()) {
                    if (line[j] == '"') break;

                    if (string escape = line.substr(j, 2); escape == "\\n")  {res += "\n"; j += 2;}
                    else if (escape == "\\r") { res += "\r"; j+=2; }
                    else if (escape == "\\\"") { res += "\""; j += 2;}
                    else if (escape == "\\\\") { res += "\\"; j += 2;}
                    else if (escape == "\\t") { res += "\t"; j += 2;}


                    if (line[j] == '"') break; // A second one was included because the next j iteration may range out and not capture the last character
                    res += line[j];
                    j++;
                }

                if (j >= line.size()) callErr("string ranged out", i+1);

                allTokens.push_back({
                    .sub_type = STRING,
                    .value = res,
                    .type = TYPE,
                    .line = i+1
                });
                j++; // Skip the last "
                continue;
            }
            case '\'': {
                j++;
                char user_ch = line[j];
                if (line[++j] != '\'') callErr("Malformed single character", i+1);
                allTokens.push_back({
                    .sub_type = STRING,
                    .value = std::string{user_ch},
                    .type = TYPE,
                    .line = i+1
                });
                j++;
                continue;
            }
            case '(': {
                j++;
                vector str_expr = {line.substr(j, line.size())};
                for (const auto& l : lines) str_expr.push_back(l); //NOTE: This may need to be refactored in the future
                auto expr_lex = lexer(str_expr, ')');
                j += expr_lex.back().line;
                expr_lex.pop_back();
                allTokens.push_back({
                    .sub_type = BOOL,
                    .meta = "expression",
                    .type = TYPE,
                    .scope = expr_lex,
                    .line = i+1
                });
            }
            j++;
            continue;
            case '{': {
                j++;
                string half_line;
                while (j < line.size()) {half_line += line[j]; j++;}
                vector<Lex> end_line;
                auto searchScopeEnd = [&end_line](vector<Lex> l) {
                    for (int c_i{}; c_i < l.size(); c_i++) if (l[c_i].sub_type == R_CURL) {
                        c_i++;
                        while (c_i < l.size()) { end_line.push_back(l[c_i]); c_i++; }
                        return true;
                    }
                    return false;
                };

                lines.insert(lines.begin()+ (++i), half_line);
                vector<Lex> lex_lines, line_;
                while (i < lines.size()) {
                    line_ = lexer(vector{lines[i]});
                    for (const auto& lex : line_) lex_lines.push_back(lex);
                    if (searchScopeEnd(line_)) break;
                    i++;
                }
                lex_lines.pop_back(); // remove }
                allTokens.push_back({
                    .type = BODY,
                    .scope = lex_lines,
                    .line = i+1
                });
                for (const auto& lex : end_line) allTokens.push_back(lex);
            }
                goto cont_outer;
            case ';':
                allTokens.push_back({
                    .meta = ";",
                    .value = std::string{ch},
                    .type = SYMBOL,
                    .line = i+1
                });
                j++;
                continue;
            case '}':
                if (stop_at == '}') {
                    allTokens.push_back({.line = j});
                    goto end;
                }
                allTokens.push_back({
                    .sub_type = R_CURL,
                    .value = std::string{ch},
                    .type = SYMBOL,
                    .line = i+1
                });
                j++;
                continue;
            case ')':
                // The last token before closure has only a line field equal to j in order to correctly parse the end of expressions
                if (stop_at == ')') {allTokens.push_back({.line = j});goto end;}
                allTokens.push_back({
                    .sub_type = R_PARA,
                    .value = std::string{ch},
                    .type = SYMBOL,
                    .line = i+1
                });
                j++;
                continue;
            default:
                allTokens.push_back({
                    .meta = string{ch},
                    .type = OPERATOR,
                    .line = i+1
                });
                j++;
            }
        }
        cont_outer:;
        i++;
        j = 0;
    }
    end:;

    return allTokens;
}