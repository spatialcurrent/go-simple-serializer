# =================================================================
#
# Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

from ctypes import *
import sys


# Load Shared Object
# gss.so must be in the LD_LIBRARY_PATH
# By default, LD_LIBRARY_PATH does not include current directory.
# You can add current directory with LD_LIBRARY_PATH=. python test.py
lib = cdll.LoadLibrary("gss.so")

# Define Function Definitions
version = lib.Version
version.argtypes = []
version.restype = c_char_p

convert = lib.Convert
convert.argtypes = [c_char_p, c_char_p, c_char_p, c_char_p, c_char_p, POINTER(c_char_p)]
convert.restype = c_char_p

# Define input and output variables
# Output must be a ctypec_char_p
input_string = "{\"a\":\"b\",\"c\":[\"d\"]}"
output_string_pointer = c_char_p()

print "Version:", version

print input_string

err = convert(input_string, "json", "", "", "yaml", byref(output_string_pointer))
if err != None:
    print("error: %s" % (str(err, encoding='utf-8')))
    sys.exit(1)

# Convert from ctype to python string
output_string = output_string_pointer.value

# Print output to stdout
print output_string
