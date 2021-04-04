import { Action, Trigger } from "../../shared/types/integration";

export function isAction(item: Trigger | Action) {
  return item.inputs !== undefined;
}
