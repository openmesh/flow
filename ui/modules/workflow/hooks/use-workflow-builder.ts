import produce from "immer";
import { useMemo, useState } from "react";
import { Node, Workflow } from "../../shared/types/workflow";

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
        const updated = produce(workflow, (draft) => {
          // Add new node to array
          draft.nodes.push(node);
          // If the node has a parent then we also need to attach the new node to
          // it.
          if (node.parentIds.length > 0) {
            const parentNode = draft.nodes.find(
              (n) => n.id === node.parentIds[0]
            );
            parentNode.childrenIds.push(node.id);
          }
        });

        setWorkflow(updated);
      },
      updateParamValue: (nodeId: string, key: string, value: string) => {
        const updated = produce(workflow, (draft) => {
          const nodeIndex = draft.nodes.findIndex((n) => n.id === nodeId);
          if (nodeIndex === -1) {
            return draft;
          }
          const paramIndex = draft.nodes[nodeIndex].params.findIndex(
            (p) => p.key === key
          );
          if (paramIndex === -1) {
            return draft;
          }
          draft.nodes[nodeIndex].params[paramIndex].value = value;
          console.log(draft);
        });
        console.log(updated);
        setWorkflow(updated);
      },
    }),
    [workflow]
  );

  return [workflow, handlers];
}
