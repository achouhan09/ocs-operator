// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	configv1 "github.com/openshift/api/config/v1"
)

// ReleaseApplyConfiguration represents a declarative configuration of the Release type for use
// with apply.
type ReleaseApplyConfiguration struct {
	Architecture *configv1.ClusterVersionArchitecture `json:"architecture,omitempty"`
	Version      *string                              `json:"version,omitempty"`
	Image        *string                              `json:"image,omitempty"`
	URL          *configv1.URL                        `json:"url,omitempty"`
	Channels     []string                             `json:"channels,omitempty"`
}

// ReleaseApplyConfiguration constructs a declarative configuration of the Release type for use with
// apply.
func Release() *ReleaseApplyConfiguration {
	return &ReleaseApplyConfiguration{}
}

// WithArchitecture sets the Architecture field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Architecture field is set to the value of the last call.
func (b *ReleaseApplyConfiguration) WithArchitecture(value configv1.ClusterVersionArchitecture) *ReleaseApplyConfiguration {
	b.Architecture = &value
	return b
}

// WithVersion sets the Version field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Version field is set to the value of the last call.
func (b *ReleaseApplyConfiguration) WithVersion(value string) *ReleaseApplyConfiguration {
	b.Version = &value
	return b
}

// WithImage sets the Image field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Image field is set to the value of the last call.
func (b *ReleaseApplyConfiguration) WithImage(value string) *ReleaseApplyConfiguration {
	b.Image = &value
	return b
}

// WithURL sets the URL field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the URL field is set to the value of the last call.
func (b *ReleaseApplyConfiguration) WithURL(value configv1.URL) *ReleaseApplyConfiguration {
	b.URL = &value
	return b
}

// WithChannels adds the given value to the Channels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Channels field.
func (b *ReleaseApplyConfiguration) WithChannels(values ...string) *ReleaseApplyConfiguration {
	for i := range values {
		b.Channels = append(b.Channels, values[i])
	}
	return b
}
