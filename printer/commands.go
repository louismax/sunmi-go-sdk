package printer

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
