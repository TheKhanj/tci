package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/thekhanj/tci"
	"github.com/thekhanj/tci/internal/app"
	"github.com/thekhanj/tci/internal/config"
)

func Status(
	ctx context.Context,
	configPath string,
	output OutputOptions,
) error {
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		return err
	}

	panels, err := app.GetPanels(ctx, cfg)
	if err != nil {
		return err
	}

	switch tci.OutputFormat(output.Format) {
	case tci.OutputJson:
		return tci.TableWriteJson(os.Stdout, panels)

	case tci.OutputTsv:
		var columns tci.Columns
		if len(panels) > 0 {
			columns = panels[0].Columns()
		}

		headers := output.Headers(columns.Headers)
		return tci.TableWriteTsv(os.Stdout, headers, panels)

	default:
		return fmt.Errorf("unsupported output format: %q", output.Format)
	}
}
