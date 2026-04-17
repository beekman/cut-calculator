# Roadmap

## V1 — Complete Product

A user can describe a pegboard job in a YAML file, run the tool,
and get an ASCII diagram showing how to cut each sheet.

### Tasks
- [ ] YAML file input parser
- [ ] Flag + file merge logic
- [ ] 2D solver: recursive guillotine packing
- [ ] 2D solver: rotation support + `--no-rotate` flag
- [ ] ASCII diagram output (1D and 2D)
- [ ] JSON output formatter
- [ ] ASCII becomes default output; plain text via `--output text`
- [ ] Fuzz tests (1D and 2D solvers)
- [ ] Golden file tests (all output formats)

---

## V2+ — Future

- Wallpaper pattern-repeat matching

---

_MVP shipped as v0.1.0 — see `archive/roadmap/MVP.md`_
