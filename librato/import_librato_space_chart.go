package librato

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/henrikhodne/go-librato/librato"
)

// Space import fans out to multiple resources due to the
// Space Charts. Instead of creating one resource with nested
// rules, we use the best practices approach of one resource per rule.
func resourceLibratoSpaceChartImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var result []*schema.ResourceData
	client := meta.(*librato.Client)

	// First query the space group
	idParts := strings.SplitN(d.Id(), ".", 2)
	spaceIDStr, chartIDStr := idParts[0], idParts[1]

	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 0)
	if err != nil {
		return nil, err
	}

	chartID, err := strconv.ParseUint(chartIDStr, 10, 0)
	if err != nil {
		return nil, err
	}

	spaceRaw, _, err := SpaceRefreshFunc(client, uint(spaceID))()
	if err != nil {
		return nil, err
	}
	if spaceRaw == nil {
		return nil, fmt.Errorf("space not found")
	}
	space := spaceRaw.(*librato.Space)

	// Fetch our chart to make sure it exists.
	// In the future, we could import all properties as well;
	// this is not necessary at this time so skipping that.
	_, _, err = client.Spaces.GetChart(*space.ID, uint(chartID))
	if err != nil {
		return nil, err
	}

	d.SetId(strconv.FormatUint(uint64(chartID), 10))
	d.SetType("librato_space_chart")
	d.Set("space_id", spaceID)
	result = append(result, d)
	return result, nil
}
