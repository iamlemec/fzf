package fzf

import (
    "bytes"
    "html"

    "github.com/junegunn/fzf/src/util"
)

type HighlightType int

const (
    HighlightHTML HighlightType = iota
    HighlightANSI HighlightType = iota
)

func PrintHighlight (item *Item, pattern *Pattern, slab *util.Slab, htype HighlightType) (string) {
    text := item.AsString(false)
    _, offsets, locs := pattern.MatchItem(item, true, slab)

    open := ""
    close := ""
    if htype == HighlightHTML {
        open = "<span class=\"match\">"
        close = "</span>"
    } else if htype == HighlightANSI {
        open = "\x1b[31m"
        close = "\x1b[0m"
    }

    lpos := len(*locs) - 1
    lval := (*locs)[lpos]

    high := false
    pos := 0

    out := new(bytes.Buffer)
    buf := new(bytes.Buffer)

    for _, off := range offsets {
        off0 := int(off[0])
        off1 := int(off[1])

        out.WriteString(html.EscapeString(text[pos:off0]))

        for i := off0; i < off1; i++ {
            if lpos >= 0 {
                lval = (*locs)[lpos]
            } else {
                lval = -1
            }

            if i == lval {
                lpos--
            }

            if high && i != lval {
                out.WriteString(html.EscapeString(buf.String()))
                out.WriteString(close)
                high = false
                buf.Reset()
            }

            if !high && i == lval {
                out.WriteString(html.EscapeString(buf.String()))
                out.WriteString(open)
                high = true
                buf.Reset()
            }

            buf.WriteByte(text[i])
        }

        if high {
            out.WriteString(html.EscapeString(buf.String()))
            out.WriteString(close)
            high = false
            buf.Reset()
        }

        pos = off1
    }

    buf.WriteString(text[pos:])
    out.WriteString(html.EscapeString(buf.String()))

    return out.String()
}
