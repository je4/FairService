package git

import (
	"bytes"
	"emperror.dev/errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/package-url/packageurl-go"
	spdxjson "github.com/spdx/tools-golang/json"
	spdxv23 "github.com/spdx/tools-golang/spdx/v2/v2_3"
	spdxyaml "github.com/spdx/tools-golang/yaml"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
)

func NewPlugin(logger zLogger.ZLogger) *Plugin {
	return &Plugin{logger: logger}
}

type Plugin struct {
	logger zLogger.ZLogger
}

var purlPKGRegexpWithVersion = regexp.MustCompile(`(?i)^(?P<base>pkg:golang.+)(?P<goversion>\/v[\d]+)(?P<version>@.*)?$`)

func (g *Plugin) Handle(_fair *fair.Fair, pid string, item *fair.ItemData) (*fair.PluginResult, error) {
	naan, qualifier, components, variants, inflection, err := fair.ArkParts(pid)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse PID")
	}
	_ = naan
	_ = qualifier
	inflection = strings.ToLower(inflection)
	source, err := _fair.GetSourceByName(item.Partition, item.Source)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get source %s", item.Source)
	}

	part, err := _fair.GetPartition(item.Partition)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get partition %s", item.Partition)
	}

	url := item.URL
	if url == "" {
		url = strings.ReplaceAll(source.DetailURL, "{signature}", item.Signature)
	}
	suffix := components
	if variants != "" {
		suffix += "." + variants
	}
	if suffix == "" {
		return &fair.PluginResult{
			Type: fair.ARKPluginRedirect,
			Data: []byte(url),
		}, nil
	}
	componentParts := strings.SplitN(suffix, "/", 2)
	if len(componentParts) < 1 {
		return nil, errors.Errorf("invalid suffix %s", suffix)
	}
	version := componentParts[0]
	var file string
	if len(componentParts) >= 2 {
		file = componentParts[1]
	}

	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	refs, err := rem.List(&git.ListOptions{
		// Returns all references, including peeled references.
		PeelingOption: git.AppendPeeled,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "cannot list refs")
	}
	// Filters the references list and only keeps tags
	var foundTag string
	var hash string
	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tag := ref.Name().Short()
			tags = append(tags, tag)
			tag = strings.ToLower(tag)
			if tag == version || tag == "v"+version {
				foundTag = tag
				hash = ref.Hash().String()
			}
		}
	}
	if foundTag == "" {
		return nil, errors.Errorf("cannot find tag %s in %v", version, tags)
	}
	if !slices.Contains([]string{"?spdx", "?spdx.yaml", "?spdx.json", "?spdx.gv"}, inflection) {
		return &fair.PluginResult{
			Type: fair.ARKPluginRedirect,
			Data: []byte(fmt.Sprintf("%s/blob/%s/%s", url, hash, file)),
		}, nil
	}

	file = "spdx.json"
	var data []byte
	var rawUrl string
	switch {
	case strings.Contains(url, "github.com"):
		rawUrl = fmt.Sprintf("%s/%s/%s", strings.ReplaceAll(url, "github.com", "raw.githubusercontent.com"), hash, file)
	case strings.Contains(url, "gitlab.com"):
		rawUrl = fmt.Sprintf("%s/-/raw/%s/%s", url, foundTag, file)
	default:
		return nil, errors.Errorf("cannot handle %s", url)
	}
	resp, err := http.Get(rawUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get %s", rawUrl)
	}
	data, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("cannot get %s: %s", rawUrl, resp.Status)
	}
	var infoPID = map[string]int64{}
	var infoLicense = map[string]int64{}
	var pidCounter int64
	var withPID = []string{}
	spdxDocument, err := spdxjson.Read(bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read spdx.json")
	}
	for _, pkg := range spdxDocument.Packages {
		var hasPID bool
		if pkg.PackageLicenseDeclared != "" {
			if _, ok := infoLicense[pkg.PackageLicenseDeclared]; !ok {
				infoLicense[pkg.PackageLicenseDeclared] = 0
			}
			infoLicense[pkg.PackageLicenseDeclared]++
		} else if pkg.PackageLicenseConcluded != "" {
			if _, ok := infoLicense[pkg.PackageLicenseConcluded]; !ok {
				infoLicense[pkg.PackageLicenseConcluded] = 0
			}
			infoLicense[pkg.PackageLicenseConcluded]++
		}
		/*
			if pkg.PackageName == spdxDocument.DocumentName {
				if pkg.PackageExternalReferences == nil {
					pkg.PackageExternalReferences = []*spdxv23.PackageExternalReference{}
				}
				pkg.PackageExternalReferences = append(pkg.PackageExternalReferences, &spdxv23.PackageExternalReference{
					Category:           "PERSISTENT-ID",
					RefType:            "ark",
					Locator:            fmt.Sprintf("ark:/%s/%s", naan, qualifier),
					ExternalRefComment: "",
				})
			}
		*/
		for _, extRefs := range pkg.PackageExternalReferences {
			switch extRefs.Category {
			case "PERSISTENT-ID":
				if _, ok := infoPID[extRefs.RefType]; !ok {
					infoPID[extRefs.RefType] = 0
				}
				infoPID[extRefs.RefType]++
				continue
			case "PACKAGE-MANAGER":
				if extRefs.RefType == "purl" {
					// get rid of golang versions
					if matches := purlPKGRegexpWithVersion.FindStringSubmatch(extRefs.Locator); matches != nil {
						extRefs.Locator = matches[1] + matches[3]
					}
					packageURL, err := packageurl.FromString(extRefs.Locator)
					if err != nil {
						g.logger.Error().Err(err).Msgf("cannot parse packageurl %s", extRefs.Locator)
						continue
					}
					var signature string
					if packageURL.Type == "golang" {
						nsParts := strings.SplitN(packageURL.Namespace, "/", 2)
						if len(nsParts) < 2 {
							g.logger.Error().Msgf("invalid namespace %s", packageURL.Namespace)
							continue
						}
						signature = nsParts[1] + "/" + packageURL.Name
					} else if packageURL.Type == "github" {
						signature = packageURL.Namespace + "/" + packageURL.Name
					}
					if signature != "" {
						i, err := _fair.GetItemSource(part, source.ID, signature)
						if err != nil {
							g.logger.Info().Err(err).Msgf("cannot get item %s", signature)
							continue
						}
						if i == nil {
							g.logger.Debug().Msgf("cannot find item %s", signature)
							continue
						}
						for _, identifier := range i.Identifier {
							parts := strings.SplitN(identifier, ":", 2)
							if len(parts) < 2 {
								g.logger.Error().Msgf("invalid identifier %s", identifier)
								continue
							}
							if _, ok := infoPID[parts[0]]; !ok {
								infoPID[parts[0]] = 0
							}
							infoPID[parts[0]]++
							pkg.PackageExternalReferences = append(pkg.PackageExternalReferences, &spdxv23.PackageExternalReference{
								Category: "PERSISTENT-ID",
								RefType:  parts[0],
								Locator:  identifier,
							})
							hasPID = true
						}
					}

				}
			}
		}
		if pkg.PackageDownloadLocation != "" {
			var gitDownloadLocation1 = regexp.MustCompile(`(?i)^(.+)(\/v\d+|\.git|\/releases\/.+)$`)
			var gitDownloadLocation2 = regexp.MustCompile(`(?i)^(git\+)?(?P<base>https:\/\/[^\/]+\/)(?P<namespace>.+)\/(?P<name>[^\/]+)$`)
			if matches := gitDownloadLocation1.FindStringSubmatch(pkg.PackageDownloadLocation); matches != nil {
				pkg.PackageDownloadLocation = matches[1]
				if matches := gitDownloadLocation1.FindStringSubmatch(pkg.PackageDownloadLocation); matches != nil {
					pkg.PackageDownloadLocation = matches[1]
				}
			}
			if matches := gitDownloadLocation2.FindStringSubmatch(pkg.PackageDownloadLocation); matches != nil {
				pkg.PackageDownloadLocation = matches[2] + matches[3] + "/" + matches[4]
				signature := strings.ToLower(matches[3] + "/" + matches[4])
				i, err := _fair.GetItemSource(part, source.ID, signature)
				if err != nil {
					g.logger.Debug().Err(err).Msgf("cannot get item %s", signature)
					continue
				}
				if i == nil {
					g.logger.Debug().Msgf("cannot find item %s", signature)
					continue
				}
				for _, identifier := range i.Identifier {
					parts := strings.SplitN(identifier, ":", 2)
					if len(parts) < 2 {
						g.logger.Error().Msgf("invalid identifier %s", identifier)
						continue
					}
					if _, ok := infoPID[parts[0]]; !ok {
						infoPID[parts[0]] = 0
					}
					infoPID[parts[0]]++
					pkg.PackageExternalReferences = append(pkg.PackageExternalReferences, &spdxv23.PackageExternalReference{
						Category: "PERSISTENT-ID",
						RefType:  parts[0],
						Locator:  identifier,
					})
					hasPID = true
				}
			}
		}
		if hasPID {
			withPID = append(withPID, pkg.PackageName)
			pidCounter++
		}

	}

	if inflection == "?spdx" {
		str := ""
		str += fmt.Sprintf("Packages: %d\n", len(spdxDocument.Packages))
		str += fmt.Sprintf("PID count: %d\n", pidCounter)
		for t, count := range infoPID {
			str += fmt.Sprintf("PID %s: %d\n", t, count)
		}
		for t, count := range infoLicense {
			str += fmt.Sprintf("License \"%s\": %d\n", t, count)
		}
		return &fair.PluginResult{
			Type: fair.ARKPluginData,
			Data: []byte(str),
			Mime: "text/plain",
		}, nil
	} else if inflection == "?spdx.yaml" {
		var buf = &bytes.Buffer{}
		if err := spdxyaml.Write(spdxDocument, buf); err != nil {
			return nil, errors.Wrap(err, "cannot write spdx.json")
		}
		return &fair.PluginResult{
			Type: fair.ARKPluginData,
			Data: buf.Bytes(),
			Mime: "application/x-yaml",
		}, nil
	} else if inflection == "?spdx.json" {
		var buf = &bytes.Buffer{}
		if err := spdxjson.Write(spdxDocument, buf); err != nil {
			return nil, errors.Wrap(err, "cannot write spdx.json")
		}
		return &fair.PluginResult{
			Type: fair.ARKPluginData,
			Data: buf.Bytes(),
			Mime: "application/json",
		}, nil
	} else if inflection == "?spdx.gv" {
		str := `digraph mygraph {
  fontname="Helvetica,Arial,sans-serif"
  node [fontname="Helvetica,Arial,sans-serif"]
  edge [fontname="Helvetica,Arial,sans-serif"]
  node [shape=box];
`
		for _, pkg := range spdxDocument.Packages {
			if slices.Contains(withPID, pkg.PackageName) {
				str += fmt.Sprintf("  \"%s\" [label=\"%s\" style=filled,color=aquamarine];\n", pkg.PackageSPDXIdentifier, pkg.PackageName)
			} else {
				str += fmt.Sprintf("  \"%s\" [label=\"%s\"];\n", pkg.PackageSPDXIdentifier, pkg.PackageName)
			}
		}
		for _, rel := range spdxDocument.Relationships {
			str += fmt.Sprintf("  \"%s\" -> \"%s\";\n", rel.RefA.ElementRefID, rel.RefB.ElementRefID)
		}
		str += "}\n"
		return &fair.PluginResult{
			Type: fair.ARKPluginData,
			Data: []byte(str),
			Mime: "text/vnd.graphviz",
		}, nil

	} else {
		return &fair.PluginResult{
			Type: fair.ARKPluginCannotHandle,
		}, nil
	}
}

var _ fair.Plugin = (*Plugin)(nil)
