// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include "gss.h"

int
main(int argc, char **argv) {
    char *err;

    char *input_string = "{\"a\":\"b\",\"c\":[\"d\"]}";
    char *output_string;

    printf("%s\n", input_string);

    err = Convert(input_string, "json", "", "", "yaml", &output_string);

    if (err != NULL) {
        fprintf(stderr, "error: %s\n", err);
        free(err);
        return 1;
    }

    printf("%s\n", output_string);

    return 0;
}
