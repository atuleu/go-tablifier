package tablifier

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

type toPrint struct {
	Node   string
	Since  time.Duration
	AnInt  int `name:"Its nice"`
	AFloat float64
}

type TablifySuite struct {
	data []toPrint
}

var _ = Suite(&TablifySuite{
	data: []toPrint{
		{
			Node:   "foo",
			Since:  time.Hour + time.Minute,
			AnInt:  12,
			AFloat: 24.23,
		},
		{
			Node:   "foobar",
			Since:  45*time.Hour - 3*time.Minute,
			AnInt:  1267678,
			AFloat: -24.23,
		},
	},
})

func Test(t *testing.T) { TestingT(t) }

func (s *TablifySuite) TestParse(c *C) {
	tdata, err := reflectSlice(s.data)
	c.Assert(err, IsNil)
	c.Check(tdata.columns, DeepEquals, []string{"  Node", "Since   ", "Its nice", "AFloat"})
	c.Check(tdata.lines, DeepEquals, [][]string{
		{
			"   foo", "1h1m0s  ", "12      ", "24.23 ",
		},
		{
			"foobar", "44h57m0s", "1267678 ", "-24.23",
		},
	})
	c.Check(tdata.columnsSize, DeepEquals, []int{6, 8, 8, 6})

}

func (s *TablifySuite) TestComputeLength(c *C) {
	testdata := []struct {
		Value    string
		Expected int
	}{
		{"foobar", 6},
		{"\033[1;96mfoo\033[m", 3},
		{"\033[1;96m∞\033[m", 1},
	}
	for _, d := range testdata {
		comment := Commentf("Value: %s, Expected: %d", d.Value, d.Expected)
		c.Check(computeLength(d.Value), Equals, d.Expected, comment)
	}
}

func ExampleTablify() {
	type ToPrint struct {
		Node   string
		Since  time.Duration
		AnInt  int     `name:"An Int"`
		AFloat float64 `name:"A Float"`
	}

	data := []ToPrint{
		{
			Node:   "foo",
			Since:  time.Hour + time.Minute,
			AnInt:  12,
			AFloat: 24.23,
		},
		{
			Node:   "\033[31mfoobar\033[m",
			Since:  45*time.Hour - 3*time.Minute,
			AnInt:  1267678,
			AFloat: -24.23,
		},
		{
			Node:   "\033[1;96m∞\033[m",
			Since:  0,
			AnInt:  0,
			AFloat: 0.0,
		},
	}

	Tablify(data)
	//output:
	//┌────────┬──────────┬─────────┬─────────┐
	//│   Node │ Since    │ An Int  │ A Float │
	//├────────┼──────────┼─────────┼─────────┤
	//│    foo │ 1h1m0s   │ 12      │ 24.23   │
	//│ [31mfoobar[m │ 44h57m0s │ 1267678 │ -24.23  │
	//│      [1;96m∞[m │ 0s       │ 0       │ 0       │
	//└────────┴──────────┴─────────┴─────────┘
}
