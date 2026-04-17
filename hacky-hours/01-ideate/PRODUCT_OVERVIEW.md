# Product Overview

## Who
A DIYer working with wood — building a deck, covering a shed with pegboard,
framing a room. Not a professional fabricator. Comfortable with basic measurements
but not an optimization expert. Wants to know "how do I cut this with as little
waste as possible?" before they go buy materials or make the first cut.

## What
A command-line tool written in Go that solves the **cutting stock problem**.
Given a list of required piece sizes and available stock material dimensions,
it calculates the optimal cutting plan to produce all required pieces with
minimum material waste.

Supports:
- **1D mode**: linear stock (lumber, pipe, rebar, rope, wallpaper rolls)
- **2D mode**: sheet stock (plywood, pegboard, drywall, sheet metal, wallpaper sheets)

**Stock inventory:** users can specify a mix of stock pieces in different sizes,
including leftover offcuts from previous projects — not just a single standard size.

**Out of scope (for now):** wallpaper pattern-repeat matching.

## Where
CLI tool — runs locally, no server or internet connection required.
Distributed as a single compiled binary.

**Output modes:**
- Visual ASCII diagram (default) — shows each stock piece with cut positions marked
- Plain text cut list — labeled cuts per piece, human-readable
- JSON export — structured output for piping or scripting (`--output json`)

## When
No hard deadline. MVP scope is intentionally small — solve the core problem well
before expanding.

**Input size constraints (at least for now):**
- 2D stock: up to 10 ft × 10 ft
- 1D stock: typical lumber lengths (8 ft, 10 ft, 12 ft, 16 ft range)
- Problem scale: small — a typical DIY project, not a production run

## Why
Cutting material by hand-calculation or guesswork wastes money and time.
This tool automates the optimization so users get the best possible cutting
plan without having to think through every permutation themselves.

---

## Constraints & Values

### Licensing
GPL-3.0 (copyleft open source). Any algorithms or dependencies chosen must
be compatible with GPL-3.0.

### Privacy
No user data collected. Fully offline — no network calls, no telemetry.

### Infrastructure
None. Single compiled binary. No server, no database, no cloud dependency.
