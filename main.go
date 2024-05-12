package main

import (
	"github.com/canghai118/posts-list/cmd/buildindex"
	"github.com/canghai118/posts-list/cmd/gen"
	"github.com/canghai118/posts-list/cmd/serve"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(buildindex.BuildIndexCmd)
	rootCmd.AddCommand(serve.ServeCmd)
	rootCmd.AddCommand(gen.GenCmd)

	rootCmd.Execute()
}
