package sr

import (
	"fmt"
)

const (
	ModelRealEsrgan = "real-esrgan"
	ModelRealCugan  = "real-cugan"
)

var (
	ModelNames = []string{
		ModelRealEsrgan,
		ModelRealCugan,
	}
	Models map[string]Model
)

func init() {
	Models = map[string]Model{
		ModelRealEsrgan: {newModelMust(ModelRealEsrgan)},
		ModelRealCugan:  {newModelMust(ModelRealCugan)},
	}
}

func newModel(name string) (model, error) {
	switch name {
	case ModelRealEsrgan:
		model, err := newEsrgan()
		return model, err
	case ModelRealCugan:
		model, err := newCugan()
		return model, err
	}
	return nil, fmt.Errorf("model %s not found", name)
}
func newModelMust(name string) model {
	model, err := newModel(name)
	if err != nil {
		log.Panicw("failed to create model",
			"model", name,
			"err", err,
		)
		return nil
	}
	return model
}
