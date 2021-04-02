import camelcaseKeys from "camelcase-keys";
import { Integration } from "../types/integration";

export async function getIntegrations(): Promise<{
  data: Integration[];
  totalItems: number;
}> {
  const res = await fetch(`/api/v1/integrations`);
  if (res.status !== 200) {
    throw new Error("Failed to fetch integrations");
  }
  return camelcaseKeys(await res.json(), { deep: true });
}
