package main

import (
    "os"
    "bufio"
    "github.com/json-iterator/go"
    "rdbviewer/back/shared"
)

// custom version of MustImportJson() that reads directly from file
// ram usage is lower, performance is almost the same
// jsoniter is used for import but not for export since it does not correctly indent during export
func MustImportJson(src string, dest interface{}) {
    f,err := os.Open(src)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    buf := bufio.NewReaderSize(f, 1024 * 1024)

    dec := jsoniter.ConfigCompatibleWithStandardLibrary.NewDecoder(buf)
    if err := dec.Decode(dest); err != nil {
        panic(err)
    }
}

func FirstKey(in map[string]struct{}) string {
    for k,_ := range in {
        return k
    }
    return ""
}

func Pagination(curPage uint32, itemsCount uint32, itemsPerPage uint32) (uint32,uint32,bool,bool) {
    if itemsCount == 0 && curPage == 0 {
        return 0, 0, true, true
    }

    start := curPage * itemsPerPage
    if start >= itemsCount {
        return 0, 0, false, false
    }

    end := shared.Uint32Min((curPage + 1) * itemsPerPage, itemsCount)
    FullyLoaded := (end == itemsCount)

    return start, end, FullyLoaded, true
}
