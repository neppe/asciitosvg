// Copyright 2012 - 2015 The ASCIIToSVG Contributors
// All rights reserved.

package asciitosvg

import (
	"image"
	"strings"
	"testing"

	"github.com/maruel/ut"
)

func TestNewCanvas(t *testing.T) {
	t.Parallel()
	data := []struct {
		input   []string
		strings []string
		texts   []string
		corners [][]image.Point
	}{
		// 0 Small box
		{
			[]string{
				"+-+",
				"| |",
				"+-+",
			},
			[]string{"Path{[(0,0) (1,0) (2,0) (2,1) (2,2) (1,2) (0,2) (0,1)]}"},
			[]string{""},
			[][]image.Point{{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}}},
		},

		// 1 Tight box
		{
			[]string{
				"++",
				"++",
			},
			[]string{"Path{[(0,0) (1,0) (1,1) (0,1)]}"},
			[]string{""},
			[][]image.Point{
				{
					{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1},
				},
			},
		},

		// 2 Indented box
		{
			[]string{
				"",
				" +-+",
				" | |",
				" +-+",
			},
			[]string{"Path{[(1,1) (2,1) (3,1) (3,2) (3,3) (2,3) (1,3) (1,2)]}"},
			[]string{""},
			[][]image.Point{{{X: 1, Y: 1}, {X: 3, Y: 1}, {X: 3, Y: 3}, {X: 1, Y: 3}}},
		},

		// 3 Free flow text
		{
			[]string{
				"",
				" foo bar ",
				"b  baz   bee",
			},
			[]string{"Text{(1,1) \"foo bar\"}", "Text{(0,2) \"b  baz\"}", "Text{(9,2) \"bee\"}"},
			[]string{"foo bar", "b  baz", "bee"},
			[][]image.Point{
				{{X: 1, Y: 1}, {X: 7, Y: 1}},
				{{X: 0, Y: 2}, {X: 5, Y: 2}},
				{{X: 9, Y: 2}, {X: 11, Y: 2}},
			},
		},

		// 4 Text in a box
		{
			[]string{
				"+--+",
				"|Hi|",
				"+--+",
			},
			[]string{"Path{[(0,0) (1,0) (2,0) (3,0) (3,1) (3,2) (2,2) (1,2) (0,2) (0,1)]}", "Text{(1,1) \"Hi\"}"},
			[]string{"", "Hi"},
			[][]image.Point{
				{{X: 0, Y: 0}, {X: 3, Y: 0}, {X: 3, Y: 2}, {X: 0, Y: 2}},
				{{X: 1, Y: 1}, {X: 2, Y: 1}},
			},
		},

		// 5 Concave pieces
		{
			[]string{
				"    +----+",
				"    |    |",
				"+---+    +----+",
				"|             |",
				"+-------------+",
				"", // 5
				"+----+",
				"|    |",
				"|    +---+",
				"|        |",
				"|    +---+", // 10
				"|    |",
				"+----+",
				"",
				"    +----+",
				"    |    |", // 15
				"+---+    |",
				"|        |",
				"+---+    |",
				"    |    |",
				"    +----+",
			},
			[]string{
				"Path{[(4,0) (5,0) (6,0) (7,0) (8,0) (9,0) (9,1) (9,2) (10,2) (11,2) (12,2) (13,2) (14,2) (14,3) (14,4) (13,4) (12,4) (11,4) (10,4) (9,4) (8,4) (7,4) (6,4) (5,4) (4,4) (3,4) (2,4) (1,4) (0,4) (0,3) (0,2) (1,2) (2,2) (3,2) (4,2) (4,1)]}",
				"Path{[(0,6) (1,6) (2,6) (3,6) (4,6) (5,6) (5,7) (5,8) (6,8) (7,8) (8,8) (9,8) (9,9) (9,10) (8,10) (7,10) (6,10) (5,10) (5,11) (5,12) (4,12) (3,12) (2,12) (1,12) (0,12) (0,11) (0,10) (0,9) (0,8) (0,7)]}",
				"Path{[(4,14) (5,14) (6,14) (7,14) (8,14) (9,14) (9,15) (9,16) (9,17) (9,18) (9,19) (9,20) (8,20) (7,20) (6,20) (5,20) (4,20) (4,19) (4,18) (3,18) (2,18) (1,18) (0,18) (0,17) (0,16) (1,16) (2,16) (3,16) (4,16) (4,15)]}",
			},
			[]string{"", "", ""},
			[][]image.Point{
				{
					{X: 4, Y: 0}, {X: 9, Y: 0}, {X: 9, Y: 2}, {X: 14, Y: 2},
					{X: 14, Y: 4}, {X: 0, Y: 4}, {X: 0, Y: 2}, {X: 4, Y: 2},
				},
				{
					{X: 0, Y: 6}, {X: 5, Y: 6}, {X: 5, Y: 8}, {X: 9, Y: 8},
					{X: 9, Y: 10}, {X: 5, Y: 10}, {X: 5, Y: 12}, {X: 0, Y: 12},
				},
				{
					{X: 4, Y: 14}, {X: 9, Y: 14}, {X: 9, Y: 20}, {X: 4, Y: 20},
					{X: 4, Y: 18}, {X: 0, Y: 18}, {X: 0, Y: 16}, {X: 4, Y: 16},
				},
			},
		},

		// 6 Inner boxes
		{
			[]string{
				"+-----+",
				"|     |",
				"| +-+ |",
				"| | | |",
				"| +-+ |",
				"|     |",
				"+-----+",
			},
			[]string{
				"Path{[(0,0) (1,0) (2,0) (3,0) (4,0) (5,0) (6,0) (6,1) (6,2) (6,3) (6,4) (6,5) (6,6) (5,6) (4,6) (3,6) (2,6) (1,6) (0,6) (0,5) (0,4) (0,3) (0,2) (0,1)]}",
				"Path{[(2,2) (3,2) (4,2) (4,3) (4,4) (3,4) (2,4) (2,3)]}",
			},
			[]string{"", ""},
			[][]image.Point{
				{{X: 0, Y: 0}, {X: 6, Y: 0}, {X: 6, Y: 6}, {X: 0, Y: 6}},
				{{X: 2, Y: 2}, {X: 4, Y: 2}, {X: 4, Y: 4}, {X: 2, Y: 4}},
			},
		},

		// 7 Real world diagram example
		{
			[]string{
				//         1         2         3
				"      +------+",
				"      |Editor|-------------+--------+",
				"      +------+             |        |",
				"          |                |        v",
				"          v                |   +--------+",
				"      +------+             |   |Document|", // 5
				"      |Window|             |   +--------+",
				"      +------+             |",
				"         |                 |",
				"   +-----+-------+         |",
				"   |             |         |", // 10
				"   v             v         |",
				"+------+     +------+      |",
				"|Window|     |Window|      |",
				"+------+     +------+      |",
				"                |          |", // 15
				"                v          |",
				"              +----+       |",
				"              |View|       |",
				"              +----+       |",
				"                |          |", // 20
				"                v          |",
				"            +--------+     |",
				"            |Document|<----+",
				"            +--------+",
			},
			[]string{
				"Path{[(6,0) (7,0) (8,0) (9,0) (10,0) (11,0) (12,0) (13,0) (13,1) (13,2) (12,2) (11,2) (10,2) (9,2) (8,2) (7,2) (6,2) (6,1)]}",
				"Path{[(14,1) (15,1) (16,1) (17,1) (18,1) (19,1) (20,1) (21,1) (22,1) (23,1) (24,1) (25,1) (26,1) (27,1) (28,1) (29,1) (30,1) (31,1) (32,1) (33,1) (34,1) (35,1) (36,1) (36,2) (36,3)]}",
				"Path{[(14,1) (15,1) (16,1) (17,1) (18,1) (19,1) (20,1) (21,1) (22,1) (23,1) (24,1) (25,1) (26,1) (27,1) (27,2) (27,3) (27,4) (27,5) (27,6) (27,7) (27,8) (27,9) (27,10) (27,11) (27,12) (27,13) (27,14) (27,15) (27,16) (27,17) (27,18) (27,19) (27,20) (27,21) (27,22) (27,23) (26,23) (25,23) (24,23) (23,23) (22,23)]}",
				"Path{[(10,3) (10,4)]}",
				"Path{[(31,4) (32,4) (33,4) (34,4) (35,4) (36,4) (37,4) (38,4) (39,4) (40,4) (40,5) (40,6) (39,6) (38,6) (37,6) (36,6) (35,6) (34,6) (33,6) (32,6) (31,6) (31,5)]}",
				"Path{[(6,5) (7,5) (8,5) (9,5) (10,5) (11,5) (12,5) (13,5) (13,6) (13,7) (12,7) (11,7) (10,7) (9,7) (8,7) (7,7) (6,7) (6,6)]}",
				"Path{[(9,8) (9,9)]}",
				"Path{[(9,9) (8,9) (7,9) (6,9) (5,9) (4,9) (3,9) (3,10) (3,11)]}",
				"Path{[(9,9) (10,9) (11,9) (12,9) (13,9) (14,9) (15,9) (16,9) (17,9) (17,10) (17,11)]}",
				"Path{[(0,12) (1,12) (2,12) (3,12) (4,12) (5,12) (6,12) (7,12) (7,13) (7,14) (6,14) (5,14) (4,14) (3,14) (2,14) (1,14) (0,14) (0,13)]}",
				"Path{[(13,12) (14,12) (15,12) (16,12) (17,12) (18,12) (19,12) (20,12) (20,13) (20,14) (19,14) (18,14) (17,14) (16,14) (15,14) (14,14) (13,14) (13,13)]}",
				"Path{[(16,15) (16,16)]}",
				"Path{[(14,17) (15,17) (16,17) (17,17) (18,17) (19,17) (19,18) (19,19) (18,19) (17,19) (16,19) (15,19) (14,19) (14,18)]}",
				"Path{[(16,20) (16,21)]}",
				"Path{[(12,22) (13,22) (14,22) (15,22) (16,22) (17,22) (18,22) (19,22) (20,22) (21,22) (21,23) (21,24) (20,24) (19,24) (18,24) (17,24) (16,24) (15,24) (14,24) (13,24) (12,24) (12,23)]}",
				"Text{(7,1) \"Editor\"}",
				"Text{(32,5) \"Document\"}",
				"Text{(7,6) \"Window\"}",
				"Text{(1,13) \"Window\"}",
				"Text{(14,13) \"Window\"}",
				"Text{(15,18) \"View\"}",
				"Text{(13,23) \"Document\"}",
			},
			[]string{
				"", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
				"Editor", "Document", "Window", "Window", "Window", "View", "Document",
			},
			[][]image.Point{
				{{X: 6, Y: 0}, {X: 13, Y: 0}, {X: 13, Y: 2}, {X: 6, Y: 2}},
				{{X: 14, Y: 1}, {X: 36, Y: 1}, {X: 36, Y: 3}},
				{{X: 14, Y: 1}, {X: 27, Y: 1}, {X: 27, Y: 23}, {X: 22, Y: 23}},
				{{X: 10, Y: 3}, {X: 10, Y: 4}},
				{{X: 31, Y: 4}, {X: 40, Y: 4}, {X: 40, Y: 6}, {X: 31, Y: 6}},
				{{X: 6, Y: 5}, {X: 13, Y: 5}, {X: 13, Y: 7}, {X: 6, Y: 7}},
				{{X: 9, Y: 8}, {X: 9, Y: 9}},
				{{X: 9, Y: 9}, {X: 3, Y: 9}, {X: 3, Y: 11}},
				{{X: 9, Y: 9}, {X: 17, Y: 9}, {X: 17, Y: 11}},
				{{X: 0, Y: 12}, {X: 7, Y: 12}, {X: 7, Y: 14}, {X: 0, Y: 14}},
				{{X: 13, Y: 12}, {X: 20, Y: 12}, {X: 20, Y: 14}, {X: 13, Y: 14}},
				{{X: 16, Y: 15}, {X: 16, Y: 16}},
				{{X: 14, Y: 17}, {X: 19, Y: 17}, {X: 19, Y: 19}, {X: 14, Y: 19}},
				{{X: 16, Y: 20}, {X: 16, Y: 21}},
				{{X: 12, Y: 22}, {X: 21, Y: 22}, {X: 21, Y: 24}, {X: 12, Y: 24}},
				{{X: 7, Y: 1}, {X: 12, Y: 1}},
				{{X: 32, Y: 5}, {X: 39, Y: 5}},
				{{X: 7, Y: 6}, {X: 12, Y: 6}},
				{{X: 1, Y: 13}, {X: 6, Y: 13}},
				{{X: 14, Y: 13}, {X: 19, Y: 13}},
				{{X: 15, Y: 18}, {X: 18, Y: 18}},
				{{X: 13, Y: 23}, {X: 20, Y: 23}},
			},
		},

		// 8 Interwined lines.
		{
			[]string{
				"             +-----+-------+",
				"             |     |       |",
				"             |     |       |",
				"        +----+-----+----   |",
				"--------+----+-----+-------+---+",
				"        |    |     |       |   |",
				"        |    |     |       |   |     |   |",
				"        |    |     |       |   |     |   |",
				"        |    |     |       |   |     |   |",
				"--------+----+-----+-------+---+-----+---+--+",
				"        |    |     |       |   |     |   |  |",
				"        |    |     |       |   |     |   |  |",
				"        |   -+-----+-------+---+-----+   |  |",
				"        |    |     |       |   |     |   |  |",
				"        |    |     |       |   +-----+---+--+",
				"             |     |       |         |   |",
				"             |     |       |         |   |",
				"     --------+-----+-------+---------+---+-----",
				"             |     |       |         |   |",
				"             +-----+-------+---------+---+",
			},
			// TODO(dhobsd): it's a tad overwhelming.
			nil,
			nil,
			nil,
		},
		// 9 Indented box
		{
			[]string{
				"",
				"\t+-+",
				"\t| |",
				"\t+-+",
			},
			[]string{"Path{[(9,1) (10,1) (11,1) (11,2) (11,3) (10,3) (9,3) (9,2)]}"},
			[]string{""},
			[][]image.Point{{{X: 9, Y: 1}, {X: 11, Y: 1}, {X: 11, Y: 3}, {X: 9, Y: 3}}},
		},
	}
	for i, line := range data {
		c, err := NewCanvas([]byte(strings.Join(line.input, "\n")), 9)
		if err != nil {
			t.Fatalf("Error creating canvas: %s", err)
		}
		objs := c.Objects()
		if line.strings != nil {
			ut.AssertEqualIndex(t, i, line.strings, getStrings(objs))
		}
		if line.texts != nil {
			ut.AssertEqualIndex(t, i, line.texts, getTexts(objs))
		}
		if line.corners != nil {
			ut.AssertEqualIndex(t, i, line.corners, getCorners(objs))
		}
	}
}

func TestNewCanvasBroken(t *testing.T) {
	// These are the ones that do not give the desired result.
	t.Parallel()
	data := []struct {
		input   []string
		strings []string
		texts   []string
		corners [][]image.Point
	}{
		// 0 URL
		{
			[]string{
				"github.com/foo/bar",
			},
			[]string{"Text{(0,0) \"github.com/foo/bar\"}"},
			[]string{"github.com/foo/bar"},
			[][]image.Point{{{X: 0, Y: 0}, {X: 17, Y: 0}}},
		},

		// 1 Merged boxes
		{
			[]string{
				"+-+-+",
				"| | |",
				"+-+-+",
			},
			[]string{"Path{[(0,0) (1,0) (2,0) (3,0) (4,0) (4,1) (4,2) (3,2) (2,2) (1,2) (0,2) (0,1)]}", "Path{[(0,0) (1,0) (2,0) (3,0) (4,0) (4,1) (4,2) (3,2) (2,2) (2,1)]}"},
			[]string{"", ""},
			// TODO(dhobsd): BROKEN.
			[][]image.Point{
				{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 2}, {X: 0, Y: 2}},
				{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}},
			},
		},

		// 2 Adjacent boxes
		{
			// TODO(dhobsd): BROKEN. This one is hard, as it can be seen as 3 boxes
			// but that is not what is desired.
			[]string{
				"+-++-+",
				"| || |",
				"+-++-+",
			},
			[]string{
				"Path{[(0,0) (1,0) (2,0) (3,0) (4,0) (5,0) (5,1) (5,2) (4,2) (3,2) (2,2) (1,2) (0,2) (0,1)]}",
				"Path{[(0,0) (1,0) (2,0) (3,0) (4,0) (5,0) (5,1) (5,2) (4,2) (3,2) (2,2) (2,1)]}",
				"Path{[(0,0) (1,0) (2,0) (3,0) (4,0) (5,0) (5,1) (5,2) (4,2) (3,2) (3,1)]}",
			},
			[]string{"", "", ""},
			[][]image.Point{
				{{X: 0, Y: 0}, {X: 5, Y: 0}, {X: 5, Y: 2}, {X: 0, Y: 2}},
				{{X: 0, Y: 0}, {X: 5, Y: 0}, {X: 5, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}},
				{{X: 0, Y: 0}, {X: 5, Y: 0}, {X: 5, Y: 2}, {X: 3, Y: 2}, {X: 3, Y: 1}},
			},
		},
	}
	for i, line := range data {
		c, err := NewCanvas([]byte(strings.Join(line.input, "\n")), 9)
		if err != nil {
			t.Fatalf("Error creating canvas: %s", err)
		}
		objs := c.Objects()
		if line.strings != nil {
			ut.AssertEqualIndex(t, i, line.strings, getStrings(objs))
		}
		if line.texts != nil {
			ut.AssertEqualIndex(t, i, line.texts, getTexts(objs))
		}
		if line.corners != nil {
			ut.AssertEqualIndex(t, i, line.corners, getCorners(objs))
		}
	}
}

func TestPointsToCorners(t *testing.T) {
	t.Parallel()
	data := []struct {
		in       []image.Point
		expected []image.Point
		closed   bool
	}{
		{
			[]image.Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
			[]image.Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
			false,
		},
		{
			[]image.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			[]image.Point{{X: 0, Y: 0}, {X: 2, Y: 0}},
			false,
		},
		{
			[]image.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
			[]image.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
			false,
		},
		{
			[]image.Point{
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2},
				{X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1},
			},
			[]image.Point{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}},
			true,
		},
		{
			[]image.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}},
			[]image.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}},
			// TODO(dhobsd): Unexpected; broken.
			false,
		},
	}
	for i, line := range data {
		p, c := pointsToCorners(line.in)
		ut.AssertEqualIndex(t, i, line.expected, p)
		ut.AssertEqualIndex(t, i, line.closed, c)
	}
}

func BenchmarkT(b *testing.B) {
	data := []string{
		"             +-----+-------+",
		"             |     |       |",
		"             |     |       |",
		"        +----+-----+----   |",
		"--------+----+-----+-------+---+",
		"        |    |     |       |   |",
		"        |    |     |       |   |     |   |",
		"        |    |     |       |   |     |   |",
		"        |    |     |       |   |     |   |",
		"--------+----+-----+-------+---+-----+---+--+",
		"        |    |     |       |   |     |   |  |",
		"        |    |     |       |   |     |   |  |",
		"        |   -+-----+-------+---+-----+   |  |",
		"        |    |     |       |   |     |   |  |",
		"        |    |     |       |   +-----+---+--+",
		"             |     |       |         |   |",
		"             |     |       |         |   |",
		"     --------+-----+-------+---------+---+-----",
		"             |     |       |         |   |",
		"             +-----+-------+---------+---+",
		"",
		"",
	}
	chunk := []byte(strings.Join(data, "\n"))
	input := make([]byte, 0, len(chunk)*b.N)
	for i := 0; i < b.N; i++ {
		input = append(input, chunk...)
	}
	expected := 30 * b.N
	b.ResetTimer()
	c, err := NewCanvas(input, 8)
	if err != nil {
		b.Fatalf("Error creating canvas: %s", err)
	}

	objs := c.Objects()
	if len(objs) != expected {
		b.Fatalf("%d != %d", len(objs), expected)
	}
}

// Private details.

func getPoints(objs []Object) [][]image.Point {
	out := [][]image.Point{}
	for _, obj := range objs {
		out = append(out, obj.Points())
	}
	return out
}

func getTexts(objs []Object) []string {
	out := []string{}
	for _, obj := range objs {
		t := obj.Text()
		if !obj.IsText() {
			out = append(out, "")
		} else if len(t) > 0 {
			out = append(out, string(t))
		} else {
			panic("failed")
		}
	}
	return out
}

func getStrings(objs []Object) []string {
	out := []string{}
	for _, obj := range objs {
		out = append(out, obj.String())
	}
	return out
}

func getCorners(objs []Object) [][]image.Point {
	out := make([][]image.Point, len(objs))
	for i, obj := range objs {
		out[i] = obj.Corners()
	}
	return out
}