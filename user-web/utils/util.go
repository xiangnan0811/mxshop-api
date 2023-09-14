package utils

import "strings"

func RemoveTopStruct(fields map[string]string) map[string]string {
    r := make(map[string]string, len(fields))
    for field, val := range fields {
        r[field[strings.Index(field, ".")+1:]] = val
    }
    return r
}
