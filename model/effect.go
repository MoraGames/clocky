package model

type (
	//TODO: Rework this shit
	Effect struct {
		Name       string
		Parameters []*EffectParameter
	}

	EffectParameter struct {
		Key   string
		Value int
	}
)
