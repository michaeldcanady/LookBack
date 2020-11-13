// +build darvin

package main

import(

)

const(
  PATHSEPARATOR = '/'
  PATHLISTSEPARATOR = ':'
  UNIT = 1000
)

func ByteCountSI(b int64) string {
    if b < UNIT {
        return fmt.Sprintf("%d B", b)
    }
    div, exp := int64(UNIT), 0
    for n := b / UNIT; n >= UNIT; n /= UNIT {
        div *= UNIT
        exp++
    }
    return fmt.Sprintf("%.1f %cB",
        float64(b)/float64(div), "kMGTPE"[exp])
}
