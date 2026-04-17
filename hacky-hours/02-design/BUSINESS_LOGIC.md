# Business Logic

## Problem Definition

Given:
- A **stock inventory** — one or more pieces of raw material (boards, sheets, offcuts)
- A **cut list** — the required output pieces and their sizes

Produce:
- An **assignment** of every required piece to a stock piece
- A **cut sequence** for each stock piece
- A **waste report** showing offcuts and total material efficiency
- A recommendation for **additional stock to purchase** if the inventory is insufficient

The optimizer must guarantee every required piece is produced. Minimizing waste is the primary objective. Using existing stock before suggesting new purchases is the secondary objective.

---

## Key Concepts

### Kerf Width
A saw blade removes a small amount of material with each cut — typically 1/8" (3mm) for a circular saw. This is called the **kerf**. For a single cut it's negligible, but across 20 cuts on one board it adds up to 2.5". The optimizer accounts for kerf on every cut.

Kerf width is user-configurable. Default: `0` (no kerf). Users set it once per session or per stock material type.

### 1D vs 2D Mode
The tool auto-detects mode based on input:
- All stock and required pieces have **one dimension** → 1D mode
- Any stock or required piece has **two dimensions** → 2D mode

Mixed 1D/2D input in a single session is not supported.

### Guillotine Cuts (2D only)
A guillotine cut goes straight across the full width or height of a panel — which is how a circular saw, table saw, or track saw physically works. The optimizer only produces guillotine-compatible layouts. Non-guillotine arrangements (L-shapes, irregular nesting) are never suggested.

### Rotation (2D only)
By default, required pieces may be rotated 90° if it produces a better fit. Rotation can be disabled globally (`--no-rotate`) for materials where orientation matters (e.g. wood grain, pegboard hole alignment, wallpaper).

### Mixed Stock
Stock pieces are tagged as either:
- **On hand** — already owned; the optimizer uses these first at no cost
- **Available** — a standard size the user can purchase; the optimizer uses these only after exhausting on-hand stock

The optimizer reports the minimum number of "available" pieces to buy to complete the cut list.

---

## 1D Algorithm

**Problem class:** 1D bin packing / cutting stock

**Algorithm: Exact Dynamic Programming**

For DIY-scale inputs (≤ ~200 required pieces, ≤ ~50 stock pieces), an exact DP solution is computationally feasible and guarantees the optimal answer.

**Approach:**

1. For each stock piece, solve a **bounded knapsack** subproblem: find the subset of required pieces that fits within the stock length (accounting for kerf between cuts) and maximizes used length.
2. Assign required pieces to stock pieces using a **branch-and-bound** search over the full assignment, guided by the knapsack values.
3. On-hand stock is tried before available (purchasable) stock.
4. If available stock is needed, find the minimum number of additional pieces required.

**Waste metric:** `(total stock length used − total required length) / total stock length used`

```mermaid
flowchart TD
    A[Input: stock inventory + cut list] --> B[Sort stock: on-hand first, largest first]
    B --> C[For each stock piece: solve knapsack subproblem]
    C --> D[Branch-and-bound assignment search]
    D --> E{All pieces assigned?}
    E -- yes --> F[Return optimal plan]
    E -- no --> G[Add minimum purchasable stock]
    G --> D
```

---

## 2D Algorithm

**Problem class:** 2D guillotine bin packing

**Algorithm: Recursive Guillotine with Branch-and-Bound**

For DIY-scale 2D inputs, a recursive guillotine approach with pruning is tractable and produces optimal or near-optimal results.

**Approach:**

1. Sort required pieces largest-first (by area).
2. For each stock sheet, attempt to pack remaining pieces using **recursive guillotine subdivision**:
   - Place the next piece in the top-left corner of the current region.
   - Make a guillotine cut (horizontal or vertical — try both).
   - Recursively pack each resulting sub-region.
   - Backtrack if no valid placement exists.
3. Rotation is tried at each placement step (unless `--no-rotate`).
4. On-hand sheets are tried before purchasable sheets.
5. Branch-and-bound pruning skips branches where remaining area is provably insufficient.

**Kerf** is applied to each guillotine cut, reducing the usable area of sub-regions accordingly.

**Waste metric:** `(total sheet area used − total required piece area) / total sheet area used`

```mermaid
flowchart TD
    A[Input: stock sheets + required pieces] --> B[Sort pieces by area descending]
    B --> C[Try on-hand sheets first]
    C --> D[Recursive guillotine packing]
    D --> E{Piece fits in region?}
    E -- yes --> F[Place piece, split region, recurse]
    E -- no --> G[Try rotation]
    G --> H{Rotated piece fits?}
    H -- yes --> F
    H -- no --> I[Backtrack]
    I --> D
    F --> J{All pieces placed?}
    J -- yes --> K[Return plan for this sheet]
    J -- no --> L[Move to next sheet]
    L --> D
    K --> M{Any pieces remaining?}
    M -- yes --> N[Add purchasable sheet, repeat]
    M -- no --> O[Return complete plan]
```

---

## Output Contract

Every successful run produces:

| Field | Description |
|-------|-------------|
| `assignments` | Each required piece → which stock piece it comes from |
| `cut_sequences` | Per stock piece: ordered list of cuts with positions |
| `offcuts` | Leftover pieces per stock, with dimensions |
| `waste_pct` | Overall material waste percentage |
| `purchased` | List of additional stock pieces the user needs to buy |

If no valid plan exists (e.g. a required piece is larger than any available stock), the tool reports which pieces could not be fit and why — it never silently drops pieces.

---

## Edge Cases

| Case | Behavior |
|------|----------|
| Required piece larger than all stock | Error: report which piece(s) cannot be cut, suggest stock size needed |
| Empty cut list | Error: nothing to do |
| Empty stock inventory | Error: no stock provided |
| Required piece exactly equals stock size (no kerf room) | Assign as a zero-waste cut; no saw cut needed |
| Duplicate required sizes | Treated as separate pieces — each must be assigned |
| Stock piece already smaller than smallest required piece | Skip it; report as unusable offcut |
