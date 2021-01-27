package tablifier

import (
	"io"
	"os"
)

type Tablifier struct {
}

var defaultTablifier = Tablifier{}

func (t Tablifier) Tablify(slice interface{}) error {
	return t.Ftablify(os.Stdout, slice)
}

func (t Tablifier) Ftablify(w io.Writer, slice interface{}) error {
	tdata, err := reflectSlice(slice)
	if err != nil {
		return err
	}
	tdata.fprintf(w)
	return nil
}

func Tablify(slice interface{}) error {
	return Ftablify(os.Stdout, slice)
}

func Ftablify(w io.Writer, slice interface{}) error {
	return defaultTablifier.Ftablify(w, slice)
}
