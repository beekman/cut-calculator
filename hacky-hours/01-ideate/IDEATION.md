# Ideation

Free-writing space — no rules. Capture everything here.

---

## Initial Brain Dump

A CLI tool written in Go that solves the **cutting stock problem**:
given stock material of fixed size(s), find the optimal way to cut pieces
of the required sizes with minimum waste.

- Supports **1D cutting** (linear material: pipes, lumber, rebar, extrusions)
- Supports **2D cutting** (sheet material: plywood, glass, sheet metal, fabric)
- Goal: minimize offcuts / wasted material
- Use the best suitable optimization algorithm for the problem

---

## Open Questions

- Who is the primary user? (DIYer, tradesperson, fabricator, engineer?)
- CLI only, or also a library/API?
- Input format? (interactive prompt, flags, file?)
- Output format? (plain text, JSON, visual diagram?)
- How are stock sizes specified? (one fixed size, multiple sizes, unlimited?)
- Blade/kerf width — does it matter?
- Should it rank multiple solutions or just return the best?
- Any constraints on algorithm runtime / problem size?
