export interface Integration {
  label: string;
  description: string;
  key: string;
  baseUrl: string;
  triggers: Trigger[];
  actions: Action[];
}

export interface Trigger {
  key: string;
  label: string;
  description: string;
  endpoint: string;
  method: string;
  inputs: InputField[];
  outputs: OutputField[];
}

export interface Action {
  key: string;
  label: string;
  description: string;
  endpoint: string;
  method: string;
  inputs: InputField[];
  outputs: OutputField[];
}

export interface InputField {
  key: string;
  label: string;
  description: string;
  required: boolean;
  type: FieldType;
  default: string;
  example: string;
}

export interface OutputField {
  label: string;
  key: string;
  description: string;
  type: FieldType;
  path: string;
}

export type FieldType =
  | "number"
  | "string"
  | "boolean"
  | "datetime"
  | "complex";
