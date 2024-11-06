package printer

import "fmt"

const (
	ALIGN_LEFT   = 0 // Left alignment
	ALIGN_CENTER = 1 // Center alignment
	ALIGN_RIGHT  = 2 // Right alignment

	HRI_POS_ABOVE = 1 // HRI above the barcode
	HRI_POS_BELOW = 2 // HRI below the barcode

	DIFFUSE_DITHER   = 0
	THRESHOLD_DITHER = 2

	COLUMN_FLAG_BW_REVERSE = 1 << 0
	COLUMN_FLAG_BOLD       = 1 << 1
	COLUMN_FLAG_DOUBLE_H   = 1 << 2
	COLUMN_FLAG_DOUBLE_W   = 1 << 3
)

type ColumnSettings struct {
	Width     int
	Alignment int
	Flag      int
}

type PrintObject struct {
	Content        string
	CharHSize      int
	DotsPerLine    int
	AsciiCharWidth int
	CjkCharWidth   int
	ColumnSettings []ColumnSettings
}

func NewPrint(param ...int) *PrintObject {
	_dotsPerLine := 384 // 以网点为单位的打印宽度,58mm为384,80mm为576
	if len(param) > 0 {
		_dotsPerLine = param[0]
	}

	return &PrintObject{
		Content:        "",
		CharHSize:      1,
		DotsPerLine:    _dotsPerLine,
		AsciiCharWidth: 12,
		CjkCharWidth:   24,
		ColumnSettings: make([]ColumnSettings, 0),
	}
}

func (p *PrintObject) numToHexStr(n int, bytes int) string {
	str := ""
	var v int

	for i := 0; i < bytes; i++ {
		v = n & 0xFF
		if v < 0x10 {
			str += fmt.Sprintf("0%x", v)
		} else {
			str += fmt.Sprintf("%x", v)
		}
		n >>= 8
	}
	return str
}

func (p *PrintObject) unicodeToUtf8(unicode int) string {
	var c1, c2, c3, c4 int
	if unicode < 0 {
		return ""
	}
	if unicode <= 0x7F {
		c1 = unicode & 0x7F
		return p.numToHexStr(c1, 1)
	}
	if unicode <= 0x7FF {
		c1 = ((unicode >> 6) & 0x1F) | 0xC0
		c2 = ((unicode) & 0x3F) | 0x80
		return p.numToHexStr(c1, 1) + p.numToHexStr(c2, 1)
	}
	if unicode <= 0xFFFF {
		c1 = ((unicode >> 12) & 0x0F) | 0xE0
		c2 = ((unicode >> 6) & 0x3F) | 0x80
		c3 = ((unicode) & 0x3F) | 0x80
		return p.numToHexStr(c1, 1) + p.numToHexStr(c2, 1) + p.numToHexStr(c3, 1)
	}
	if unicode <= 0x10FFFF {
		c1 = ((unicode >> 18) & 0x07) | 0xF0
		c2 = ((unicode >> 12) & 0x3F) | 0x80
		c3 = ((unicode >> 6) & 0x3F) | 0x80
		c4 = ((unicode) & 0x3F) | 0x80
		return p.numToHexStr(c1, 1) + p.numToHexStr(c2, 1) + p.numToHexStr(c3, 1) + p.numToHexStr(c4, 1)
	}
	return ""
}
