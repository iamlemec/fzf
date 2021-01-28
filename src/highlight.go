package fzf

import (
    "fmt"
    "html"

    "github.com/junegunn/fzf/src/util"
)

type HighlightType int

const (
    HighlightHTML HighlightType = iota
    HighlightANSI HighlightType = iota
)

func PrintHighlight (item *Item, pattern *Pattern, slab *util.Slab, htype HighlightType) {
    text := item.AsString(false)
    _, offsets, locs := pattern.MatchItem(item, true, slab)

    /*
    for _, off := range offsets {
        fmt.Print("(", off[0], off[1], ") ")
    }
    for _, p := range *locs {
        fmt.Print(p, " ")
    }
    fmt.Print("\n")
    */

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
    buf := ""

    pos := 0
    out := ""

    for _, off := range offsets {
        off0 := int(off[0])
        off1 := int(off[1])

        out += html.EscapeString(text[pos:off0])

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
                out += html.EscapeString(buf)
                out += close
                high = false
                buf = ""
            }

            if !high && i == lval {
                out += html.EscapeString(buf)
                out += open
                high = true
                buf = ""
            }

            buf += string(text[i])
        }

        if high {
            out += html.EscapeString(buf)
            out += close
            high = false
            buf = ""
        }

        pos = off1
    }

    buf += text[pos:]
    out += html.EscapeString(buf)

    fmt.Println(out)
}

// fmt.Print("\x1b[31m")
// fmt.Print("\x1b[0m")
