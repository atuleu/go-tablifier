package tablifier

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type tableData struct {
	columns     []string
	columnsSize []int
	lines       [][]string
}

func checkType(slice interface{}) (reflect.Type, error) {
	sType := reflect.TypeOf(slice)
	if sType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("tablifier: argument is not a slice of struct")
	}
	elemType := sType.Elem()
	if elemType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("tablifier: argument is not a slice of struct")

	}
	return elemType, nil
}

func (d *tableData) parseColumns(elemType reflect.Type) {
	d.columns = make([]string, elemType.NumField())
	d.columnsSize = make([]int, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		fieldType := elemType.Field(i)
		columnName := fieldType.Name
		nameTag := fieldType.Tag.Get("name")
		if len(nameTag) > 0 {
			columnName = nameTag
		}
		d.columns[i] = columnName
		d.columnsSize[i] = len(columnName)
	}
}

func (d *tableData) parseLine(elem reflect.Value) ([]string, error) {
	res := make([]string, len(d.columns))
	if len(d.columns) != elem.NumField() {
		return nil, fmt.Errorf("unexpected number of field %d, expected %d", elem.NumField(), len(d.columns))
	}
	for i := 0; i < len(d.columns); i++ {
		value := fmt.Sprintf("%v", elem.Field(i))
		res[i] = value
		if len(value) > d.columnsSize[i] {
			d.columnsSize[i] = len(value)
		}
	}
	return res, nil
}

func (d *tableData) parseLines(slice reflect.Value) error {
	d.lines = make([][]string, slice.Len())
	var err error
	for i := 0; i < slice.Len(); i++ {
		d.lines[i], err = d.parseLine(slice.Index(i))
		if err != nil {
			return fmt.Errorf("tablify: line %d: %s", i, err)
		}
	}
	return nil
}

func reflectSlice(slice interface{}) (*tableData, error) {
	elemType, err := checkType(slice)
	if err != nil {
		return nil, err
	}

	res := &tableData{}
	res.parseColumns(elemType)

	err = res.parseLines(reflect.ValueOf(slice))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d tableData) lineFormat() string {
	format := "│"
	prefix := ""

	for _, size := range d.columnsSize {
		format += " %" + prefix + strconv.Itoa(size) + "s │"
		prefix = "-"
	}
	return format + "\n"
}

func (d tableData) sepLine(st, md, end string) string {
	res := ""
	for i, size := range d.columnsSize {
		if i == 0 {
			res += st
		} else {
			res += md
		}
		res += strings.Repeat("─", size+2)
	}
	return res + end
}

func fprintSlice(w io.Writer, format string, args ...string) {
	argsi := make([]interface{}, 0, len(args))
	for _, a := range args {
		argsi = append(argsi, a)
	}
	fmt.Fprintf(w, format, argsi...)
}

func (d tableData) fprintf(w io.Writer) {
	fmt.Fprintln(w, d.sepLine("┌", "┬", "┐"))
	lineFormat := d.lineFormat()
	fprintSlice(w, lineFormat, d.columns...)
	fmt.Fprintln(w, d.sepLine("├", "┼", "┤"))
	for _, l := range d.lines {
		fprintSlice(w, lineFormat, l...)
	}
	fmt.Fprintln(w, d.sepLine("└", "┴", "┘"))

}
