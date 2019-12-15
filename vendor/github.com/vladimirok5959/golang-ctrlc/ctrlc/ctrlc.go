package ctrlc

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
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
	var ParamColor string
	flag.StringVar(&ParamColor, "color", "auto", "color output (auto/always/never)")
	flag.Parse()

	useColors := !IS_WIN_PLATFORM && ParamColor == "always"

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	fmt.Printf(
		icon_start(useColors)+"%s\n",
		cly(
			useColors,
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
			"\r"+icon_warn(useColors)+"%s\n",
			cly(
				useColors,
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
				"\r"+icon_warn(useColors)+"%s\n",
				cly(
					useColors,
					fmt.Sprintf(
						"Shutting down (interrupt) (%d sec)",
						t/time.Second,
					),
				),
			)
		case syscall.SIGTERM:
			fmt.Printf(
				icon_warn(useColors)+"%s\n",
				cly(
					useColors,
					fmt.Sprintf(
						"Shutting down (terminate) (%d sec)",
						t/time.Second,
					),
				),
			)
		default:
			fmt.Printf(
				icon_warn(useColors)+"%s\n",
				cly(
					useColors,
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
				icon_hot(useColors)+"%s\n",
				clr(
					useColors,
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
			icon_mag(useColors)+"%s\n",
			cly(
				useColors,
				fmt.Sprintf(
					"Application exited with errors (%d sec)",
					t/time.Second,
				),
			),
		)
		os.Exit(1)
	} else {
		fmt.Printf(
			icon_sc(useColors)+"%s\n",
			clg(
				useColors,
				fmt.Sprintf(
					"Application exited successfully (%d sec)",
					t/time.Second,
				),
			),
		)
	}
}
