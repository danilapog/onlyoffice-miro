package document

import (
	"context"
	"encoding/json"
)

type BuilderService interface {
	Build(ctx context.Context, callbackUrl string,
		configurer DocumentConfigurer, opts ...BuilderOption) (*Config, error)
}

type builderService struct {
	keyGenerator       KeyGenerator
	signatureGenerator SignatureGenerator
	formatManager      FormatManager
}

func NewBuilderService(
	keyGenerator KeyGenerator,
	signatureGenerator SignatureGenerator,
	formatManager FormatManager,
) BuilderService {
	return &builderService{
		keyGenerator:       keyGenerator,
		signatureGenerator: signatureGenerator,
		formatManager:      formatManager,
	}
}

func (s *builderService) Build(
	ctx context.Context,
	callbackUrl string,
	configurer DocumentConfigurer,
	opts ...BuilderOption,
) (*Config, error) {
	options := BuilderOptions{mode: Desktop}
	for _, option := range opts {
		option(&options)
	}

	title := configurer.Title()
	ext := s.formatManager.GetFileExt(title)
	key, err := s.keyGenerator.Generate(ctx, configurer)
	if err != nil {
		return nil, err
	}

	format, exists := s.formatManager.GetFormatByName(ext)
	if !exists {
		return nil, ErrUnsupportedFormat
	}

	config := &Config{
		Document: Document{
			Key:      key,
			Title:    title,
			URL:      configurer.URL(),
			FileType: ext,
			Permissions: Permissions{
				Edit: format.IsEditable(),
			},
		},
		Editor: Editor{
			CallbackURL: callbackUrl,
		},
		DocumentType: format.Type,
		Type:         string(options.mode),
	}

	if options.userConfigurer != nil {
		uconfigurer := options.userConfigurer
		config.Editor.User = User{
			ID:   uconfigurer.ID(),
			Name: uconfigurer.Name(),
		}

		config.Editor.Lang = uconfigurer.Language()
	}

	if len(options.key) > 0 {
		if err := s.signConfig(config, options.key); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func (s *builderService) signConfig(config *Config, secret []byte) error {
	buf, err := json.Marshal(config)
	if err != nil {
		return err
	}

	token, err := s.signatureGenerator.Sign(secret, buf)
	if err != nil {
		return err
	}

	config.Token = token
	return nil
}
