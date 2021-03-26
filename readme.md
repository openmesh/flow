# OpenMesh Events

## Draft

### Flow

- User adds application 
- Each application has one or more events
- Once an application has been added then a sensor can be created
- User creates a sensor which listens for a given event. The sensor can filter results and trigger a workflow
- A workflow is made up of actions
- Actions have an input and output as well as a behaviour that is unique for each action 

Sounds good

## Workflow overview

- Persistent representation of the workflow gets created 
- Trigger comes in from event bus
- If it matches one of the user's stored workflow runner creates a concrete `workflow.Workflow` from the persisted version