package app

import (
	"context"
	"fmt"
	"sort"

	"github.com/thekhanj/tci/internal/config"
	"github.com/thekhanj/tci/panel"
)

func GetAccountPanel(
	ctx context.Context,
	api *panel.Api,
	username string,
	acc config.Account,
) (panel.Panel, error) {
	session, err := NewAccountSession(ctx, username, acc)
	if err != nil {
		return panel.Panel{}, err
	}

	return api.Get(ctx, session)
}

func GetPanels(
	ctx context.Context,
	cfg config.Config,
) ([]panel.Panel, error) {
	api, err := panel.NewApi()
	if err != nil {
		return nil, fmt.Errorf("create panel API: %w", err)
	}
	defer api.Close()

	usernames := make([]string, 0, len(cfg.Accounts))
	for username := range cfg.Accounts {
		usernames = append(usernames, username)
	}
	sort.Strings(usernames)

	panels := make([]panel.Panel, 0, len(usernames))

	for _, username := range usernames {
		p, err := GetAccountPanel(
			ctx,
			api,
			username,
			cfg.Accounts[username],
		)
		if err != nil {
			return nil, fmt.Errorf("account %s: %w", username, err)
		}

		panels = append(panels, p)
	}

	return panels, nil
}
