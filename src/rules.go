package main

// Чтение файлов правил, разбор их структуры и поиск.
// Этот файл не зависит от Windows — его можно тестировать на любой системе.

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type Item struct {
	Kind  string // "li", "note", "p", "sub"
	Label string // для заметок: "Примечание", "Исключение" ...
	Note  string // вид заметки: prim, utoch, iskl, poyas, other
	Text  string
}

type Point struct {
	Num    string
	Text   string
	Pun    []string
	Items  []Item
	Doc    *Doc
	Sec    *Section
	Index  int    // порядковый номер в документе
	search string // нормализованный текст для поиска
}

// Block — либо пункт, либо служебная строка раздела (абзац, подзаголовок, список).
type Block struct {
	Point *Point
	Item  *Item
}

type Section struct {
	Title  string
	Blocks []Block
}

type Doc struct {
	File     string
	Title    string
	Short    string
	Updated  string
	Source   string
	Sections []*Section
	Points   []*Point
}

var (
	pointRe = regexp.MustCompile(`^(\d+(?:\.\d+)+)\.?\s+(.+)$`)
	noteRe  = regexp.MustCompile(`(?i)^(Примечание|Уточнение|Исключение|Пояснение|Примеры|Пример|Нарушением считается|Запрещены|Агрессивными действиями считается|IC причины для убийства)(\s*\d*)\s*:\s*`)
)

// LoadDocs читает все .txt из папки правил (по алфавиту).
func LoadDocs(dir string) ([]*Doc, []string) {
	var docs []*Doc
	var errs []string
	files, _ := filepath.Glob(filepath.Join(dir, "*.txt"))
	sort.Strings(files)
	if len(files) == 0 {
		errs = append(errs, "В папке rules нет файлов с правилами (*.txt).")
	}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			errs = append(errs, "Не удалось прочитать "+filepath.Base(f)+": "+err.Error())
			continue
		}
		d := ParseDoc(string(b))
		d.File = filepath.Base(f)
		if d.Title == "" {
			d.Title = strings.TrimSuffix(d.File, filepath.Ext(d.File))
		}
		if d.Short == "" {
			d.Short = d.Title
		}
		docs = append(docs, d)
	}
	return docs, errs
}

// ParseDoc разбирает текст файла: шапка до строки "---", затем правила.
func ParseDoc(src string) *Doc {
	src = strings.TrimPrefix(src, "\uFEFF")
	src = strings.ReplaceAll(src, "\u200b", "")
	src = strings.ReplaceAll(src, "\r\n", "\n")
	d := &Doc{}
	lines := strings.Split(src, "\n")
	body := lines
	for i, l := range lines {
		if strings.TrimSpace(l) == "---" {
			for _, h := range lines[:i] {
				h = strings.TrimSpace(h)
				if h == "" || strings.HasPrefix(h, "#") {
					continue
				}
				k, v, ok := strings.Cut(h, ":")
				if !ok {
					continue
				}
				v = strings.TrimSpace(v)
				switch strings.ToLower(strings.TrimSpace(k)) {
				case "название":
					d.Title = v
				case "кратко":
					d.Short = v
				case "обновлено":
					d.Updated = v
				case "ссылка":
					d.Source = v
				}
			}
			body = lines[i+1:]
			break
		}
	}

	var sec *Section
	var pt *Point
	ensure := func() {
		if sec == nil {
			sec = &Section{}
			d.Sections = append(d.Sections, sec)
		}
	}
	for _, raw := range body {
		line := strings.TrimSpace(raw)
		if line == "" {
			pt = nil
			continue
		}
		if strings.HasPrefix(line, "## ") {
			sec = &Section{Title: strings.TrimSpace(line[3:])}
			d.Sections = append(d.Sections, sec)
			pt = nil
			continue
		}
		if strings.HasPrefix(line, "### ") {
			ensure()
			sec.Blocks = append(sec.Blocks, Block{Item: &Item{Kind: "sub", Text: strings.TrimSpace(line[4:])}})
			pt = nil
			continue
		}
		if m := pointRe.FindStringSubmatch(line); m != nil {
			ensure()
			text, pun := SplitPun(m[2])
			pt = &Point{Num: m[1], Text: text, Pun: pun, Doc: d, Sec: sec, Index: len(d.Points)}
			sec.Blocks = append(sec.Blocks, Block{Point: pt})
			d.Points = append(d.Points, pt)
			continue
		}
		ensure()
		it := classify(line)
		if pt != nil {
			pt.Items = append(pt.Items, it)
		} else {
			sec.Blocks = append(sec.Blocks, Block{Item: &it})
		}
	}
	for _, p := range d.Points {
		parts := []string{p.Num, p.Text, strings.Join(p.Pun, " ")}
		for _, it := range p.Items {
			parts = append(parts, it.Label+" "+it.Text)
		}
		p.search = Norm(strings.Join(parts, "\n"))
	}
	return d
}

func classify(line string) Item {
	if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "• ") {
		return Item{Kind: "li", Text: strings.TrimSpace(strings.TrimLeft(line, "-• "))}
	}
	if m := noteRe.FindStringSubmatchIndex(line); m != nil {
		word := line[m[2]:m[3]]
		label := strings.TrimSpace(line[m[2]:m[5]])
		return Item{Kind: "note", Label: label, Note: noteKind(word), Text: line[m[1]:]}
	}
	return Item{Kind: "p", Text: line}
}

func noteKind(w string) string {
	w = strings.ToLower(w)
	switch {
	case strings.HasPrefix(w, "примечание"):
		return "prim"
	case strings.HasPrefix(w, "уточнение"):
		return "utoch"
	case strings.HasPrefix(w, "исключение"):
		return "iskl"
	case strings.HasPrefix(w, "пояснение"):
		return "poyas"
	}
	return "other"
}

// SplitPun отделяет наказания (всё после первого «|»).
func SplitPun(s string) (string, []string) {
	i := strings.Index(s, "|")
	if i < 0 {
		return s, nil
	}
	var pun []string
	for _, p := range strings.Split(s[i+1:], "|") {
		p = strings.TrimSuffix(strings.TrimSpace(p), ".")
		if p != "" {
			pun = append(pun, p)
		}
	}
	return strings.TrimSpace(s[:i]), pun
}

// Norm — нижний регистр и ё→е, длина в рунах сохраняется.
func Norm(s string) string {
	r := []rune(s)
	for i, c := range r {
		c = unicode.ToLower(c)
		if c == 'ё' {
			c = 'е'
		}
		r[i] = c
	}
	return string(r)
}

var numQueryRe = regexp.MustCompile(`^\d+(\.\d+)*\.?$`)

// Search ищет пункты, где встречаются все слова запроса (или номер пункта).
func Search(docs []*Doc, q string) (found []*Point, words []string) {
	nq := strings.TrimSpace(Norm(q))
	if nq == "" {
		return nil, nil
	}
	words = strings.Fields(nq)
	if numQueryRe.MatchString(nq) {
		n := strings.TrimSuffix(nq, ".")
		for _, d := range docs {
			for _, p := range d.Points {
				if p.Num == n || strings.HasPrefix(p.Num, n+".") {
					found = append(found, p)
				}
			}
		}
		return found, []string{n}
	}
	for _, d := range docs {
		for _, p := range d.Points {
			ok := true
			for _, w := range words {
				if !strings.Contains(p.search, w) {
					ok = false
					break
				}
			}
			if ok {
				found = append(found, p)
			}
		}
	}
	return found, words
}

// FavKey — ключ избранного, устойчивый к перестановке пунктов.
func FavKey(p *Point) string {
	t := []rune(p.Text)
	if len(t) > 40 {
		t = t[:40]
	}
	return p.Doc.File + "|" + p.Num + "|" + string(t)
}

// PointText — текст пункта для копирования.
func PointText(p *Point) string {
	var b strings.Builder
	b.WriteString("п. " + p.Num + " (" + p.Doc.Title + "): " + p.Text)
	if len(p.Pun) > 0 {
		b.WriteString(" — " + strings.Join(p.Pun, " | "))
	}
	for _, it := range p.Items {
		b.WriteString("\r\n")
		switch it.Kind {
		case "li":
			b.WriteString("• " + it.Text)
		case "note":
			b.WriteString(it.Label + ": " + it.Text)
		default:
			b.WriteString(it.Text)
		}
	}
	return b.String()
}
