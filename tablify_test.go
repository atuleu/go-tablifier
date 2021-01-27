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
		toPrint{
			Node:   "foo",
			Since:  time.Hour + time.Minute,
			AnInt:  12,
			AFloat: 24.23,
		},
		toPrint{
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
	c.Check(tdata.columns, DeepEquals, []string{"Node", "Since", "Its nice", "AFloat"})
	c.Check(tdata.lines, DeepEquals, [][]string{
		[]string{
			"foo", "1h1m0s", "12", "24.23",
		},
		[]string{
			"foobar", "44h57m0s", "1267678", "-24.23",
		},
	})
	c.Check(tdata.columnsSize, DeepEquals, []int{6, 8, 8, 6})

}

func (s *TablifySuite) TestParseLines(c *C) {

}
