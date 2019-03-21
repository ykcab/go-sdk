package web

import (
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

var (
	_ configutil.ConfigResolver = (*ViewCacheConfig)(nil)
)

// ViewCacheConfig is a config for the view cache.
type ViewCacheConfig struct {
	// LiveReload indicates if we should store compiled views in memory for re-use (default), or read them from disk each load.
	LiveReload bool `json:"liveReload,omitempty" yaml:"liveReload,omitempty" env:"WEB_VIEW_LIVE_RELOAD"`
	// Paths are a list of view paths to include in the templates list.
	Paths []string `json:"paths,omitempty" yaml:"paths,omitempty" env:"WEB_VIEW_CACHE_PATHS,csv"`
	// BufferPoolSize is the size of the re-usable buffer pool for rendering views.
	BufferPoolSize int `json:"bufferPoolSize,omitempty" yaml:"bufferPoolSize,omitempty"`

	// InternalErrorTemplateName is the template name to use for the view result provider `InternalError` result.
	InternalErrorTemplateName string `json:"internalErrorTemplateName,omitempty" yaml:"internalErrorTemplateName,omitempty"`
	// BadRequestTemplateName is the template name to use for the view result provider `BadRequest` result.
	BadRequestTemplateName string `json:"badRequestTemplateName,omitempty" yaml:"badRequestTemplateName,omitempty"`
	// NotFoundTemplateName is the template name to use for the view result provider `NotFound` result.
	NotFoundTemplateName string `json:"notFoundTemplateName,omitempty" yaml:"notFoundTemplateName,omitempty"`
	// NotAuthorizedTemplateName is the template name to use for the view result provider `NotAuthorized` result.
	NotAuthorizedTemplateName string `json:"notAuthorizedTemplateName,omitempty" yaml:"notAuthorizedTemplateName,omitempty"`
	// StatusTemplateName is the template name to use for the view result provider status result.
	StatusTemplateName string `json:"statusTemplateName,omitempty" yaml:"statusTemplateName,omitempty"`
}

// Resolve adds extra resolution steps when we setup the config.
func (vcc *ViewCacheConfig) Resolve() error {
	return env.Env().ReadInto(vcc)
}

// BufferPoolSizeOrDefault gets the buffer pool size or a default.
func (vcc ViewCacheConfig) BufferPoolSizeOrDefault() int {
	return configutil.CoalesceInt(vcc.BufferPoolSize, DefaultViewBufferPoolSize)
}

// InternalErrorTemplateNameOrDefault returns the internal error template name for the app.
func (vcc ViewCacheConfig) InternalErrorTemplateNameOrDefault() string {
	return configutil.CoalesceString(vcc.InternalErrorTemplateName, DefaultTemplateNameInternalError)
}

// BadRequestTemplateNameOrDefault returns the bad request template name for the app.
func (vcc ViewCacheConfig) BadRequestTemplateNameOrDefault() string {
	return configutil.CoalesceString(vcc.BadRequestTemplateName, DefaultTemplateNameBadRequest)
}

// NotFoundTemplateNameOrDefault returns the not found template name for the app.
func (vcc ViewCacheConfig) NotFoundTemplateNameOrDefault() string {
	return configutil.CoalesceString(vcc.NotFoundTemplateName, DefaultTemplateNameNotFound)
}

// NotAuthorizedTemplateNameOrDefault returns the not authorized template name for the app.
func (vcc ViewCacheConfig) NotAuthorizedTemplateNameOrDefault() string {
	return configutil.CoalesceString(vcc.NotAuthorizedTemplateName, DefaultTemplateNameNotAuthorized)
}

// StatusTemplateNameOrDefault returns the not authorized template name for the app.
func (vcc ViewCacheConfig) StatusTemplateNameOrDefault() string {
	return configutil.CoalesceString(vcc.StatusTemplateName, DefaultTemplateNameStatus)
}
