package main

import (
	"strconv"
	"strings"
)

// ParseHotkey разбирает строку вида "Alt+Q", "Ctrl+Shift+F", "F10".
func ParseHotkey(s string) (mods, vk uintptr, label string, ok bool) {
	parts := strings.Split(strings.ReplaceAll(s, " ", ""), "+")
	if len(parts) == 0 || parts[len(parts)-1] == "" {
		return 0, 0, "", false
	}
	var labels []string
	for _, p := range parts[:len(parts)-1] {
		switch strings.ToLower(p) {
		case "alt":
			mods |= 0x0001
			labels = append(labels, "Alt")
		case "ctrl", "control", "ctl":
			mods |= 0x0002
			labels = append(labels, "Ctrl")
		case "shift":
			mods |= 0x0004
			labels = append(labels, "Shift")
		case "win":
			mods |= 0x0008
			labels = append(labels, "Win")
		default:
			return 0, 0, "", false
		}
	}
	key := strings.ToUpper(parts[len(parts)-1])
	// русские буквы → та же клавиша на английской раскладке
	ru := map[string]string{"Й": "Q", "Ц": "W", "У": "E", "К": "R", "Е": "T", "Н": "Y", "Г": "U", "Ш": "I", "Щ": "O", "З": "P",
		"Ф": "A", "Ы": "S", "В": "D", "А": "F", "П": "G", "Р": "H", "О": "J", "Л": "K", "Д": "L",
		"Я": "Z", "Ч": "X", "С": "C", "М": "V", "И": "B", "Т": "N", "Ь": "M"}
	if v, ok := ru[key]; ok {
		key = v
	}
	named := map[string]uintptr{"PAUSE": 0x13, "INSERT": 0x2D, "INS": 0x2D, "DELETE": 0x2E, "DEL": 0x2E,
		"HOME": 0x24, "END": 0x23, "PAGEUP": 0x21, "PGUP": 0x21, "PAGEDOWN": 0x22, "PGDN": 0x22,
		"SPACE": 0x20, "ПРОБЕЛ": 0x20, "`": 0xC0, "~": 0xC0, "Ё": 0xC0, "TAB": 0x09}
	switch {
	case len(key) == 1 && key[0] >= 'A' && key[0] <= 'Z', len(key) == 1 && key[0] >= '0' && key[0] <= '9':
		vk = uintptr(key[0])
	case strings.HasPrefix(key, "F") && len(key) >= 2:
		n, err := strconv.Atoi(key[1:])
		if err != nil || n < 1 || n > 24 {
			return 0, 0, "", false
		}
		vk = uintptr(0x70 + n - 1)
	case strings.HasPrefix(key, "NUM") && len(key) == 4 && key[3] >= '0' && key[3] <= '9':
		vk = uintptr(0x60 + key[3] - '0')
	default:
		v, ok := named[key]
		if !ok {
			return 0, 0, "", false
		}
		vk = v
	}
	if mods == 0 && vk >= 'A' && vk <= 'Z' {
		return 0, 0, "", false // одна буква без Alt/Ctrl помешает печатать
	}
	labels = append(labels, key)
	return mods, vk, strings.Join(labels, "+"), true
}
