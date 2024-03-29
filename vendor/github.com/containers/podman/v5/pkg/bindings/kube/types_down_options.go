// Code generated by go generate; DO NOT EDIT.
package kube

import (
	"net/url"

	"github.com/containers/podman/v5/pkg/bindings/internal/util"
)

// Changed returns true if named field has been set
func (o *DownOptions) Changed(fieldName string) bool {
	return util.Changed(o, fieldName)
}

// ToParams formats struct fields to be passed to API service
func (o *DownOptions) ToParams() (url.Values, error) {
	return util.ToParams(o)
}

// WithForce set field Force to given value
func (o *DownOptions) WithForce(value bool) *DownOptions {
	o.Force = &value
	return o
}

// GetForce returns value of field Force
func (o *DownOptions) GetForce() bool {
	if o.Force == nil {
		var z bool
		return z
	}
	return *o.Force
}
