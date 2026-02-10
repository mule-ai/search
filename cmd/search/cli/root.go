// Package cli implements the command-line interface for the search tool.
//
// It uses the Cobra framework for command parsing and provides the main
// search command with all its flags and subcommands.
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/mule-ai/search/internal/browser"
	"github.com/mule-ai/search/internal/cache"
	"github.com/mule-ai/search/internal/config"
	"github.com/mule-ai/search/internal/formatter"
	searxnglib "github.com/mule-ai/search/internal/searxng"
	"github.com/mule-ai/search/internal/ui"
	"github.com/mule-ai/search/internal/validation"
	"github.com/mule-ai/search/pkg/version"
)

// RootCommand wraps a Cobra Command with version information.
type RootCommand struct {
	*cobra.Command
	Version string
}

// ConfigFlags holds CLI flag values.
//
// These represent the command-line arguments that override config settings.
type ConfigFlags struct {
	Instance     string
	Results      int
	Format       string
	Category     string
	Timeout      int
	Language     string
	SafeSearch   int
	ConfigPath   string
	Verbose      bool
	Page         int
	TimeRange    string
	Open         bool
	OpenAll      bool
	NoColor      bool
	APIKey       string
	// Cache flags
	CacheEnabled *bool
	NoCache      bool
	CacheSize    int
	CacheTTL     int
	ClearCache   bool
	CacheStats   bool
}

func NewRootCommand() *RootCommand {
	var cfgFlags ConfigFlags

	rc := &RootCommand{
		Version: version.Version,
	}
	cmd := &cobra.Command{
		Use:           "search",
		Short:         "Search CLI - A command-line search tool using SearXNG",
		Long: `Search is a powerful command-line search tool that queries
SearXNG instances and returns formatted results.

Usage:
  search [flags] <query>

Examples:
  search "golang tutorials"
  search -n 20 "machine learning"
  search -f json "rust programming" | jq '.results[] | .title'`,
		PersistentPreRunE: persistentPreRun(&cfgFlags),
		RunE:              run(&cfgFlags),
		Args:            cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}

	addGlobalFlags(cmd.Flags(), &cfgFlags)
	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(newCategoriesCommand())
	AddCompletionCommand(cmd)

	// Set version template for --version flag
	versionOutput := fmt.Sprintf("search version %s", version.Version)
	if version.GitCommit != "unknown" {
		versionOutput += fmt.Sprintf("\nGit commit: %s", version.GitCommit)
	}
	if version.BuildDate != "unknown" {
		versionOutput += fmt.Sprintf("\nBuilt: %s", version.BuildDate)
	}
	if version.GoVersion != "unknown" {
		versionOutput += fmt.Sprintf("\nGo version: %s", version.GoVersion)
	}
	cmd.SetVersionTemplate(versionOutput + "\n")
	cmd.Version = version.Version

	rc.Command = cmd
	return rc
}

func addGlobalFlags(fs *pflag.FlagSet, cfg *ConfigFlags) {
	fs.StringVarP(&cfg.Instance, "instance", "i",
		"https://search.butler.ooo", "SearXNG instance URL")
	fs.IntVarP(&cfg.Results, "results", "n",
		10, "Number of results to return")
	fs.StringVarP(&cfg.Format, "format", "f",
		"text", "Output format: json, markdown, text")
	fs.StringVarP(&cfg.Category, "category", "c",
		"general", "Search category")
	fs.IntVarP(&cfg.Timeout, "timeout", "t",
		30, "Request timeout in seconds")
	fs.StringVarP(&cfg.Language, "language", "l",
		"en", "Language code")
	fs.IntVarP(&cfg.SafeSearch, "safe", "s",
		1, "Safe search level (0, 1, 2)")
	fs.StringVar(&cfg.ConfigPath, "config", "",
		"Custom config file path")
	fs.BoolVarP(&cfg.Verbose, "verbose", "v",
		false, "Enable verbose output")
	fs.IntVar(&cfg.Page, "page", 1, "Page number for pagination")
	fs.StringVar(&cfg.TimeRange, "time", "",
		"Time range filter: day, week, month, year")
	fs.BoolVar(&cfg.Open, "open", false,
		"Open first result in browser")
	fs.BoolVar(&cfg.OpenAll, "open-all", false,
		"Open all results in browser")
	fs.BoolVar(&cfg.NoColor, "no-color", false,
		"Disable colored output")
	fs.StringVar(&cfg.APIKey, "api-key", "",
		"API key for SearXNG authentication")
	// Cache flags
	cacheEnabled := true
	fs.BoolVar(&cacheEnabled, "cache", true,
		"Enable request caching")
	cfg.CacheEnabled = &cacheEnabled
	fs.BoolVar(&cfg.NoCache, "no-cache", false,
		"Disable caching for this request")
	fs.IntVar(&cfg.CacheSize, "cache-size", 100,
		"Maximum number of cache entries")
	fs.IntVar(&cfg.CacheTTL, "cache-ttl", 300,
		"Cache TTL in seconds (default: 300)")
	fs.BoolVar(&cfg.ClearCache, "clear-cache", false,
		"Clear the cache before searching")
	fs.BoolVar(&cfg.CacheStats, "cache-stats", false,
		"Show cache statistics")
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("search version %s\n", version.Version)
			if version.GitCommit != "unknown" {
				fmt.Printf("Git commit: %s\n", version.GitCommit)
			}
			if version.BuildDate != "unknown" {
				fmt.Printf("Built: %s\n", version.BuildDate)
			}
			if version.GoVersion != "unknown" {
				fmt.Printf("Go version: %s\n", version.GoVersion)
			}
			return nil
		},
	}
}

func newCategoriesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "categories",
		Short: "List available search categories",
		Long: `List all available search categories supported by SearXNG.

Categories allow you to narrow your search to specific types of content:
  - general: Web search across all engines
  - images: Image search results
  - videos: Video search results
  - news: News articles
  - map: Map results
  - music: Music search
  - it: Information technology
  - science: Scientific research
  - files: File and document search
  - social media: Social media content`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Available Search Categories:")
			fmt.Println()
			
			categories := searxnglib.GetCategoryNames()
			for _, name := range categories {
				cat, _ := searxnglib.GetCategory(name)
				fmt.Printf("  %-15s %s\n", cat.Name, cat.DisplayName)
				fmt.Printf("                 %s\n", cat.Description)
				fmt.Printf("                 Example: search -c %s \"%s\"\n", cat.Name, cat.ExampleQuery)
				fmt.Println()
			}
			
			return nil
		},
	}
}

func persistentPreRun(cfg *ConfigFlags) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Loading configuration...\n")
		}
		return nil
	}
}

func run(cfgFlags *ConfigFlags) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		query := args[0]

		// Sanitize input to remove potentially dangerous characters
		query = ui.SanitizeInput(query)

		// Validate inputs
		if err := validation.ValidateQuery(query); err != nil {
			return err
		}
		if err := validation.ValidateResultCount(cfgFlags.Results); err != nil {
			return err
		}
		if err := validation.ValidateTimeout(cfgFlags.Timeout); err != nil {
			return err
		}
		if err := validation.ValidateFormat(cfgFlags.Format); err != nil {
			return err
		}
		if err := validation.ValidateCategory(cfgFlags.Category); err != nil {
			return err
		}
		if err := validation.ValidateSafeSearch(cfgFlags.SafeSearch); err != nil {
			return err
		}
		if err := validation.ValidateLanguage(cfgFlags.Language); err != nil {
			return err
		}
		if err := validation.ValidatePageNumber(cfgFlags.Page); err != nil {
			return err
		}
		if err := validation.ValidateTimeRange(cfgFlags.TimeRange); err != nil {
			return err
		}

		if cfgFlags.Verbose {
			fmt.Fprintf(os.Stderr, "Query: %s\n", query)
		}

		// Load config
		cfgOverride := &config.CliConfig{
			ConfigPath:  cfgFlags.ConfigPath,
			Verbose:     cfgFlags.Verbose,
			Page:        cfgFlags.Page,
			TimeRange:   cfgFlags.TimeRange,
		}

		// Only override config with CLI flags if they were explicitly set
		if cmd.Flags().Changed("instance") {
			cfgOverride.Instance = cfgFlags.Instance
		}
		if cmd.Flags().Changed("results") {
			cfgOverride.Results = cfgFlags.Results
		}
		if cmd.Flags().Changed("format") {
			cfgOverride.Format = cfgFlags.Format
		}
		if cmd.Flags().Changed("category") {
			cfgOverride.Category = cfgFlags.Category
		}
		if cmd.Flags().Changed("timeout") {
			cfgOverride.Timeout = cfgFlags.Timeout
		}
		if cmd.Flags().Changed("language") {
			cfgOverride.Language = cfgFlags.Language
		}
		if cmd.Flags().Changed("safe") {
			cfgOverride.SafeSearch = cfgFlags.SafeSearch
		}
		if cmd.Flags().Changed("cache") {
			cfgOverride.CacheEnabled = cfgFlags.CacheEnabled
		}
		if cmd.Flags().Changed("no-cache") {
			cfgOverride.NoCache = cfgFlags.NoCache
		}
		if cmd.Flags().Changed("cache-size") && cfgFlags.CacheSize > 0 {
			cfgOverride.CacheSize = &cfgFlags.CacheSize
		}
		if cmd.Flags().Changed("cache-ttl") && cfgFlags.CacheTTL > 0 {
			cfgOverride.CacheTTL = &cfgFlags.CacheTTL
		}
		if cmd.Flags().Changed("api-key") {
			cfgOverride.APIKey = cfgFlags.APIKey
		}

		cfg, err := config.LoadConfig(cfgOverride)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Validate instance URL from final config
		if err := validation.ValidateInstanceURL(cfg.Instance); err != nil {
			return err
		}

		if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Using instance: %s\n", cfg.Instance)
		}

		// Create SearXNG client
		client := searxnglib.NewClient(cfg)

		// Wrap with caching if enabled
		var searchClient interface {
			SearchWithConfig(query string, results int, format string, category string, timeout int, language string, safeSearch int, page int, timeRange string) (*searxnglib.SearchResponse, error)
		}
		searchClient = client

		var cachedClient *cache.CachedClient
		if cfg.CacheEnabled {
			cachedClient = cache.NewCachedClient(
				client,
				cfg.CacheSize,
				time.Duration(cfg.CacheTTL)*time.Second,
			)

			// Handle cache clearing if requested
			if cfgFlags.ClearCache {
				cachedClient.ClearCache()
				if cfg.Verbose {
					fmt.Fprintf(os.Stderr, "Cache cleared\n")
				}
			}

			// Handle cache stats
			if cfgFlags.CacheStats {
				stats := cachedClient.GetStats()
				fmt.Fprintf(os.Stderr, "Cache stats: %d/%d entries\n", stats.Size, stats.MaxSize)
			}

			// Create a wrapper that implements the SearchWithConfig interface
			searchClient = &cachedSearchClient{cached: cachedClient}
		}

		// Create spinner for search operation
		spinner := ui.NewSearchSpinner(cfg.Verbose && !cfgFlags.NoColor)

		// Start spinner
		spinner.Start()
		startTime := time.Now()

		// Perform search
		results, err := searchClient.SearchWithConfig(
			query,
			cfg.Results,
			cfg.Format,
			cfg.Categories[0],
			cfg.Timeout,
			cfg.Language,
			cfg.SafeSearch,
			cfgFlags.Page,
			cfgFlags.TimeRange,
		)

		// Calculate search duration
		duration := time.Since(startTime)

		// Stop spinner with results
		if err != nil {
			spinner.StopWithError(err)
			return fmt.Errorf("search failed: %w", err)
		}

		spinner.Stop(len(results.Results), fmt.Sprintf("%.2fs", duration.Seconds()))

		if !cfg.Verbose {
			// If not verbose, spinner already showed the results count
		} else {
			fmt.Fprintf(os.Stderr, "Found %d results\n", len(results.Results))
		}

		// Format and output results - use category-aware formatter
		category := ""
		if len(cfg.Categories) > 0 {
			category = cfg.Categories[0]
		}
		outputFormatter, err := formatter.NewFormatterForCategory(cfg.Format, category, cfgFlags.NoColor)
		if err != nil {
			return fmt.Errorf("failed to create formatter: %w", err)
		}

		output, err := outputFormatter.Format(results)
		if err != nil {
			return fmt.Errorf("failed to format results: %w", err)
		}

		fmt.Print(output)

		// Handle browser opening flags
		if cfgFlags.Open || cfgFlags.OpenAll {
			if err := openResults(results, cfgFlags.OpenAll, cfgFlags.Verbose); err != nil {
				return fmt.Errorf("failed to open results in browser: %w", err)
			}
		}

		return nil
	}
}

// openResults opens search results in the browser
func openResults(results *searxnglib.SearchResponse, openAll bool, verbose bool) error {
	if len(results.Results) == 0 {
		return fmt.Errorf("no results to open")
	}

	// Collect URLs to open
	var urls []string
	if openAll {
		// Open all results
		for _, result := range results.Results {
			urls = append(urls, result.URL)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "\nOpening %d results in browser...\n", len(urls))
		}
	} else {
		// Open only the first result
		urls = append(urls, results.Results[0].URL)
		if verbose {
			fmt.Fprintf(os.Stderr, "\nOpening first result in browser: %s\n", results.Results[0].URL)
		}
	}

	// Check if browser opening is supported
	if !browser.IsSupported() {
		return fmt.Errorf("browser opening is not supported on this system")
	}

	// Open URLs
	if err := browser.OpenURLs(urls); err != nil {
		return err
	}

	return nil
}

// Execute runs the root command
func Execute() error {
	rootCmd := NewRootCommand()
	rootCmd.SetArgs(os.Args[1:])
	return rootCmd.Execute()
}

// cachedSearchClient wraps a CachedClient to implement the SearchWithConfig interface.
type cachedSearchClient struct {
	cached *cache.CachedClient
}

// SearchWithConfig executes a search using individual request parameters.
func (c *cachedSearchClient) SearchWithConfig(query string, results int, format string, category string, timeout int, language string, safeSearch int, page int, timeRange string) (*searxnglib.SearchResponse, error) {
	req := searxnglib.NewSearchRequest(query)
	req.Page = page
	req.Format = "json" // API always returns JSON
	req.Categories = []string{category}
	req.Languages = []string{language}
	req.SafeSearch = safeSearch
	req.TimeRange = timeRange
	return c.cached.Search(req)
}
