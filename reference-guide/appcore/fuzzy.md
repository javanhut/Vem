# fuzzy.go

**Path:** `/home/javanhut/Development/Vem/internal/appcore/fuzzy.go`
**Lines:** 179
**Purpose:** Fuzzy file finder algorithm implementation

## Overview

Implements fuzzy matching for file finder with:
- Case-insensitive matching
- Scoring based on match quality
- Match position tracking for highlighting
- Pre-computed lowercase runes for performance

## Code Blocks

### Lines 1-6: Package and Imports

```go
package appcore
import (
    "sort"
    "unicode"
)
```

### Lines 8-29: Pre-processing

#### fuzzyItem (Lines 8-12)
```go
type fuzzyItem struct {
    Path     string    // Original file path
    Lower    []rune    // Lowercase runes for matching
    Original []rune    // Original runes for case bonus
}
```

#### buildFuzzyItems (Lines 14-29)
Pre-processes file paths into fuzzyItems for efficient matching.

### Lines 31-123: FuzzyScore Algorithm

#### FuzzyScore (Lines 31-59)
Public entry point that calls `fuzzyScoreRunes()`.

#### fuzzyScoreRunes (Lines 61-123)
Core scoring algorithm:

**Matching Phase (Lines 66-80):**
1. Sequential scan through target string
2. Matches pattern characters in order
3. Returns 0 if not all characters matched

**Scoring Phase (Lines 82-121):**

| Bonus | Points | Condition |
|-------|--------|-----------|
| Base match | +10 | Each matched character |
| Consecutive | +15 | Match follows previous match |
| Start of string | +10 | Match at index 0 |
| Word boundary | +5 | After `/`, `_`, `-`, ` `, `.` |
| Case match | +2 | Exact case match |
| Gap penalty | -N | Characters between matches |
| Short path | +bonus | `1000 - len(target)` |

### Lines 125-173: PerformFuzzyMatch

#### PerformFuzzyMatch (Lines 125-173)
Performs fuzzy matching on list of items:
1. Returns all items if pattern is empty
2. Scores each item using `fuzzyScoreRunes()`
3. Filters items with score > 0
4. Sorts by score descending
5. Limits to `maxResults`

### Lines 175-178: isWordBoundary

```go
func isWordBoundary(r rune) bool {
    return r == '/' || r == '_' || r == '-' || r == ' ' || r == '.' || unicode.IsUpper(r)
}
```

## Known Issues / Potential Bugs

1. **Line 175-178: isWordBoundary unused**
   - Function is defined but never called
   - The logic is duplicated inline in fuzzyScoreRunes

## Dead/Unused Code

1. **isWordBoundary function** (Lines 175-178) - Not used

## Integration Points

- Called from fuzzy finder in app.go
- Works with `filesystem.FindFiles()` for path listing
- Results displayed via fuzzy finder UI

## Scoring Examples

```
Pattern: "bf"
Target:  "internal/buffer/buffer.go"

Matches: b(uffer/)b(uffer).go  -> "buffer" not matched, no match

Target:  "internal/editor/buffer.go"
Matches: (b)u(f)fer.go
Score = 10 + 10           = 20 (base)
      + 0                 = 0  (not consecutive)
      + 0 + 5             = 5  (word boundary for 'b')
      + 0                 = 0  (no case match)
      + 1000 - 27         = 973 (short path bonus)
      - (18-11-1)         = -6 (gap penalty)
Total: ~992
```

---
*Last Updated: Reference guide creation*
