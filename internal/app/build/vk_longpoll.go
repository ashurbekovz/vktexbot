package build

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

	"os"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/events"
	"github.com/SevereCloud/vksdk/v3/longpoll-bot"
	"github.com/ashurbekovz/vktexbot/internal/app/parsers"
	"github.com/ashurbekovz/vktexbot/internal/app/usecases"
	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img"
	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"os/signal"
	"syscall"
)

type VkApp struct {
	logger *slog.Logger

	configPath string
	secretPath string
}

type Secret struct {
	GroupID string `yaml:"vk_group_id"`
	Token   string `yaml:"vk_token"`
}

type Config struct {
	Packages string `yaml:"packages"`
}

func NewVkApp(configPath string, secretPath string) *VkApp {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	return &VkApp{
		configPath: configPath,
		secretPath: secretPath,
		logger:     logger,
	}
}

func (a *VkApp) Run() {
	a.logger.Info("starting application", "config", a.configPath, "secret", a.secretPath)

	config, err := load[Config](a.configPath)
	if err != nil {
		a.logger.Error("failed to load config", "error", err)
		return
	}

	secret, err := load[Secret](a.secretPath)
	if err != nil {
		a.logger.Error("failed to load secrets", "error", err)
		return
	}

	vk := api.NewVK(secret.Token)

	group, err := vk.GroupsGetByID(map[string]any{"group_id": secret.GroupID})
	if err != nil {
		a.logger.Error("failed to get group info", "error", err, "group_id", secret.GroupID)
		return
	} else if len(group.Groups) != 1 {
		a.logger.Error("unexpected number of groups", "count", len(group.Groups), "expected", 1)
		return
	}

	lp, err := longpoll.NewLongPoll(vk, group.Groups[0].ID)
	if err != nil {
		a.logger.Error("failed to create Long Poll", "error", err, "group_id", group.Groups[0].ID)
		return
	}
	a.logger.Info("long poll initialized")

	l2i := latex2img.NewLatexToImgConverter("/tmp/", false, decimal.NewFromInt(400))
	t2i := template2img.NewLatexTemplateToImgConverter(&l2i, config.Packages)
	parser := parsers.NewVk(group.Groups[0].ID, group.Groups[0].ScreenName)
	vkUC := usecases.NewVkUsecase(vk, &t2i, parser)

	lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		opt := usecases.VkOpt{
			Message:      obj.Message.Text,
			PeerID:       obj.Message.PeerID,
			Payload:      obj.Message.Payload,
			IsNewMessage: true,
		}

		logger := a.logger.With("request_id", uuid.New().String())
		_, err := vkUC.Execute(ctx, logger, opt)
		if err != nil {
			logger.Error("failed to execute message", "error", err, "peer_id", opt.PeerID)
			return
		}
	})

	lp.MessageEdit(func(ctx context.Context, obj events.MessageEditObject) {
		opt := usecases.VkOpt{
			Message:      obj.Text,
			PeerID:       obj.PeerID,
			Payload:      obj.Payload,
			IsNewMessage: false,
		}

		logger := a.logger.With("request_id", uuid.New().String())
		_, err := vkUC.Execute(ctx, logger, opt)
		if err != nil {
			logger.Error("failed to process message edit", "error", err, "peer_id", opt.PeerID)
			return
		}
	})

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := <-shutdown
		a.logger.Info("shutdown signal received", "signal", sig)
		cancel()
	}()

	a.logger.Info("starting VK Long Poll server")
	if err := lp.RunWithContext(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			a.logger.Info("server stopped by signal")
		} else {
			a.logger.Error("server crashed", "error", err)
		}
	}

	a.logger.Info("VK Long Poll Server stopped")
}

func load[config any](path string) (*config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error opening config file: %w", err)
	}
	defer file.Close()

	var cfg config
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("Error decoding config file: %w", err)
	}

	return &cfg, nil
}
