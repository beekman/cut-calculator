# cut-calculator

A command-line tool that solves the cutting stock problem. Give it a list of pieces you need and the stock you have (or are willing to buy), and it calculates the optimal cutting plan with minimum material waste.

Works for:
- **1D stock** — lumber, pipe, rebar, rope, wallpaper rolls
- **2D stock** — plywood, drywall, pegboard, sheet metal, wallpaper sheets

## Install

```bash
go install github.com/beekman/cut-calculator/cmd/cut-calculator@latest
```

Or download a binary from the [Releases](https://github.com/beekman/cut-calculator/releases) page.

## Quick start

```bash
# 1D: three 8-foot boards on hand, need four 3-foot and two 4-foot pieces
cut-calculator --stock 96:3:onhand --need 36:4 --need 48:2

# 1D: account for blade width (kerf)
cut-calculator --stock 96:3:onhand --need 36:4 --need 48:2 --kerf 0.125

# 1D: mix on-hand offcuts with purchasable stock
cut-calculator --stock 44:1:onhand --stock 96:10 --need 36:4 --need 48:2 --kerf 0.125

# 2D: one 4×8 sheet on hand, need four 24×36 pieces
cut-calculator --stock 48x96:1:onhand --need 24x36:4

# 2D: disable rotation (pieces won't be turned 90°)
cut-calculator --stock 48x96:2:onhand --need 24x36:4 --need 12x12:6 --no-rotate

# 2D: multiple purchasable sheets available if on-hand isn't enough
cut-calculator --stock 48x96:1:onhand --stock 48x96:5 --need 24x36:8 --kerf 0.125

# Plain text output instead of ASCII diagram
cut-calculator --stock 96:3:onhand --need 36:4 --need 48:2 --output text

# JSON output for scripting
cut-calculator --stock 96:3:onhand --need 36:4 --output json

# From a YAML file
cut-calculator -f myproject.yaml

# File provides the base; flags add stock and override kerf
cut-calculator -f myproject.yaml --stock 44:1:onhand --kerf 0.0625
```

## Flags

| Flag | Description |
|------|-------------|
| `-f`, `-file` | Path to a YAML input file |
| `--stock` | Stock piece: `96` (1D) or `48x96` (2D). Append `:N` for count, `:onhand` if you own it |
| `--need` | Required piece: `36` (1D) or `24x36` (2D). Append `:N` for count |
| `--kerf` | Blade kerf width in inches (default: `0`) |
| `--no-rotate` | Disable 90° rotation of pieces in 2D mode |
| `--output` | Output format: `ascii` (default), `text`, `json` |

## Input file format

For anything beyond a quick one-liner, a YAML file is easier to work with:

```yaml
# 1D example
kerf: 0.125

stock:
  - length: 96
    count: 3
    on_hand: true
  - length: 96       # buy more if needed
    on_hand: false

need:
  - length: 36
    count: 4
  - length: 48
    count: 2
```

```yaml
# 2D example
kerf: 0.125
rotate: true

stock:
  - width: 48
    height: 96
    count: 1
    on_hand: true

need:
  - width: 24
    height: 36
    count: 4
  - width: 12
    height: 12
    count: 6
```

When both a file and flags are provided, flags **add** to the stock/need lists and **override** kerf, rotate, and output settings.

## Pattern repeat

For stock with a repeating pattern (wallpaper, patterned fabric), declare a `repeat_distance` so pieces are snapped to pattern boundaries. Alignment gaps are reported as required waste.

```yaml
rotate: false

stock:
  - width: 27
    height: 240
    on_hand: true
    repeat_distance: 24   # pattern repeats every 24"
    repeat_axis: height   # repeat runs vertically

need:
  - width: 27
    height: 96
    count: 3
```

## Join groups

Tag pieces with a `join_group` to let the solver try edge-gluing or joining them as a single combined cut. It runs both options and keeps whichever wastes less.

```yaml
stock:
  - length: 96
    count: 2
    on_hand: true

need:
  - length: 48
    label: left-panel
    join_group: tabletop
  - length: 48
    label: right-panel
    join_group: tabletop
  - length: 36
    count: 1
```

Combined cuts render inline in ASCII output: `|--left-panel(48"):right-panel(48")--|`

## Output formats

**ASCII diagram** (default) — visual layout of each stock piece:

```
Stock #1 (96") — on hand
|--A(36")--|--A(36")--|--B(22")--|  [waste: 1.75"]

Overall waste: 2.1%
Nothing to purchase — all pieces fit within on-hand stock.
```

**Plain text** (`--output text`) — human-readable cut list per piece.

**JSON** (`--output json`) — structured output for scripting or further processing.

## License

GPL-3.0 — see [LICENSE](LICENSE).
