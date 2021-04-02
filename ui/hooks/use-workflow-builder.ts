import { useMemo, useState } from "react";
import { Node, Workflow } from "../types/workflow";

export function useWorkflowBuilder(): [Workflow, typeof handlers] {
  const [workflow, setWorkflow] = useState<Workflow>({
    description: "",
    name: "",
    nodes: [],
  });

  const handlers = useMemo(
    () => ({
      addNode: (node: Node) => {
        // Create new array of nodes
        const nodes = [...workflow.nodes, node];
        // If the node has a parent then we also need to attach the new node to
        // it.
        if (node.parentIds.length > 0) {
          const parentNode = nodes.find((n) => n.id === node.parentIds[0]);
          parentNode.childrenIds.push(node.id);
        }

        setWorkflow({ ...workflow, nodes });
      },
    }),
    [workflow]
  );

  return [workflow, handlers];
}
