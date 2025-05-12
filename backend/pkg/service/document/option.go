package document

type BuilderOptions struct {
	key            []byte
	userConfigurer UserConfigurer
	mode           EditorMode
}

type BuilderOption func(*BuilderOptions)

func WithKey(val []byte) BuilderOption {
	return func(o *BuilderOptions) {
		if len(val) > 0 {
			o.key = val
		}
	}
}

func WithUserConfigurer(val UserConfigurer) BuilderOption {
	return func(o *BuilderOptions) {
		if val != nil {
			o.userConfigurer = val
		}
	}
}

func WithEditorMode(val EditorMode) BuilderOption {
	return func(o *BuilderOptions) {
		o.mode = val
	}
}
