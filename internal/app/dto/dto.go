package dto

import (
	"fmt"
)

type Asset struct {
	Kind               string
	Name               string
	Variants           map[string]Variant
	UpscalableVariants []Variant
}
type Variant struct {
	Kind     string
	Name     string
	Variant  string
	Filename string
}

func (a Asset) String() string {
	return fmt.Sprintf("%s/%s", a.Kind, a.Name)
}
func (v Variant) String() string {
	return fmt.Sprintf("%s/%s/%s", v.Kind, v.Name, v.Variant)
}
func (v Variant) Path() string {
	return fmt.Sprintf("%s/%s/%s", v.Variant, v.Kind, v.Filename)
}
func (v Variant) URL(api string) string {
	return fmt.Sprintf("%s/asset/variants/%s/%s/%s", api, v.Kind, v.Name, v.Variant)
}
func (v Variant) FileURL(api string) string {
	return fmt.Sprintf("%s/asset/variants/%s/%s/%s/file", api, v.Kind, v.Name, v.Variant)
}
