package runner

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping"
	"github.com/projectdiscovery/gologger"
	"golang.org/x/exp/maps"
)

func printStatistics(stats map[string]subscraping.Statistics) {

	sources := maps.Keys(stats)
	sort.Strings(sources)

	var lines []string
	var skipped []string

	for _, source := range sources {
		sourceStats := stats[source]
		if sourceStats.Skipped {
			skipped = append(skipped, fmt.Sprintf(" %s", source))
		} else {
			lines = append(lines, fmt.Sprintf(" %-20s %-10s %10d %10d", source, sourceStats.TimeTaken.Round(time.Millisecond).String(), sourceStats.Results, sourceStats.Errors))
		}
	}

	if len(lines) > 0 {
		gologger.Print().Msgf("\n Source               Duration      Results     Errors\n%s\n", strings.Repeat("─", 56))
		gologger.Print().Msgf(strings.Join(lines, "\n"))
		gologger.Print().Msgf("\n")
	}

	if len(skipped) > 0 {
		gologger.Print().Msgf("\n The following sources were included but skipped...\n\n")
		gologger.Print().Msgf(strings.Join(skipped, "\n"))
		gologger.Print().Msgf("\n\n")
	}
}
