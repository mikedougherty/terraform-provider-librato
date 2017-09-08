package librato

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/henrikhodne/go-librato/librato"
)

// Space import fans out to multiple resources due to the
// Space Charts. Instead of creating one resource with nested
// rules, we use the best practices approach of one resource per rule.
func resourceLibratoSpaceImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*librato.Client)

	// First query the space group
	spaceID, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		return nil, err
	}

	spaceRaw, _, err := SpaceRefreshFunc(client, uint(spaceID))()

	if err != nil {
		return nil, err
	}

	space := spaceRaw.(*librato.Space)
	spaceCharts, _, err := client.Spaces.ListCharts(*space.ID)

	// Start building our results
	results := make([]*schema.ResourceData, 1,
		1+len(spaceCharts))
	results[0] = d

	for _, chartSlug := range spaceCharts {
		chartResource := resourceLibratoSpaceChart()
		chart := chartResource.Data(nil)
		chart.Set("space_id", spaceID)
		// For importing, we use the form '{spaceID}.{chartID}' so it can be found properly.
		chart.SetId(fmt.Sprintf("%s.%s", strconv.FormatUint(uint64(*space.ID), 10), strconv.FormatUint(uint64(*chartSlug.ID), 10)))
		chartResults, err := resourceLibratoSpaceChartImportState(chart, meta)
		if err != nil {
			return nil, err
		}
		chart.SetId(strconv.FormatUint(uint64(*chartSlug.ID), 10))
		results = append(results, chartResults...)
	}

	return results, nil
}
