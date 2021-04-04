import { useMemo, useState } from "react";

export function useNodeSelector(): [string, typeof handlers] {
  const [selectedNodeId, setSelectedNodeId] = useState<string>();

  const handlers = useMemo(
    () => ({
      setSelectedNodeId,
      deselectNode: () => setSelectedNodeId(undefined),
    }),
    []
  );

  return [selectedNodeId, handlers];
}
