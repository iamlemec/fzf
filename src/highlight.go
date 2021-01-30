package fzf

import (
    // "fmt"
    "bytes"
    "html"

    "github.com/junegunn/fzf/src/util"
)

type HighlightType int

const (
    HighlightHTML HighlightType = iota
    HighlightANSI HighlightType = iota
)

func makeSegments (length int, locs *[]int) (c chan [2]int) {
    c = make(chan [2]int)
    go func() {
        lpos := len(*locs) - 1
        lval := (*locs)[lpos]
        high := false

        s0 := -1
        s1 := -1

        for i := 0; i < length; i++ {
            // fmt.Println(i, lval, high, s0, s1)

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
                    lval = (*locs)[lpos]
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
    slab *util.Slab
    openTag string
    closeTag string
    escapeFunc func (string) string
    out *bytes.Buffer
}

func NewHighlighter (htype HighlightType, slab *util.Slab) *Highlighter {
    high := Highlighter{
        slab: slab,
        openTag: "",
        closeTag: "",
        escapeFunc: func (x string) string { return x },
        out: new(bytes.Buffer),
    }

    if htype == HighlightHTML {
        high.openTag = "<span class=\"match\">"
        high.closeTag = "</span>"
        high.escapeFunc = html.EscapeString
    } else if htype == HighlightANSI {
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

    /*
    fmt.Println(text)
    for i, x := range *locs {
        fmt.Println(i, x)
    }
    */

    high.reset()
    pos := 0

    for off := range makeSegments(length, locs) {
        off0 := off[0]
        off1 := off[1]
        high.write(chars[pos:off0])
        high.open()
        high.write(chars[off0:off1])
        high.close()
        pos = off1
    }

    high.write(chars[pos:])

    return high.string()
}
