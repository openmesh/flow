# OpenMesh Events

## Draft

### Flow

- User adds application 
- Each application has one or more events
- Once an application has been added, then a sensor can be created
- User creates a sensor which listens for a given event. The sensor can filter results and trigger a workflow
- A workflow is made up of actions
- Actions have an input and output as well as a behaviour that is unique for each action 

Sounds good

## Workflow overview

- Persistent representation of the workflow gets created 
- Trigger comes in from event bus
- If it matches one of the user's stored workflow runner creates a concrete `workflow.Workflow` from the persisted version

## Looking into being able to pass ids on create / update

#### Update workflow 

- When a workflow is updated if at least one node is included in the payload then the existing nodes
  will be replaced with the given nodes.
  
#### Update node

- If at least one parent id or child id is provided with the payload then the parent and child
  edges should be replaced with the new ones specified in the payload.
  
- If at least one param is included in the param then the existing params should be replaced with
  ones supplied in the payload.
  

#### Auth endpoints

- `/auth/oauth/{auth_source}` create auth
- `/auth/signup`

# MVP TODO

## Backlog


## In progress


## Complete

- [x] Implement authorization
- [x] Implement authentication
- [x] Implement sign up flow
- [x] Allow users to create a workflow while supplying a UUID
