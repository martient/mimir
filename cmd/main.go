package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/taultek/mimir/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "mimir",
	Short: "Mimir - AI agent orchestrator",
	Long:  `Mimir is a personal AI assistant and orchestrator for managing multiple opencode instances across different projects, contexts, and use cases.`,
	Run:   runServer,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Mimir gateway",
	Long:  `Start the Mimir gateway server for handling agent orchestration, webhooks, and API requests.`,
	Run:   runServer,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show gateway status",
	Long:  `Show the current status of the Mimir gateway and its components.`,
	Run:   runStatus,
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long:  `Manage registered projects in Mimir.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered projects",
	Long:  `List all projects registered in Mimir.`,
	Run:   runProjectsList,
}

var projectsAddCmd = &cobra.Command{
	Use:   "add <path>",
	Short: "Register a new project",
	Long:  `Register a new project in Mimir for agent management.`,
	Args:  cobra.MinimumNArgs(1),
	Run:   runProjectsAdd,
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message to opencode",
	Long:  `Send a message to a specific opencode instance.`,
	Run:   runSend,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(projectsListCmd)
	projectsCmd.AddCommand(projectsAddCmd)
	rootCmd.AddCommand(sendCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	fmt.Println("Starting Mimir Gateway...")
	fmt.Println("Server would start on HTTP port 8080 and WS port 8081")
	fmt.Println("Server implementation will be completed in later phases")
}

func runStatus(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Mimir Gateway Status:")
	fmt.Printf("  Config: %s\n", cfg.Database.Path)
	fmt.Println("  Server: Not running (use 'mimir serve')")
	fmt.Printf("  Projects: %d registered\n", len(cfg.Projects))
	fmt.Println("  Active Sessions: N/A (database access not yet implemented in CLI)")
}

func runProjectsList(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Projects) == 0 {
		fmt.Println("No projects registered.")
		fmt.Println("Use 'mimir projects add <path>' to register a project.")
		return
	}

	fmt.Println("Registered Projects:")
	for i, project := range cfg.Projects {
		fmt.Printf("  %d. %s\n", i+1, project.Name)
		fmt.Printf("     Path: %s\n", project.Path)
	}
}

func runProjectsAdd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: project path is required")
		return
	}

	projectPath := args[0]
	fmt.Printf("Registering project: %s\n", projectPath)
	fmt.Println("Project will be added to Mimir configuration")
	fmt.Println("Project registration will be implemented in Phase 3 (with database integration)")
}

func runSend(cmd *cobra.Command, args []string) {
	project, _ := cmd.Flags().GetString("project")
	message, _ := cmd.Flags().GetString("message")

	if project == "" {
		fmt.Println("Error: --project flag is required")
		return
	}

	if message == "" {
		fmt.Println("Error: --message flag is required")
		return
	}

	fmt.Printf("Sending message to project '%s': %s\n", project, message)
	fmt.Println("Message will be sent to opencode instance")
	fmt.Println("Message sending will be implemented in Phase 3 (with opencode SDK integration)")
}
