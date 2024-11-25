package git

import (
	"bytes"
	"emperror.dev/errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/je4/FairService/v2/pkg/fair"
	spdxjson "github.com/spdx/tools-golang/json"
	spdxv23 "github.com/spdx/tools-golang/spdx/v2/v2_3"
	"io"
	"net/http"
	"strings"
)

type Plugin struct {
}

func (g *Plugin) Handle(_fair *fair.Fair, pid string, item *fair.ItemData) (*fair.PluginResult, error) {
	naan, qualifier, components, variants, inflection, err := fair.ArkParts(pid)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse PID")
	}
	// cannot deal with inflections
	if inflection != "" {
		return &fair.PluginResult{
			Type: fair.ARKPluginCannotHandle,
		}, nil
	}
	source, err := _fair.GetSourceByName(item.Partition, item.Source)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get source %s", item.Source)
	}

	url := item.URL
	if url == "" {
		url = strings.ReplaceAll(source.DetailURL, "{signature}", item.Signature)
	}
	suffix := components
	if variants != "" {
		suffix += "." + variants
	}
	componentParts := strings.SplitN(suffix, "/", 2)
	if len(componentParts) != 2 {
		return nil, errors.Errorf("invalid suffix %s", suffix)
	}
	version := componentParts[0]
	file := componentParts[1]

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
	if file != "SPDX" {
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
	if resp.StatusCode == http.StatusOK {
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read response body")
		}
	} else {
		return nil, errors.Errorf("cannot get %s: %s", rawUrl, resp.Status)
	}
	spdxDocument, err := spdxjson.Read(bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read spdx.json")
	}
	for _, pkg := range spdxDocument.Packages {
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
			continue
		}

	}

	return nil, nil
}

var _ fair.Plugin = (*Plugin)(nil)
