package transliterator

// Only for Cyrillic characters
//
// Data from https://github.com/alexsergivan/transliterator
// and https://en.wikipedia.org/wiki/Cyrillic_script_in_Unicode
// and https://www.gosuslugi.ru/help/faq/foreign_passport/100359

var x004 = []string{
	"Ie",   // 0x00
	"E",    // 0x01
	"Dj",   // 0x02
	"Gj",   // 0x03
	"Ie",   // 0x04
	"Dz",   // 0x05
	"I",    // 0x06
	"Yi",   // 0x07
	"J",    // 0x08
	"Lj",   // 0x09
	"Nj",   // 0x0a
	"Tsh",  // 0x0b
	"Kj",   // 0x0c
	"I",    // 0x0d
	"U",    // 0x0e
	"Dzh",  // 0x0f
	"A",    // 0x10
	"B",    // 0x11
	"V",    // 0x12
	"G",    // 0x13
	"D",    // 0x14
	"E",    // 0x15
	"Zh",   // 0x16
	"Z",    // 0x17
	"I",    // 0x18
	"I",    // 0x19
	"K",    // 0x1a
	"L",    // 0x1b
	"M",    // 0x1c
	"N",    // 0x1d
	"O",    // 0x1e
	"P",    // 0x1f
	"R",    // 0x20
	"S",    // 0x21
	"T",    // 0x22
	"U",    // 0x23
	"F",    // 0x24
	"Kh",   // 0x25
	"Ts",   // 0x26
	"Ch",   // 0x27
	"Sh",   // 0x28
	"Shch", // 0x29
	"Ie",   // 0x2a
	"Y",    // 0x2b
	"",     // 0x2c
	"E",    // 0x2d
	"Iu",   // 0x2e
	"Ia",   // 0x2f
	"a",    // 0x30
	"b",    // 0x31
	"v",    // 0x32
	"g",    // 0x33
	"d",    // 0x34
	"e",    // 0x35
	"zh",   // 0x36
	"z",    // 0x37
	"i",    // 0x38
	"i",    // 0x39
	"k",    // 0x3a
	"l",    // 0x3b
	"m",    // 0x3c
	"n",    // 0x3d
	"o",    // 0x3e
	"p",    // 0x3f
	"r",    // 0x40
	"s",    // 0x41
	"t",    // 0x42
	"u",    // 0x43
	"f",    // 0x44
	"kh",   // 0x45
	"ts",   // 0x46
	"ch",   // 0x47
	"sh",   // 0x48
	"shch", // 0x49
	"ie",   // 0x4a
	"y",    // 0x4b
	"",     // 0x4c
	"e",    // 0x4d
	"iu",   // 0x4e
	"ia",   // 0x4f
	"ie",   // 0x50
	"e",    // 0x51
	"dj",   // 0x52
	"gj",   // 0x53
	"ie",   // 0x54
	"dz",   // 0x55
	"i",    // 0x56
	"yi",   // 0x57
	"j",    // 0x58
	"lj",   // 0x59
	"nj",   // 0x5a
	"tsh",  // 0x5b
	"kj",   // 0x5c
	"i",    // 0x5d
	"u",    // 0x5e
	"dzh",  // 0x5f
}