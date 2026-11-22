package main

/*
#cgo LDFLAGS: -lsqlite3
#include <sqlite3.h>
#include <stdlib.h>
*/
import "C"
import (
	"encoding/hex"
	"log"
	"strings"
	"unsafe"
)

// Encode dynamic text as a SQLite blob literal cast to text. Its alphabet is
// strictly hexadecimal, including for quotes, newlines and embedded NULs.
func q(value string) string { return "CAST(X'" + hex.EncodeToString([]byte(value)) + "' AS TEXT)" }

func (a *App) sql(query string) string {
	path := C.CString(a.db)
	defer C.free(unsafe.Pointer(path))
	var db *C.sqlite3
	code := C.sqlite3_open(path, &db)
	if db != nil {
		defer C.sqlite3_close(db)
	}
	if code != C.SQLITE_OK {
		panic("open database")
	}
	C.sqlite3_busy_timeout(db, 5000)
	source := C.CString("PRAGMA foreign_keys=ON; " + query)
	defer C.free(unsafe.Pointer(source))
	var lines []string
	for next := source; *next != 0; {
		var stmt *C.sqlite3_stmt
		if C.sqlite3_prepare_v2(db, next, -1, &stmt, &next) != C.SQLITE_OK {
			log.Print(C.GoString(C.sqlite3_errmsg(db)))
			panic("prepare")
		}
		if stmt == nil {
			continue
		}
		for {
			code = C.sqlite3_step(stmt)
			if code != C.SQLITE_ROW {
				break
			}
			lines = append(lines, C.GoString((*C.char)(unsafe.Pointer(C.sqlite3_column_text(stmt, 0)))))
		}
		C.sqlite3_finalize(stmt)
		if code != C.SQLITE_DONE {
			log.Print(C.GoString(C.sqlite3_errmsg(db)))
			panic("step")
		}
	}
	return strings.Join(lines, "\n")
}
