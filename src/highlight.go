package fzf

import (
    "sort"
    "bytes"
    "html"

    "github.com/junegunn/fzf/src/util"
)

func makeSegments (length int, locs *[]int) (c chan [2]int) {
    c = make(chan [2]int)
    go func() {
        slocs := make([]int, len(*locs))
        copy(slocs, *locs)
        sort.Sort(sort.Reverse(sort.IntSlice(slocs)))

        lpos := len(slocs) - 1
        lval := slocs[lpos]
        high := false

        s0 := -1
        s1 := -1

        for i := 0; i < length; i++ {
            if high {
                if i != lval {
                    s1 = i
                    high = false
                    c <- [2]int{s0, s1}
                    s0 = -1
                }
            } else {
                if i == lval {
                    s0 = i
                    high = true
                }
            }

            if i == lval {
                lpos--
                if lpos >= 0 {
                    lval = slocs[lpos]
                } else {
                    lval = -1
                }
            }
        }

        if s0 != -1 {
            c <- [2]int{s0, length}
        }

        close(c)
    }()
    return c
}

type Highlighter struct {
    openTag string
    closeTag string
    escapeFunc func (string) string
    slab *util.Slab
    out *bytes.Buffer
}

func NewHighlighter (htype highlightType, slab *util.Slab) *Highlighter {
    high := Highlighter{
        openTag: "",
        closeTag: "",
        escapeFunc: func (x string) string { return x },
        slab: slab,
        out: new(bytes.Buffer),
    }

    if htype == highlightHtml {
        high.openTag = "<span class=\"match\">"
        high.closeTag = "</span>"
        high.escapeFunc = html.EscapeString
    } else if htype == highlightAnsi {
        high.openTag = "\x1b[31m"
        high.closeTag = "\x1b[0m"
    }

    return &high
}

func (high *Highlighter) reset() {
    high.out.Reset()
}

func (high *Highlighter) open() {
    high.out.WriteString(high.openTag)
}

func (high *Highlighter) close() {
    high.out.WriteString(high.closeTag)
}

func (high *Highlighter) write(text []rune) {
    high.out.WriteString(high.escapeFunc(string(text)))
}

func (high *Highlighter) string() (string) {
    return high.out.String()
}

func (high *Highlighter) RenderHighlight (item *Item, pattern *Pattern) (string) {
    _, _, locs := pattern.MatchItem(item, true, high.slab)

    text := item.AsString(false)
    chars := []rune(text)
    length := len(chars)

    high.reset()
    pos := 0

    for off := range makeSegments(length, locs) {
        high.write(chars[pos:off[0]])
        high.open()
        high.write(chars[off[0]:off[1]])
        high.close()
        pos = off[1]
    }

    high.write(chars[pos:])

    return high.string()
}
