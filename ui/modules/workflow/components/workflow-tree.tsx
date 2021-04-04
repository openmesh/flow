import { useMemo } from "react";
import { Node, Workflow } from "../../shared/types/workflow";
import { WorkflowItem } from "./workflow-item";

export function WorkflowTree({
  workflow,
  addNode,
  selectNode,
}: {
  workflow: Workflow;
  addNode: (node: Node) => void;
  selectNode: (id: string) => void;
}) {
  // Get the top level node in the workflow graph.
  const rootId = useMemo(
    () => workflow.nodes.find((x) => x.parentIds.length === 0)?.id,
    [workflow]
  );

  return rootId ? (
    <WorkflowItem
      workflow={workflow}
      itemId={rootId}
      addNode={addNode}
      selectNode={selectNode}
    />
  ) : null;
}
