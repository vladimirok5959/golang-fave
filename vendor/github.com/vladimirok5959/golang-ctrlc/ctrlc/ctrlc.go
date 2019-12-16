package ctrlc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const C_ICON_START = "🌟"
const C_ICON_WARN = "⚡️"
const C_ICON_HOT = "🔥"
const C_ICON_MAG = "✨"
const C_ICON_SC = "🌳"

type Iface interface {
	Shutdown(ctx context.Context) error
}

type CallbackFunc func(ctx context.Context, shutdown context.CancelFunc) *[]Iface

func App(t time.Duration, f CallbackFunc) {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	fmt.Printf(
		icon_start(UseColors())+"%s\n",
		cly(
			UseColors(),
			fmt.Sprintf(
				"Application started (%d sec)",
				t/time.Second,
			),
		),
	)

	sctx, shutdown := context.WithCancel(context.Background())
	ifaces := f(sctx, shutdown)

	select {
	case <-sctx.Done():
		fmt.Printf(
			"\r"+icon_warn(UseColors())+"%s\n",
			cly(
				UseColors(),
				fmt.Sprintf(
					"Shutting down (application) (%d sec)",
					t/time.Second,
				),
			),
		)
	case val := <-stop:
		switch val {
		case syscall.SIGINT:
			fmt.Printf(
				"\r"+icon_warn(UseColors())+"%s\n",
				cly(
					UseColors(),
					fmt.Sprintf(
						"Shutting down (interrupt) (%d sec)",
						t/time.Second,
					),
				),
			)
		case syscall.SIGTERM:
			fmt.Printf(
				icon_warn(UseColors())+"%s\n",
				cly(
					UseColors(),
					fmt.Sprintf(
						"Shutting down (terminate) (%d sec)",
						t/time.Second,
					),
				),
			)
		default:
			fmt.Printf(
				icon_warn(UseColors())+"%s\n",
				cly(
					UseColors(),
					fmt.Sprintf(
						"Shutting down (%d sec)",
						t/time.Second,
					),
				),
			)
		}
	}

	shutdown()

	errors := false
	ctx, cancel := context.WithTimeout(context.Background(), t)
	for _, iface := range *ifaces {
		if err := iface.Shutdown(ctx); err != nil {
			errors = true
			fmt.Printf(
				icon_hot(UseColors())+"%s\n",
				clr(
					UseColors(),
					fmt.Sprintf(
						"Shutdown error (%T): %s",
						iface,
						err.Error(),
					),
				),
			)
		}
	}
	cancel()

	if errors {
		fmt.Printf(
			icon_mag(UseColors())+"%s\n",
			cly(
				UseColors(),
				fmt.Sprintf(
					"Application exited with errors (%d sec)",
					t/time.Second,
				),
			),
		)
		os.Exit(1)
	} else {
		fmt.Printf(
			icon_sc(UseColors())+"%s\n",
			clg(
				UseColors(),
				fmt.Sprintf(
					"Application exited successfully (%d sec)",
					t/time.Second,
				),
			),
		)
	}
}

func UseColors() bool {
	useColors := strings.Contains(
		fmt.Sprintf("%s", os.Args),
		"--color=always",
	)
	if !useColors {
		useColors = strings.Contains(
			fmt.Sprintf("%s", os.Args),
			"-color=always",
		)
	}
	if !useColors {
		useColors = strings.Contains(
			fmt.Sprintf("%s", os.Args),
			"color=always",
		)
	}
	if !useColors {
		useColors = strings.Contains(
			fmt.Sprintf("%s", os.Args),
			"--color always",
		)
	}
	if !useColors {
		useColors = strings.Contains(
			fmt.Sprintf("%s", os.Args),
			"-color always",
		)
	}
	if !useColors {
		useColors = strings.Contains(
			fmt.Sprintf("%s", os.Args),
			"color always",
		)
	}
	return !IS_WIN_PLATFORM && useColors
}
