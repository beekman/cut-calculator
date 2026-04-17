package input

// Merge combines a file-based config with flag-based overrides.
// File provides the base; flags add to stock/need lists and override scalar settings.
func Merge(file, flags *Config) *Config {
	out := &Config{
		Kerf:   file.Kerf,
		Rotate: file.Rotate,
		// file rotate defaults to true; if the file didn't set it, keep true
		OutputFormat: file.OutputFormat,
	}

	// scalar flags override file values when explicitly set
	if flags.Kerf != 0 {
		out.Kerf = flags.Kerf
	}
	if !flags.Rotate {
		out.Rotate = false
	}
	if flags.OutputFormat != "" {
		out.OutputFormat = flags.OutputFormat
	}

	// stock and need: flags add to the file list
	out.Stock = append(out.Stock, file.Stock...)
	out.Stock = append(out.Stock, flags.Stock...)

	// re-label need items sequentially across both sources
	allNeed := append(file.Need, flags.Need...)
	labels := labelSeq()
	for i := range allNeed {
		allNeed[i].Label = labels()
	}
	out.Need = allNeed

	return out
}
