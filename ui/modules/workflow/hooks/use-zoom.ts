import { useMemo, useState } from "react";

export function useZoom(): [number, typeof handlers] {
  const [zoom, setZoom] = useState(100);

  const handlers = useMemo(
    () => ({
      zoomIn: () => {
        setZoom(zoom + 10);
      },
      zoomOut: () => {
        setZoom(zoom - 10);
      },
    }),
    [zoom]
  );

  return [zoom, handlers];
}
