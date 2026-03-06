package variant

type VariantModel struct {
	Source   string `json:"src"`
	Format   string `json:"format"`
	Codecs   string `json:"codecs,omitempty"`
	Manbytes uint16 `json:"manbytes"`
	Width    uint16 `json:"w"`
	Height   uint16 `json:"h"`
}

type VariantLayersModel struct {
	Original VariantModel     `json:"original"`
	Fallback *VariantModel    `json:"fallback,omitempty"`
	Variants [][]VariantModel `json:"variants"`
}
