package main

import (
	"fmt"
	"strings"

	"github.com/bugsnag/bugsnag-go/v2"
)

const metadataPrefix string = "BUGSNAG_METADATA_"
const metadataPrefixLen int = len(metadataPrefix)

func addMetadata(event *bugsnag.Event, keypath string, value string) {
	if len(keypath) == 0 {
		return
	}
	tab, key := splitTabKeyValues(keypath)
	event.MetaData.Add(formatMetadataKey(tab), formatMetadataKey(key), value)
}

func splitTabKeyValues(keypath string) (string, string) {
	key_components := strings.SplitN(keypath, ".", 2)
	if len(key_components) > 1 {
		return key_components[0], key_components[1]
	}
	return "custom", keypath
}

func parseMetadataKeypath(key string) (string, error) {
	if strings.HasPrefix(key, metadataPrefix) && len(key) > metadataPrefixLen {
		return strings.TrimPrefix(key, metadataPrefix), nil
	}
	return "", fmt.Errorf("No metadata prefix found")
}

func formatMetadataKey(key string) string {
	return strings.Replace(key, "_", " ", -1)
}
