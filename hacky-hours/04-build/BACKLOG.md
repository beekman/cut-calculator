# Backlog

## V1.1 — Pattern Repeat + Join Groups

### Pattern Repeat
- [x] Add `RepeatDistance` and `RepeatAxis` fields to `StockPiece` model (`internal/model/types.go`)
- [ ] Parse `repeat_distance` / `repeat_axis` in YAML input (`internal/input/file.go`)
- [ ] Update 1D solver: snap piece placement to repeat boundaries; report alignment gaps as required waste
- [ ] Update 2D solver: snap placement on the specified axis to repeat boundaries
- [ ] Update ASCII output: show repeat boundary markers (dashed lines at each repeat interval)
- [ ] Add unit tests for repeat-boundary placement logic (1D and 2D)
- [ ] Add golden file tests for repeat output

### Join Groups
- [x] Add `JoinGroup` field to `RequiredPiece` model (`internal/model/types.go`)
- [ ] Parse `join_group` in YAML input (`internal/input/file.go`)
- [ ] Implement join-group pre-processor: for each group, generate a candidate combined piece and try both options (combined vs. individual); keep the lower-waste result
- [ ] Update ASCII output: combined pieces render as one rectangle with a dashed dividing line and both labels
- [ ] Add unit tests for join-group pre-processor
- [ ] Add golden file tests for join-group output
