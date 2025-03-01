package passive

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/alienvault"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/anubis"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/bevigil"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/binaryedge"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/bufferover"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/c99"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/censys"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/certspotter"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/chaos"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/chinaz"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/commoncrawl"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/crtsh"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/digitorus"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/dnsdb"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/dnsdumpster"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/dnsrepo"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/fofa"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/fullhunt"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/github"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/hackertarget"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/hunter"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/intelx"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/leakix"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/netlas"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/passivetotal"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/quake"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/rapiddns"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/riddler"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/robtex"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/securitytrails"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/shodan"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/sitedossier"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/threatbook"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/virustotal"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/waybackarchive"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/whoisxmlapi"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/zoomeye"
	"github.com/dtsmcnee/subfinder/v2/pkg/subscraping/sources/zoomeyeapi"
	"github.com/projectdiscovery/gologger"
)

var AllSources = [...]subscraping.Source{
	&alienvault.Source{},
	&anubis.Source{},
	&bevigil.Source{},
	&binaryedge.Source{},
	&bufferover.Source{},
	&c99.Source{},
	&censys.Source{},
	&certspotter.Source{},
	&chaos.Source{},
	&chinaz.Source{},
	&commoncrawl.Source{},
	&crtsh.Source{},
	&digitorus.Source{},
	&dnsdb.Source{},
	&dnsdumpster.Source{},
	&dnsrepo.Source{},
	&fofa.Source{},
	&fullhunt.Source{},
	&github.Source{},
	&hackertarget.Source{},
	&hunter.Source{},
	&intelx.Source{},
	&netlas.Source{},
	&leakix.Source{},
	&passivetotal.Source{},
	&quake.Source{},
	&rapiddns.Source{},
	&riddler.Source{},
	&robtex.Source{},
	&securitytrails.Source{},
	&shodan.Source{},
	&sitedossier.Source{},
	&threatbook.Source{},
	&virustotal.Source{},
	&waybackarchive.Source{},
	&whoisxmlapi.Source{},
	&zoomeye.Source{},
	&zoomeyeapi.Source{},
	// &threatminer.Source{}, // failing  api
	// &reconcloud.Source{}, // failing due to cloudflare bot protection
}

var NameSourceMap = make(map[string]subscraping.Source, len(AllSources))

func init() {
	for _, currentSource := range AllSources {
		NameSourceMap[strings.ToLower(currentSource.Name())] = currentSource
	}
}

// Agent is a struct for running passive subdomain enumeration
// against a given host. It wraps subscraping package and provides
// a layer to build upon.
type Agent struct {
	sources []subscraping.Source
}

// New creates a new agent for passive subdomain discovery
func New(sourceNames, excludedSourceNames []string, useAllSources, useSourcesSupportingRecurse bool) *Agent {
	sources := make(map[string]subscraping.Source, len(AllSources))

	if useAllSources {
		maps.Copy(sources, NameSourceMap)
	} else {
		if len(sourceNames) > 0 {
			for _, source := range sourceNames {
				if NameSourceMap[source] == nil {
					gologger.Warning().Msgf("There is no source with the name: %s", source)
				} else {
					sources[source] = NameSourceMap[source]
				}
			}
		} else {
			for _, currentSource := range AllSources {
				if currentSource.IsDefault() {
					sources[currentSource.Name()] = currentSource
				}
			}
		}
	}

	if len(excludedSourceNames) > 0 {
		for _, sourceName := range excludedSourceNames {
			delete(sources, sourceName)
		}
	}

	if useSourcesSupportingRecurse {
		for sourceName, source := range sources {
			if !source.HasRecursiveSupport() {
				delete(sources, sourceName)
			}
		}
	}

	gologger.Debug().Msgf(fmt.Sprintf("Selected source(s) for this search: %s", strings.Join(maps.Keys(sources), ", ")))

	// Create the agent, insert the sources and remove the excluded sources
	agent := &Agent{sources: maps.Values(sources)}

	return agent
}
