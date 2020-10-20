package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// V3Droplet is the result of staging an application package.
// There are two types (lifecycles) of droplets: buildpack and
// docker. In the case of buildpacks, the droplet contains the
// bits produced by the buildpack.
type V3Droplet struct {
	State             string            `json:"state,omitempty"`
	Error             string            `json:"error,omitempty"`
	Lifecycle         V3Lifecycle       `json:"lifecycle,omitempty"`
	GUID              string            `json:"guid,omitempty"`
	CreatedAt         string            `json:"created_at,omitempty"`
	UpdatedAt         string            `json:"updated_at,omitempty"`
	Links             map[string]Link   `json:"links,omitempty"`
	ExecutionMetadata string            `json:"execution_metadata,omitempty"`
	ProcessTypes      map[string]string `json:"process_types,omitempty"`
	Metadata          V3Metadata        `json:"metadata,omitempty"`

	// Only specified when the droplet is using the Docker lifecycle.
	Image string `json:"image,omitempty"`

	// The following fields are specified when the droplet is using
	// the buildpack lifecycle.
	Checksum struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"checksum,omitempty"`
	Stack      string                `json:"stack,omitempty"`
	Buildpacks []V3DetectedBuildpack `json:"buildpacks,omitempty"`
}

type V3DetectedBuildpack struct {
	Name          string `json:"name,omitempty"`           // system buildpack name
	BuildpackName string `json:"buildpack_name,omitempty"` // name reported by the buildpack
	DetectOutput  string `json:"detect_output,omitempty"`  // output during detect process
	Version       string `json:"version,omitempty"`
}

type SetCurrentDropletV3Response struct {
	Data  V3Relationship  `json:"data,omitempty"`
	Links map[string]Link `json:"links,omitempty"`
}

func (c *Client) SetCurrentDropletForV3App(guid string) (*SetCurrentDropletV3Response, error) {
	req := c.NewRequest("PATCH", "/v3/apps/"+guid+"/relationships/current_droplet")
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error setting droplet for v3 app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error setting droplet for v3 app with GUID [%s], response code: %d", guid, resp.StatusCode)
	}

	var r SetCurrentDropletV3Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading droplet response JSON")
	}

	return &r, nil
}
