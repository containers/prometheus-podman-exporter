// Code generated by go generate; DO NOT EDIT.
package manifests

import (
	"net/url"

	"github.com/containers/podman/v5/pkg/bindings/internal/util"
)

// Changed returns true if named field has been set
func (o *AddArtifactOptions) Changed(fieldName string) bool {
	return util.Changed(o, fieldName)
}

// ToParams formats struct fields to be passed to API service
func (o *AddArtifactOptions) ToParams() (url.Values, error) {
	return util.ToParams(o)
}

// WithAnnotation set field Annotation to given value
func (o *AddArtifactOptions) WithAnnotation(value map[string]string) *AddArtifactOptions {
	o.Annotation = value
	return o
}

// GetAnnotation returns value of field Annotation
func (o *AddArtifactOptions) GetAnnotation() map[string]string {
	if o.Annotation == nil {
		var z map[string]string
		return z
	}
	return o.Annotation
}

// WithArch set field Arch to given value
func (o *AddArtifactOptions) WithArch(value string) *AddArtifactOptions {
	o.Arch = &value
	return o
}

// GetArch returns value of field Arch
func (o *AddArtifactOptions) GetArch() string {
	if o.Arch == nil {
		var z string
		return z
	}
	return *o.Arch
}

// WithFeatures set field Features to given value
func (o *AddArtifactOptions) WithFeatures(value []string) *AddArtifactOptions {
	o.Features = value
	return o
}

// GetFeatures returns value of field Features
func (o *AddArtifactOptions) GetFeatures() []string {
	if o.Features == nil {
		var z []string
		return z
	}
	return o.Features
}

// WithOS set field OS to given value
func (o *AddArtifactOptions) WithOS(value string) *AddArtifactOptions {
	o.OS = &value
	return o
}

// GetOS returns value of field OS
func (o *AddArtifactOptions) GetOS() string {
	if o.OS == nil {
		var z string
		return z
	}
	return *o.OS
}

// WithOSVersion set field OSVersion to given value
func (o *AddArtifactOptions) WithOSVersion(value string) *AddArtifactOptions {
	o.OSVersion = &value
	return o
}

// GetOSVersion returns value of field OSVersion
func (o *AddArtifactOptions) GetOSVersion() string {
	if o.OSVersion == nil {
		var z string
		return z
	}
	return *o.OSVersion
}

// WithOSFeatures set field OSFeatures to given value
func (o *AddArtifactOptions) WithOSFeatures(value []string) *AddArtifactOptions {
	o.OSFeatures = value
	return o
}

// GetOSFeatures returns value of field OSFeatures
func (o *AddArtifactOptions) GetOSFeatures() []string {
	if o.OSFeatures == nil {
		var z []string
		return z
	}
	return o.OSFeatures
}

// WithVariant set field Variant to given value
func (o *AddArtifactOptions) WithVariant(value string) *AddArtifactOptions {
	o.Variant = &value
	return o
}

// GetVariant returns value of field Variant
func (o *AddArtifactOptions) GetVariant() string {
	if o.Variant == nil {
		var z string
		return z
	}
	return *o.Variant
}

// WithType set field Type to given value
func (o *AddArtifactOptions) WithType(value *string) *AddArtifactOptions {
	o.Type = &value
	return o
}

// GetType returns value of field Type
func (o *AddArtifactOptions) GetType() *string {
	if o.Type == nil {
		var z *string
		return z
	}
	return *o.Type
}

// WithConfigType set field ConfigType to given value
func (o *AddArtifactOptions) WithConfigType(value string) *AddArtifactOptions {
	o.ConfigType = &value
	return o
}

// GetConfigType returns value of field ConfigType
func (o *AddArtifactOptions) GetConfigType() string {
	if o.ConfigType == nil {
		var z string
		return z
	}
	return *o.ConfigType
}

// WithConfig set field Config to given value
func (o *AddArtifactOptions) WithConfig(value string) *AddArtifactOptions {
	o.Config = &value
	return o
}

// GetConfig returns value of field Config
func (o *AddArtifactOptions) GetConfig() string {
	if o.Config == nil {
		var z string
		return z
	}
	return *o.Config
}

// WithLayerType set field LayerType to given value
func (o *AddArtifactOptions) WithLayerType(value string) *AddArtifactOptions {
	o.LayerType = &value
	return o
}

// GetLayerType returns value of field LayerType
func (o *AddArtifactOptions) GetLayerType() string {
	if o.LayerType == nil {
		var z string
		return z
	}
	return *o.LayerType
}

// WithExcludeTitles set field ExcludeTitles to given value
func (o *AddArtifactOptions) WithExcludeTitles(value bool) *AddArtifactOptions {
	o.ExcludeTitles = &value
	return o
}

// GetExcludeTitles returns value of field ExcludeTitles
func (o *AddArtifactOptions) GetExcludeTitles() bool {
	if o.ExcludeTitles == nil {
		var z bool
		return z
	}
	return *o.ExcludeTitles
}

// WithSubject set field Subject to given value
func (o *AddArtifactOptions) WithSubject(value string) *AddArtifactOptions {
	o.Subject = &value
	return o
}

// GetSubject returns value of field Subject
func (o *AddArtifactOptions) GetSubject() string {
	if o.Subject == nil {
		var z string
		return z
	}
	return *o.Subject
}

// WithAnnotations set field Annotations to given value
func (o *AddArtifactOptions) WithAnnotations(value map[string]string) *AddArtifactOptions {
	o.Annotations = value
	return o
}

// GetAnnotations returns value of field Annotations
func (o *AddArtifactOptions) GetAnnotations() map[string]string {
	if o.Annotations == nil {
		var z map[string]string
		return z
	}
	return o.Annotations
}

// WithFiles set field Files to given value
func (o *AddArtifactOptions) WithFiles(value []string) *AddArtifactOptions {
	o.Files = value
	return o
}

// GetFiles returns value of field Files
func (o *AddArtifactOptions) GetFiles() []string {
	if o.Files == nil {
		var z []string
		return z
	}
	return o.Files
}
