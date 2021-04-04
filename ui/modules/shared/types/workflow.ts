import { InputField } from "./integration";

export interface Workflow {
  name: string;
  description: string;
  nodes: Node[];
}

export interface Node {
  id: string;
  integration: string;
  action: string;
  label: string;
  description: string;
  params: Param[];
  parentIds: string[];
  childrenIds: string[];
  type: "ACTION" | "TRIGGER";
}

export interface Param {
  key: string;
  value: string;
  type: "value" | "reference";
  input: InputField;
}
