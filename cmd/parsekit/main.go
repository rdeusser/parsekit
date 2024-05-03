package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rdeusser/parsekit/lexer"

	"github.com/k0kubun/pp/v3"

	"github.com/rdeusser/parsekit/internal/logging"
	"github.com/rdeusser/parsekit/lang/golang"
	"github.com/rdeusser/parsekit/version"
)

type rootOptions struct {
	Debug    bool
	Lang     string
	Filename string
}

func (o *rootOptions) Init() {
	o.Debug = false
	o.Lang = ""
	o.Filename = ""
}

func main() {
	options := &rootOptions{}
	options.Init()

	logger := logging.New(os.Stderr, logging.Info)

	cmd := &cobra.Command{
		Use:     "parsekit",
		Short:   "Tools to help with lexing and parsing languages",
		Version: version.GetHumanVersion(),
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		Args: cobra.MaximumNArgs(1),
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if options.Debug {
				logger.SetLevel(logging.Debug)
			}
		},
		RunE: func(_ *cobra.Command, args []string) error {
			options.Filename = args[0]
			return run(logger, *options, args)
		},
		SilenceUsage:  true,
		SilenceErrors: true, // we log and return our own errors and if this is false the errors are printed twice
	}

	cmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	cmd.PersistentFlags().BoolVar(&options.Debug, "debug", options.Debug, "Run in debug mode")
	cmd.Flags().StringVarP(&options.Lang, "lang", "l", options.Lang, "Language to lex/parse")

	if err := cmd.Execute(); err != nil {
		logger.Error(err.Error())
	}
}

func run(logger logging.Logger, options rootOptions, args []string) error {
	l := lexer.New(lexer.DefaultConfig, lexer.WithLogger(logger))

	switch strings.ToLower(options.Lang) {
	case "go", "golang":
		l = golang.NewLexer(lexer.WithLogger(logger))
	}

	prompt := "> "

	if len(args) == 0 {
		fmt.Printf("Welcome to parsekit! To exit, either type \"exit\" or press enter twice.\n\n")

		count := 0
		scanner := bufio.NewScanner(os.Stdin)

		for {
			fmt.Print(prompt)
			if !scanner.Scan() {
				break // Exit loop if no more input
			}

			input := scanner.Text()
			if input == "" {
				count++
				if count == 2 {
					fmt.Println("Detected two consecutive ENTER presses. Exiting.")
					break
				}
			} else if input == "exit" {
				break
			} else {
				count = 0 // reset the counter if input is not empty

				tokens, err := l.Lex(input)
				if err != nil {
					return err
				}

				pp.Println(tokens)
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	} else {
		input, err := os.ReadFile(options.Filename)
		if err != nil {
			return err
		}

		tokens, err := l.Lex(string(input))
		if err != nil {
			return err
		}

		pp.Println(tokens)
	}

	return nil
}
