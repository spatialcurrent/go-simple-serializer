// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

#include <iostream>
#include <string>
#include <cstring>
#include "gss.h"

// Conv is a example of a C++ function that can convert between formats using some std::string variables.
// In production, you would want to write the function definition to match the use case.
char* conv(std::string input_string, std::string input_format, std::string input_header, std::string input_comment, std::string output_format, char** output_string_c) {

  char *input_string_c = new char[input_string.length() + 1];
  std::strcpy(input_string_c, input_string.c_str());
  char *input_format_c = new char[input_format.length() + 1];
  std::strcpy(input_format_c, input_format.c_str());
  char *input_header_c = new char[input_header.length() + 1];
  std::strcpy(input_header_c, input_header.c_str());
  char *input_comment_c = new char[input_comment.length() + 1];
  std::strcpy(input_comment_c, input_comment.c_str());
  char *output_format_c = new char[output_format.length() + 1];
  std::strcpy(output_format_c, output_format.c_str());

  char *err = Convert(input_string_c, input_format_c, input_header_c, input_comment_c, output_format_c, output_string_c);

  free(input_string_c);
  free(input_format_c);
  free(input_header_c);
  free(input_comment_c);
  free(output_format_c);

  return err;

}

int main(int argc, char **argv) {

  // Since Go requires non-const values, we must define our parameters as variables
  // https://stackoverflow.com/questions/4044255/passing-a-string-literal-to-a-function-that-takes-a-stdstring
  std::string input_string("{\"a\":\"b\",\"c\":[\"d\"]}");
  std::string input_format("json");
  std::string input_header("");
  std::string input_comment("");
  std::string output_format("yaml");
  char *output_char_ptr;

  // Write version to stderr
  std::cout << "Version" << Version() <<std::endl;

  // Write input to stderr
  std::cout << input_string << std::endl;

  char *err = conv(input_string, input_format, input_header, input_comment, output_format, &output_char_ptr);
  if (err != NULL) {
    // Write output to stderr
    std::cerr << std::string(err) << std::endl;
    // Return exit code indicating error
    return 1;
  }
  std::string output_string = std::string(output_char_ptr);

  // Write output to stdout
  std::cout << output_string << std::endl;

  // Return exit code indicating success
  return 0;
}
