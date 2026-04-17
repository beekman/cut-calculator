# Iteration Log — post-v1.0.0

## Raw Brain-Dump

### Pattern-repeat matching (wallpaper / fabric / tiled materials)
When cutting materials with a repeating pattern (wallpaper, patterned fabric, some tile),
each required piece needs to start at a pattern-repeat boundary on the stock.

**Design decisions:**
- `repeat_distance` is a property of the **stock piece** (global to that material roll/sheet)
- In 2D, the repeat applies to **one axis** — user specifies `repeat_axis: height` or `repeat_axis: width`
- In 1D, axis is implicit (the only axis)
- Pieces that don't need alignment use stock with no `repeat_distance` set
- Almost always paired with `--no-rotate` (pattern repeat implies orientation matters)

### Prefer joining pieces (combined-cut optimization)
Two or more required pieces tagged with the same `join_group` will be joined in the final
product. The solver tries cutting them as a single combined piece; if that wastes less
material, the combined cut is used.

**Design decisions:**
- `join_group` is a label on `RequiredPiece`
- Output: combined piece shows as one rectangle with a dashed dividing line and both labels
- Works in 1D and 2D

## Triage

| Item | Category | Design doc affected |
|------|----------|-------------------|
| Pattern-repeat matching | Next milestone | BUSINESS_LOGIC.md, ARCHITECTURE.md |
| Join groups | Next milestone | BUSINESS_LOGIC.md, ARCHITECTURE.md |
