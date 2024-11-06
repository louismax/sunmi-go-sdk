package printer

import "unicode/utf8"

// Append raw data.
func (p *PrintObject) appendRawData(data string) {
	p.Content += data
}

// Append unicode character.
func (p *PrintObject) appendUnicode(unicode, count int) {
	for i := 0; i < count; i++ {
		p.Content += p.unicodeToUtf8(unicode)
	}
}

// AppendText Append text.
func (p *PrintObject) AppendText(str string) {
	for _, r := range str {
		p.Content += p.unicodeToUtf8(int(r))
	}
}

// LineFeed [LF]打印缓冲区和进纸行中的数据
func (p *PrintObject) LineFeed(n ...int) {
	if len(n) > 0 {
		for i := 0; i < n[0]; i++ {
			p.Content += "0a"
		}
	} else {
		p.Content += "0a"
	}
}

// RestoreDefaultSettings [ESC @] 恢复默认设置
func (p *PrintObject) RestoreDefaultSettings() {
	p.CharHSize = 1
	p.Content += "1b40"
}

// RestoreDefaultLineSpacing [ESC 2] 恢复默认行距
func (p *PrintObject) RestoreDefaultLineSpacing() {
	p.Content += "1b32"
}

// SetLineSpacing [ESC 3] 设置行距
func (p *PrintObject) SetLineSpacing(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1b33" + p.numToHexStr(n, 1)
	}
}

// SetPrintModes [ESC !] 设置打印模式
func (p *PrintObject) SetPrintModes(bold, doubleH, doubleW bool) {
	var n uint8
	if bold {
		n |= 8
	}
	if doubleH {
		n |= 16
	}
	if doubleW {
		n |= 32
	}
	if doubleW {
		p.CharHSize = 2
	} else {
		p.CharHSize = 1
	}
	p.Content += "1b21" + p.numToHexStr(int(n), 1)
}

// SetCharacterSize [GS !] 设置字符大小
func (p *PrintObject) SetCharacterSize(h, w int) {
	var n uint8 = 0

	if h >= 1 && h <= 8 {
		n |= uint8(h - 1)
	}
	if w >= 1 && w <= 8 {
		n |= uint8(w-1) << 4
		p.CharHSize = w
	}
	p.Content += "1d21" + p.numToHexStr(int(n), 1)
}

// HorizontalTab [HT] 插入水平制表符
func (p *PrintObject) HorizontalTab(n ...int) {
	if len(n) > 0 {
		for i := 0; i < n[0]; i++ {
			p.Content += "09"
		}
	} else {
		p.Content += "09"
	}
}

// SetAbsolutePrintPosition [ESC $] 设置绝对打印位置.
func (p *PrintObject) SetAbsolutePrintPosition(n int) {
	if n >= 0 && n <= 65535 {
		p.Content += "1b24" + p.numToHexStr(n, 2)
	}
}

// SetRelativePrintPosition [ESC \] 设置相对打印位置.
func (p *PrintObject) SetRelativePrintPosition(n int) {
	if n >= -32768 && n <= 32767 {
		p.Content += "1b5c" + p.numToHexStr(n, 2)
	}
}

// SetAlignment [ESC a] 设置对齐方式.
func (p *PrintObject) SetAlignment(n int) {
	if n >= 0 && n <= 2 {
		p.Content += "1b61" + p.numToHexStr(n, 1)
	}
}

// SetUnderlineMode [ESC -] 设置下划线模式.
func (p *PrintObject) SetUnderlineMode(n int) {
	if n >= 0 && n <= 2 {
		p.Content += "1b2d" + p.numToHexStr(n, 1)
	}
}

// SetBlackWhiteReverseMode [GS B] 设置黑白倒转模式.
func (p *PrintObject) SetBlackWhiteReverseMode(ok bool) {
	if ok {
		p.Content += "1d4201"
	} else {
		p.Content += "1d4200"
	}
}

// SetUpsideDownMode [ESC {] 设置倒立模式.
func (p *PrintObject) SetUpsideDownMode(ok bool) {
	if ok {
		p.Content += "1b7b01"
	} else {
		p.Content += "1b7b00"
	}
}

// CutPaper [GS V m] 切纸.
func (p *PrintObject) CutPaper(ok bool) {
	if ok {
		p.Content += "1d5630"
	} else {
		p.Content += "1d5631"
	}
}

// PostponedCutPaper [GS V m n] 延期裁纸. 在收到此命令后，打印机将不执行切割，直到(d + n)点线馈送，其中d是打印位置和切割位置之间的距离
func (p *PrintObject) PostponedCutPaper(ok bool, n int) {
	if n >= 0 && n <= 255 {
		if ok {
			p.Content += "1d5661" + p.numToHexStr(n, 1)
		} else {
			p.Content += "1d5662" + p.numToHexStr(n, 1)
		}
	}
}

//////////////////////////////////////////////////
// Sunmi专用命令
//////////////////////////////////////////////////

// SetCjkEncoding (effective when UTF-8 mode is disabled).
//
//	n  encoding
//
// ---  --------
//
//	 0  GB18030
//	 1  BIG5
//	11  Shift_JIS
//	12  JIS 0208.
//	21  KS C 5601.
//
// 128  Disable CJK mode
// 255  Restore to default
func (p *PrintObject) SetCjkEncoding(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1d284503000601" + p.numToHexStr(n, 1)
	}
}

// SetUtf8Mode .
//
//	n  mode
//
// ---  ----
//
//	0  Disabled
//	1  Enabled
//
// 255  Restore to default
func (p *PrintObject) SetUtf8Mode(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1d284503000603" + p.numToHexStr(n, 1)
	}
}

// SetHalfBuzzAsciiCharSize Set Latin character size of vector font.
func (p *PrintObject) SetHalfBuzzAsciiCharSize(n int) {
	if n >= 0 && n <= 255 {
		p.AsciiCharWidth = n
		p.Content += "1d28450300060a" + p.numToHexStr(n, 1)
	}
}

// SetHalfBuzzCjkCharSize Set CJK character size of vector font.
func (p *PrintObject) SetHalfBuzzCjkCharSize(n int) {
	if n >= 0 && n <= 255 {
		p.CjkCharWidth = n
		p.Content += "1d28450300060b" + p.numToHexStr(n, 1)
	}
}

// SetHalfBuzzOtherCharSize Set other character size of vector font.
func (p *PrintObject) SetHalfBuzzOtherCharSize(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1d28450300060c" + p.numToHexStr(n, 1)
	}
}

// SelectAsciiCharFont Select font for Latin characters.
func (p *PrintObject) SelectAsciiCharFont(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1d284503000614" + p.numToHexStr(n, 1)
	}
}

// SelectCjkCharFont Select font for CJK characters.
//
//	n  font
//
// -----  ----
//
//	0  Built-in lattice font
//	1  Built-in vector font
//
// >=128  The (n-128)th custom vector font
func (p *PrintObject) SelectCjkCharFont(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1d284503000615" + p.numToHexStr(n, 1)
	}
}

// SelectOtherCharFont Select font for other characters.
//
//	n  font
//
// -----  ----
//
//	0,1  Built-in vector font
//
// >=128  The (n-128)th custom vector font
func (p *PrintObject) SelectOtherCharFont(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1d284503000616" + p.numToHexStr(n, 1)
	}
}

// SetPrintDensity Set print density.
func (p *PrintObject) SetPrintDensity(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1d2845020007" + p.numToHexStr(n, 1)
	}
}

// SetPrintSpeed Set print speed.
func (p *PrintObject) SetPrintSpeed(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1d2845020008" + p.numToHexStr(n, 1)
	}
}

// SetCutterMode  Set cutter mode.
// n  mode
// -  ----
// 0  Perform full-cut or partial-cut according to the cutting command
// 1  Perform partial-cut always on any cutting command
// 2  Perform full-cut always on any cutting command
// 3  Never cut on any cutting command
func (p *PrintObject) SetCutterMode(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1d2845020010" + p.numToHexStr(n, 1)
	}
}

// ClearPaperNotTakenAlarm Clear paper-not-taken alarm.
func (p *PrintObject) ClearPaperNotTakenAlarm(n int) {
	if n >= 0 && n <= 255 {
		p.Content += "1d2854010004" + p.numToHexStr(n, 1)
	}
}

// WidthOfChar 字符宽度
func (p *PrintObject) WidthOfChar(charCode rune) int {
	if charCode >= 0x00020 && charCode <= 0x0036f {
		return p.AsciiCharWidth
	}
	if charCode >= 0x0ff61 && charCode <= 0x0ff9f {
		return p.CjkCharWidth / 2
	}
	if charCode == 0x02010 ||
		(charCode >= 0x02013 && charCode <= 0x02016) ||
		(charCode >= 0x02018 && charCode <= 0x02019) ||
		(charCode >= 0x0201c && charCode <= 0x0201d) ||
		(charCode >= 0x02025 && charCode <= 0x02026) ||
		(charCode >= 0x02030 && charCode <= 0x02033) ||
		charCode == 0x02035 ||
		charCode == 0x0203b {
		return p.CjkCharWidth
	}
	if (charCode >= 0x01100 && charCode <= 0x011ff) ||
		(charCode >= 0x02460 && charCode <= 0x024ff) ||
		(charCode >= 0x025a0 && charCode <= 0x027bf) ||
		(charCode >= 0x02e80 && charCode <= 0x02fdf) ||
		(charCode >= 0x03000 && charCode <= 0x0318f) ||
		(charCode >= 0x031a0 && charCode <= 0x031ef) ||
		(charCode >= 0x03200 && charCode <= 0x09fff) ||
		(charCode >= 0x0ac00 && charCode <= 0x0d7ff) ||
		(charCode >= 0x0f900 && charCode <= 0x0faff) ||
		(charCode >= 0x0fe30 && charCode <= 0x0fe4f) ||
		(charCode >= 0x1f000 && charCode <= 0x1f9ff) {
		return p.CjkCharWidth
	}
	if (charCode >= 0x0ff01 && charCode <= 0x0ff5e) ||
		(charCode >= 0x0ffe0 && charCode <= 0x0ffe5) {
		return p.CjkCharWidth
	}
	return p.AsciiCharWidth
}

// SetupColumns 设置列
func (p *PrintObject) SetupColumns(columns [][]int) {
	p.ColumnSettings = []ColumnSettings{}
	remain := p.DotsPerLine
	for _, col := range columns {
		width := col[0]
		alignment := col[1]
		flag := col[2]
		if width == 0 || width > remain {
			width = remain
		}
		p.ColumnSettings = append(p.ColumnSettings, ColumnSettings{Width: width, Alignment: alignment, Flag: flag})
		remain -= width
		if remain == 0 {
			return
		}
	}
}

// PrintInColumns 按列打印
func (p *PrintObject) PrintInColumns(texts []string) {
	if len(p.ColumnSettings) == 0 || len(texts) == 0 {
		return
	}

	strcur := make([]string, 0)
	strrem := make([]string, 0)
	strwidth := make([]int, 0)
	numOfColumns := min(len(p.ColumnSettings), len(texts))

	for i := 0; i < numOfColumns; i++ {
		strcur = append(strcur, "")
		strrem = append(strrem, texts[i])
		strwidth = append(strwidth, 0)
	}

	for {
		done := true
		pos := 0

		for i := 0; i < numOfColumns; i++ {
			width := p.ColumnSettings[i].Width
			alignment := p.ColumnSettings[i].Alignment
			flag := p.ColumnSettings[i].Flag

			if len(strrem[i]) == 0 {
				pos += width
				continue
			}

			done = false
			strcur[i] = ""
			strwidth[i] = 0
			j := 0
			for j < len(strrem[i]) {
				r, size := utf8.DecodeRuneInString(strrem[i][j:])
				if r == '\n' {
					j += size
					break
				} else {
					w := p.WidthOfChar(r) * p.CharHSize
					if flag&COLUMN_FLAG_DOUBLE_W != 0 {
						w *= 2
					}
					if strwidth[i]+w > width {
						break
					} else {
						strcur[i] += string(r)
						strwidth[i] += w
					}
				}
				j += size
			}
			if j < len(strrem[i]) {
				strrem[i] = strrem[i][j:]
			} else {
				strrem[i] = ""
			}

			switch alignment {
			case 1:
				p.SetAbsolutePrintPosition(pos + (width-strwidth[i])/2)
			case 2:
				p.SetAbsolutePrintPosition(pos + (width - strwidth[i]))
			default:
				p.SetAbsolutePrintPosition(pos)
			}
			if flag&COLUMN_FLAG_BW_REVERSE != 0 {
				p.SetBlackWhiteReverseMode(true)
			}
			if flag&(COLUMN_FLAG_BOLD|COLUMN_FLAG_DOUBLE_H|COLUMN_FLAG_DOUBLE_W) != 0 {
				bold := flag&COLUMN_FLAG_BOLD != 0
				doubleH := flag&COLUMN_FLAG_DOUBLE_H != 0
				doubleW := flag&COLUMN_FLAG_DOUBLE_W != 0
				p.SetPrintModes(bold, doubleH, doubleW)
			}
			p.AppendText(strcur[i])
			if flag&(COLUMN_FLAG_BOLD|COLUMN_FLAG_DOUBLE_H|COLUMN_FLAG_DOUBLE_W) != 0 {
				p.SetPrintModes(false, false, false)
			}
			if flag&COLUMN_FLAG_BW_REVERSE != 0 {
				p.SetBlackWhiteReverseMode(false)
			}
			pos += width
		}

		if !done {
			p.LineFeed()
		} else {
			break
		}
	}
}

//////////////////////////////////////////////////
// 条码和二维码打印
//////////////////////////////////////////////////

func (p *PrintObject) AppendBarcode(hriPos, height, moduleSize, barcodeType int, text string) {
	textLength := len(text)

	if textLength == 0 {
		return
	}
	if textLength > 255 {
		textLength = 255
	}
	if height < 1 {
		height = 1
	} else if height > 255 {
		height = 255
	}
	if moduleSize < 1 {
		moduleSize = 1
	} else if moduleSize > 6 {
		moduleSize = 6
	}

	p.Content += "1d48" + p.numToHexStr(hriPos&3, 1)
	p.Content += "1d6600"
	p.Content += "1d68" + p.numToHexStr(height, 1)
	p.Content += "1d77" + p.numToHexStr(moduleSize, 1)
	p.Content += "1d6b" + p.numToHexStr(barcodeType, 1) + p.numToHexStr(textLength, 1)
	for _, r := range text {
		//p.Content += p.unicodeToUtf8(int(r))
		p.Content += p.numToHexStr(int(r), 1)
	}
}

func (p *PrintObject) AppendQRCode(moduleSize, ecLevel int, text string) {
	content := ""
	for _, r := range text {
		content += p.unicodeToUtf8(int(r))
	}
	textLength := len(content) / 2

	if textLength == 0 {
		return
	}
	if textLength > 65535 {
		textLength = 65535
	}
	if moduleSize < 1 {
		moduleSize = 1
	} else if moduleSize > 16 {
		moduleSize = 16
	}
	if ecLevel < 0 {
		ecLevel = 0
	} else if ecLevel > 3 {
		ecLevel = 3
	}

	p.Content += "1d286b040031410000"
	p.Content += "1d286b03003143" + p.numToHexStr(moduleSize, 1)
	p.Content += "1d286b03003145" + p.numToHexStr(ecLevel+48, 1)
	p.Content += "1d286b" + p.numToHexStr(textLength+3, 2) + "315030"
	p.Content += content
	p.Content += "1d286b0300315130"

}
