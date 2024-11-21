package printer

import (
	"github.com/mattn/go-runewidth"
	"golang.org/x/image/draw"
	"image"
	"os"
	"regexp"
	"unicode/utf8"
)

var stripAnsiEscapeRegexp = regexp.MustCompile(`(\x9B|\x1B\[)[0-?]*[ -/]*[@-~]`)

func stripAnsiEscape(s string) string {
	return stripAnsiEscapeRegexp.ReplaceAllString(s, "")
}

func realLength(s string) int {
	return runewidth.StringWidth(stripAnsiEscape(s))
}

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

// AppendDivider 添加分割线(设置字符大小1-1,设置居中,LF)
func (p *PrintObject) AppendDivider(data ...string) {
	p.SetCharacterSize(1, 1)
	p.SetAlignment(AlignCenter)
	if p.DotsPerLine == 384 {
		if len(data) < 1 {
			str := ""
			for i := 0; i < 32; i++ {
				str += "-"
			}
			p.AppendText(str)
		} else {
			s := ""
			if realLength(data[0])%2 == 0 {
				s = data[0]
			} else { //如果是奇数,需要补一个空格保证左右分割线一致
				s = " " + data[0]
			}
			if realLength(s) > 32 {
				return
			} else if realLength(s) == 32 {
				p.AppendText(s)
			} else {
				for i := 0; i < (32-realLength(s))/2; i++ {
					p.AppendText("-")
				}
				p.AppendText(s)
				for i := 0; i < (32-realLength(s))/2; i++ {
					p.AppendText("-")
				}
			}
		}
	} else if p.DotsPerLine == 576 {
		if len(data) < 1 {
			str := ""
			for i := 0; i < 48; i++ {
				str += "-"
			}
			p.AppendText(str)
		} else {
			s := ""
			if realLength(data[0])%2 == 0 {
				s = data[0]
			} else { //如果是奇数,需要补一个空格保证左右分割线一致
				s = " " + data[0]
			}
			if realLength(s) > 48 {
				return
			} else if realLength(s) == 48 {
				p.AppendText(s)
			} else {
				for i := 0; i < (48-realLength(s))/2; i++ {
					p.AppendText("-")
				}
				p.AppendText(s)
				for i := 0; i < (48-realLength(s))/2; i++ {
					p.AppendText("-")
				}
			}
		}
	}
	p.LineFeed()
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

// SetBold [ESC E] 设置加粗.
func (p *PrintObject) SetBold(ok bool) {
	if ok {
		p.Content += "1b4501"
	} else {
		p.Content += "1b4500"
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
//	 0  GB18030
//	 1  BIG5
//	11  Shift_JIS
//	12  JIS 0208
//	21  KS C 5601
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
					if flag&ColumnFlagDoubleW != 0 {
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
			if flag&ColumnFlagBwReverse != 0 {
				p.SetBlackWhiteReverseMode(true)
			}
			if flag&(ColumnFlagBold|ColumnFlagDoubleH|ColumnFlagDoubleW) != 0 {
				bold := flag&ColumnFlagBold != 0
				doubleH := flag&ColumnFlagDoubleH != 0
				doubleW := flag&ColumnFlagDoubleW != 0
				p.SetPrintModes(bold, doubleH, doubleW)
			}
			p.AppendText(strcur[i])
			if flag&(ColumnFlagBold|ColumnFlagDoubleH|ColumnFlagDoubleW) != 0 {
				p.SetPrintModes(false, false, false)
			}
			if flag&ColumnFlagBwReverse != 0 {
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

// AppendBarcode 添加条形码 常用barcodeType 73 Code_128
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

// AppendQRCode 添加二维码
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

// ////////////////////////////////////////////////
// 图像打印
// ////////////////////////////////////////////////
// Grayscale to monochrome - diffuse dithering algorithm.
func (p *PrintObject) diffuseDither(srcData [][]int, width int, height int) []int {
	if width <= 0 || height <= 0 {
		return []int{}
	}

	bmwidth := (width + 7) / 8
	dstData := make([]int, bmwidth*height)
	lineBuffer := make([][]int, 2)
	lineBuffer[0] = make([]int, width)
	lineBuffer[1] = make([]int, width)
	line1 := 0
	line2 := 1

	for i := 0; i < width; i++ {
		lineBuffer[0][i] = 0
		lineBuffer[1][i] = srcData[0][i]
	}

	for y := 0; y < height; y++ {
		tmp := line1
		line1 = line2
		line2 = tmp
		notLastLine := y < height-1

		if notLastLine {
			for i := 0; i < width; i++ {
				lineBuffer[line2][i] = srcData[y+1][i]
			}
		}

		q := y * bmwidth
		for i := 0; i < bmwidth; i++ {
			dstData[q] = 0
			q++
		}

		b1 := 0
		b2 := 0
		q = y * bmwidth
		mask := 0x80

		for x := 1; x <= width; x++ {
			var err int
			if lineBuffer[line1][b1] < 128 { // Black pixel
				err = lineBuffer[line1][b1]
				dstData[q] |= mask
			} else {
				err = lineBuffer[line1][b1] - 255
			}
			b1++
			if mask == 1 {
				q++
				mask = 0x80
			} else {
				mask >>= 1
			}
			e7 := (err*7 + 8) >> 4
			e5 := (err*5 + 8) >> 4
			e3 := (err*3 + 8) >> 4
			e1 := err - (e7 + e5 + e3)
			if x < width {
				lineBuffer[line1][b1] += e7
			}
			if notLastLine {
				lineBuffer[line2][b2] += e5
				if x > 1 {
					lineBuffer[line2][b2-1] += e3
				}
				if x < width {
					lineBuffer[line2][b2+1] += e1
				}
			}
			b2++
		}
	}
	return dstData
}

// Grayscale to monochrome - threshold dithering algorithm.
func (p *PrintObject) thresholdDither(srcData [][]int, width int, height int) []int {
	if width <= 0 || height <= 0 {
		return []int{}
	}

	bmwidth := (width + 7) / 8
	dstData := make([]int, bmwidth*height)
	q := 0

	for y := 0; y < height; y++ {
		k := q
		mask := 0x80
		for x := 0; x < width; x++ {
			if srcData[y][x] < 128 { // Black pixel
				dstData[k] |= mask
			}
			if mask == 1 {
				k++
				mask = 0x80
			} else {
				mask >>= 1
			}
		}
		q += bmwidth
	}
	return dstData
}

// Convert image pixel data from RGB to grayscale.
func (p *PrintObject) convertToGray(img image.Image) [][]int {
	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	data := make([][]int, height)
	grayData := make([][]int, height)
	for i := range data {
		data[i] = make([]int, width)
		grayData[i] = make([]int, width)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			grayData[y][x] = int((r*11+g*16+b*5)/32) & 0xff
		}
	}
	return grayData
}

// AppendImage an image.
func (p *PrintObject) AppendImage(imageFile string, mode int, width int) {
	imgOrg, err := openImage(imageFile)
	if err != nil {
		return
	}

	if width == 0 {
		width = p.DotsPerLine
	}
	w, h := imgOrg.Bounds().Max.X, imgOrg.Bounds().Max.Y
	if w > width {
		h = width * h / w
		w = width
		imgRes := resizeImage(imgOrg, w, h)
		grayData := p.convertToGray(imgRes)
		var monoData []int
		if mode == DiffuseDither {
			monoData = p.diffuseDither(grayData, w, h)
		} else if mode == ThresholdDither {
			monoData = p.thresholdDither(grayData, w, h)
		} else {
			return
		}

		w = (w + 7) / 8
		p.Content += "1d763000"
		p.Content += p.numToHexStr(w, 2)
		p.Content += p.numToHexStr(h, 2)
		for _, r := range monoData {
			p.Content += p.numToHexStr(r, 1)
		}
	}
}

func openImage(imageFile string) (image.Image, error) {
	file, err := os.Open(imageFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	img, _, err := image.Decode(file)
	return img, err
}

func resizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	return dst
}

//////////////////////////////////////////////////
// 页面模式命令
//////////////////////////////////////////////////

// EnterPageMode [ESC L] 进入页面模式。
func (p *PrintObject) EnterPageMode() {
	p.Content += "1b4c"
}

// SetPrintAreaInPageMode [ESC W] 在页面模式下设置打印区域。
// x, y: 打印区域的起源
// w, h: 打印区域的宽度和高度
func (p *PrintObject) SetPrintAreaInPageMode(x, y, w, h int) {
	p.Content += "1b57"
	p.Content += p.numToHexStr(x, 2)
	p.Content += p.numToHexStr(y, 2)
	p.Content += p.numToHexStr(w, 2)
	p.Content += p.numToHexStr(h, 2)
}

// SetPrintDirectionInPageMode [ESC T] 在页面模式下选择打印方向。
// dir: 0:正常的; 1:顺时针旋转90度; 2:顺时针旋转180度; 3:顺时针旋转270度
func (p *PrintObject) SetPrintDirectionInPageMode(dir int) {
	if dir >= 0 && dir <= 3 {
		p.Content += "1b54" + p.numToHexStr(dir, 1)
	}
}

// SetAbsoluteVerticalPrintPositionInPageMode [GS $] 在页面模式下设置绝对垂直打印位置
func (p *PrintObject) SetAbsoluteVerticalPrintPositionInPageMode(n int) {
	if n >= 0 && n <= 65535 {
		p.Content += "1d24" + p.numToHexStr(n, 2)
	}
}

// SetRelativeVerticalPrintPositionInPageMode [GS \] 在页面模式下设置相对垂直打印位置.
func (p *PrintObject) SetRelativeVerticalPrintPositionInPageMode(n int) {
	if n >= -32768 && n <= 32767 {
		p.Content += "1d5c" + p.numToHexStr(n, 2)
	}
}

// printAndExitPageMode [FF] 在缓冲区中打印数据并退出页面模式. NT211、NT212不支持页模式
func (p *PrintObject) printAndExitPageMode() {
	p.Content += "0c"
}

// PrintInPageMode [ESC FF] 在缓冲区中打印数据(并保持在页面模式). NT211、NT212不支持页模式
func (p *PrintObject) PrintInPageMode() {
	p.Content += "1b0c"
}

// ClearInPageMode [CAN] 清除缓冲区中的数据(并保持页面模式). NT211、NT212不支持页模式。
func (p *PrintObject) ClearInPageMode() {
	p.Content += "18"
}

// ExitPageMode [ESC S] 退出页面模式并丢弃缓冲区中的数据而不打印.
func (p *PrintObject) ExitPageMode() {
	p.Content += "1b53"
}
