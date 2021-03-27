package pg

import (
	"context"
	"fmt"
	"github.com/openmesh/flow"
)

func getNodes(ctx context.Context, tx *Tx, req flow.GetNodesRequest) ([]*flow.Node, int, error) {
	// build query to be executed
	baseQuery := fmt.Sprintf(
		"SELECT * FROM nodes %s",
		where("workflow_id = '%s'", req.WorkflowID),
	)

	// Get count of base query.
	var count int
	err := tx.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count;", baseQuery))
	if err != nil {
		return nil, 0, err
	}

	// Append limit and offset to query if required.
	paginatedQuery := fmt.Sprintf("%s %s", baseQuery, formatLimitOffset(req.Limit, req.Page))
	nodes := make([]*flow.Node, 0)

	if err := tx.Select(&nodes, paginatedQuery); err != nil {
		return nodes, 0, err
	}

	return nodes, count, nil
}
