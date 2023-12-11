package model

type (
	Effect struct {
		Name       string
		Parameters []*EffectParameter
	}

	EffectParameter struct {
		Key   string
		Value int
	}
)
