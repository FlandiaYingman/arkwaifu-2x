package dto

import (
	"fmt"
	"net/url"
	"path/filepath"
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
func (v Variant) ActualPath(dirPath string) string {
	return filepath.Join(dirPath, v.Path())
}
func (v Variant) URL(api string) string {
	return fmt.Sprintf("%s/asset/variants/%s/%s/%s", api, url.PathEscape(v.Kind), url.PathEscape(v.Name), url.PathEscape(v.Variant))
}
func (v Variant) FileURL(api string) string {
	return fmt.Sprintf("%s/asset/variants/%s/%s/%s/file", api, url.PathEscape(v.Kind), url.PathEscape(v.Name), url.PathEscape(v.Variant))
}
