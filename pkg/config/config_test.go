package config

import (
	"strings"
	"testing"
)

func TestConfigParsing(t *testing.T) {
	rawConfig := `version: v1beta1
export:  # Old from virtual cluster
- apiVersion: cert-manager.io/v1
  kind: Issuer
- apiVersion: cert-manager.io/v1
  kind: Certificates
  patches:
    - op: rewriteName
      path: spec.ca.secretName
  reversePatches:       
    - op: copyFromObject # Sync status back by default
      fromPath: status
      path: status`

	_, err := Parse(rawConfig)
	if err != nil {
		t.Fatalf("Error parsing config %v", err)
	}
}

func TestConfigParsingHooks(t *testing.T) {
	rawConfig := `version: v1beta1
hooks:
  hostToVirtual:
    - apiVersion: v1
      kind: Pod
      verbs: ["create", "update", "patch"]
      patches:
        - op: add
          path: metadata.annotations
          value:
            import-annotation: testing-annotation-import
  virtualToHost:
    - apiVersion: v1
      kind: Pod
      verbs: ["create", "update", "patch"]
      patches:
        - op: add
          path: metadata.annotations
          value:
            export-annotation: testing-annotation-export
`

	_, err := Parse(rawConfig)
	if err != nil {
		t.Fatalf("Error parsing config %v", err)
	}
}

func TestConfigParsingHooksUnknownVerb(t *testing.T) {
	rawConfig := `version: v1beta1
hooks:
  hostToVirtual:
    - apiVersion: v1
      kind: Pod
      verbs: ["create", "update", "patch", "unknown"]
      patches:
        - op: add
          path: metadata.annotations
          value:
            import-annotation: testing-annotation-import
  virtualToHost:
    - apiVersion: v1
      kind: Pod
      verbs: ["create", "update", "patch"]
      patches:
        - op: add
          path: metadata.annotations
          value:
            export-annotation: testing-annotation-export
`

	_, err := Parse(rawConfig)
	if err == nil {
		t.Fatalf("Error parsing config %v", err)
	}

	if !strings.Contains(err.Error(), "invalid verb \"unknown\";") {
		t.Fatalf("Error parsing config %v", err)
	}
}
