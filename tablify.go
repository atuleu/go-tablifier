// package tablifier expose two function to print nice table with
// UTF-8 border for slices of struct.  It uses exposed field name to
// determine the column name, and then print each lines.
package tablifier

import (
	"io"
	"os"
)

// Tablify takes a slice and prints a table to stdout
func Tablify(slice interface{}) error {
	return Ftablify(os.Stdout, slice)
}

// Tablify takes a slice and prints a table to stdout
func Ftablify(w io.Writer, slice interface{}) error {
	tdata, err := reflectSlice(slice)
	if err != nil {
		return err
	}
	tdata.fprintf(w)
	return nil
}
