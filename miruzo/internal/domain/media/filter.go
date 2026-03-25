package media

import "slices"

var supportedFormats = []ImageFormat{
	ImageFormatGIF,
	ImageFormatJPEG,
	ImageFormatPNG,
	ImageFormatWebP,
}

type ImageFormatSet map[ImageFormat]struct{}

var defaultAllowedFormats = ImageFormatSet{
	ImageFormatGIF:  {},
	ImageFormatPNG:  {},
	ImageFormatJPEG: {},
	ImageFormatWebP: {},
}

type VariantFilterOptions struct {
	IncludeFormatSet ImageFormatSet
	KeepFallback     bool
}

// ComputeAllowedFormats returns the effective allowed formats by removing
// excluded formats from supported formats.
func ComputeAllowedFormats(disallowed []ImageFormat) ImageFormatSet {
	if disallowed == nil {
		return defaultAllowedFormats
	}

	allowedFormats := make(ImageFormatSet, 4)
	for _, format := range supportedFormats {
		if !slices.Contains(disallowed, format) {
			allowedFormats[format] = struct{}{}
		}
	}

	return allowedFormats
}

// FilterWith returns variants that remain after applying IncludeFormatSet.
//
// IncludeFormatSet defines formats to keep.
// If KeepFallback is true, fallback-layer variants are always kept regardless
// of IncludeFormatSet.
func (v Variants) FilterWith(options VariantFilterOptions) Variants {
	if len(v) == 0 {
		return nil
	}

	var result Variants
	for _, variant := range v {
		// Fallback layers bypass filtering so legacy clients still work.
		if options.KeepFallback && variant.IsFallback() {
			result = append(result, variant)
			continue
		}

		if _, ok := options.IncludeFormatSet[variant.Format]; ok {
			result = append(result, variant)
		}
	}

	return result
}
