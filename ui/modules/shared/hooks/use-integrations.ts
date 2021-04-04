import { useEffect, useState } from "react";
import { getIntegrations } from "../../../api/integrations";
import { Integration } from "../types/integration";

export function useIntegrations() {
  const [integrations, setIntegrations] = useState<Integration[]>([]);

  async function fetchIntegrations() {
    const res = await getIntegrations();
    setIntegrations(res.data);
  }

  useEffect(() => {
    fetchIntegrations();
  }, []);

  return integrations;
}
