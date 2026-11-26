package app

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/thekhanj/tci"
	"github.com/thekhanj/tci/internal/config"
)

func NewAccountSession(
	ctx context.Context,
	username string,
	acc config.Account,
) (*tci.Session, error) {
	serializer, err := tci.NewClientSerializer()
	if err != nil {
		return nil, fmt.Errorf("create client serializer: %w", err)
	}

	store, err := tci.NewStore[tci.Client](
		username+".json",
		serializer,
	)
	if err != nil {
		return nil, fmt.Errorf("create session store: %w", err)
	}

	client, err := tci.NewClient()
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	saved, err := store.Load(ctx)
	if err == nil {
		client = saved
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("load session: %w", err)
	}

	session, err := tci.NewSession(client)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	authenticated, err := session.Authenticated(ctx)
	if err != nil {
		return nil, fmt.Errorf("check authentication: %w", err)
	}

	if authenticated {
		return session, nil
	}

	if err := session.Login(ctx, username, acc.Password); err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	if err := store.Save(ctx, client); err != nil {
		return nil, fmt.Errorf("save session: %w", err)
	}

	return session, nil
}
