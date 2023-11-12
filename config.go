package filesystem

import (
	"errors"
	"sync"
)

type Config interface {
	Get(key any, defaultValue ...any) (any, error)
	Extend(config Config) Config
	WithDefault(key any, value any) Config
	WithDefaults(config map[string]any) Config
	ToMap() map[string]any
	WithSetting(key any, value any) Config
	WithoutSettings(keys ...any) Config
}

type BaseConfig struct {
	config *sync.Map
}

func NewConfig(config map[string]any) *BaseConfig {
	m := &sync.Map{}

	for k, v := range config {
		m.Store(k, v)
	}

	return &BaseConfig{
		config: m,
	}
}

func (f *BaseConfig) Get(key any, defaultValue ...any) (any, error) {
	val, ok := f.config.Load(key)
	if ok {
		return val, nil
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return nil, errors.New("key not found")
}

func (f *BaseConfig) Extend(config Config) Config {
	m := config.ToMap()

	for k, v := range m {
		f.config.Store(k, v)
	}

	return f
}

func (f *BaseConfig) WithDefault(key any, value any) Config {
	f.config.Store(key, value)

	return f
}

func (f *BaseConfig) WithDefaults(config map[string]any) Config {
	for k, v := range config {
		f.config.Store(k, v)
	}

	return f
}

func (f *BaseConfig) ToMap() map[string]any {
	m := map[string]any{}

	f.config.Range(func(key any, value any) bool {
		m[key.(string)] = value

		return true
	})

	return m
}

func (f *BaseConfig) WithSetting(key any, value any) Config {
	f.config.Store(key, value)

	return f
}

func (f *BaseConfig) WithoutSettings(keys ...any) Config {
	for _, key := range keys {
		f.config.Delete(key)
	}

	return f
}
