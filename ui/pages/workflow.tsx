import { Input, InputGroup, InputLeftElement } from "@chakra-ui/input";
import { Box, Center, Flex, Stack, Text } from "@chakra-ui/layout";
import { useTheme } from "@chakra-ui/system";
import { Tab, TabList, TabPanel, TabPanels, Tabs } from "@chakra-ui/tabs";
import { FilterIcon } from "@heroicons/react/outline";
import { LightningBoltIcon } from "@heroicons/react/solid";
import { useMemo, useRef, useState } from "react";
import { ArcherContainer, ArcherElement } from "react-archer";
import { useDrag, useDrop } from "react-dnd";
import { v4 as uuidv4 } from "uuid";
import { DropdownCombobox } from "../components/dropdown-combobox";
import { useIntegrations } from "../hooks/use-integrations";
import { useWorkflowBuilder } from "../hooks/use-workflow-builder";
import { Action, Integration, Trigger } from "../types/integration";
import { Node, Workflow } from "../types/workflow";

export default () => {
  const [workflow, { addNode }] = useWorkflowBuilder();

  const [{ canDrop, isOver }, drop] = useDrop(
    () => ({
      accept: "WORKFLOW_ITEM",
      canDrop: () => workflow.nodes.length === 0,
      collect: (monitor) => ({
        isOver: monitor.isOver() && monitor.canDrop(),
        canDrop: monitor.canDrop(),
      }),
      drop: (item: (Action | Trigger) & { type: "ACTION" | "TRIGGER" }) => {
        addNode({
          id: uuidv4(),
          action: item.key,
          label: item.label,
          description: item.description,
          childrenIds: [],
          parentIds: [],
          params: [],
          integration: "",
          type: item.type,
        });
      },
    }),
    [workflow]
  );

  const archerContainerRef = useRef<ArcherContainer>();

  return (
    <ArcherContainer
      endShape={{
        arrow: {
          arrowLength: 3,
        },
      }}
      ref={archerContainerRef}
    >
      <Box bg="gray.100" minH="100vh" display="flex" flexDirection="column">
        <Box
          as="header"
          bg="white"
          p="4"
          borderBottomColor="gray.300"
          borderBottomWidth="1px"
        >
          <Flex justify="space-between">
            <h1>Create workflow</h1>
          </Flex>
        </Box>
        <Box display="flex" maxH="full" flex="1">
          <Sidebar hasTrigger={workflow.nodes.length > 0} />
          <Center
            ref={drop}
            role="Dustbin"
            bg={isOver && canDrop ? "blue.50" : undefined}
            onScroll={() => archerContainerRef.current?.refreshScreen()}
            overflow="auto"
            flex="1"
          >
            <WorkflowTree workflow={workflow} addNode={addNode} />
          </Center>
        </Box>
      </Box>
    </ArcherContainer>
  );
};

function Sidebar({ hasTrigger }: { hasTrigger: boolean }) {
  const integrations = useIntegrations();
  const [selectedIntegration, setSelectedIntegration] = useState<Integration>();

  const theme = useTheme();

  return (
    <Box bg="white" maxW="xs">
      <Stack px="4" py="6" spacing="4">
        <DropdownCombobox
          items={integrations}
          onChange={(change) => setSelectedIntegration(change.selectedItem)}
        />
        <Tabs isFitted>
          <TabList>
            <Tab isDisabled={hasTrigger}>Triggers</Tab>
            <Tab isDisabled={!hasTrigger}>Actions</Tab>
          </TabList>
          <TabPanels>
            <TabPanel px="0">
              <Stack>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    children={
                      <FilterIcon height="20" color={theme.colors.gray[300]} />
                    }
                  />
                  <Input placeholder="Filter triggers" />
                </InputGroup>
                {selectedIntegration?.triggers?.map((trigger) => (
                  <SidebarItem item={trigger} type="TRIGGER" />
                ))}
              </Stack>
            </TabPanel>
            <TabPanel px="0">
              <Stack>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    children={
                      <FilterIcon height="20" color={theme.colors.gray[300]} />
                    }
                  />
                  <Input placeholder="Filter actions" />
                </InputGroup>
                {selectedIntegration?.actions?.map((action) => (
                  <SidebarItem item={action} type="ACTION" />
                ))}
              </Stack>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Stack>
    </Box>
  );
}

function SidebarItem({
  item,
  type,
}: {
  item: Action | Trigger;
  type: "ACTION" | "TRIGGER";
}) {
  const theme = useTheme();
  const [{ isDragging }, drag, dragPreview] = useDrag(() => ({
    type: "WORKFLOW_ITEM",
    collect: (monitor) => ({
      isDragging: monitor.isDragging(),
    }),
    item: {
      ...item,
      type,
    },
  }));

  return (
    <Box
      _hover={{ shadow: "md" }}
      p="4"
      borderRadius="md"
      ref={dragPreview}
      opacity={isDragging ? "0.5" : 1}
    >
      <Stack direction="row" alignItems="flex-start" spacing="4" ref={drag}>
        <Box p="2" bg="gray.200" borderRadius="md">
          <LightningBoltIcon height="20" color={theme.colors.gray[400]} />
        </Box>
        <Box>
          <Text fontWeight="semibold" fontSize="md" textColor="gray.900">
            {item.label}
          </Text>
          <Text textColor="gray.500" fontSize="sm">
            {item.description}
          </Text>
        </Box>
      </Stack>
    </Box>
  );
}

function WorkflowItem({
  workflow,
  itemId,
  addNode,
}: {
  workflow: Workflow;
  itemId: string;
  addNode: (node: Node) => void;
}) {
  const theme = useTheme();

  const item = useMemo(() => workflow.nodes.find((x) => x.id === itemId), [
    workflow,
  ]);

  const [{ canDrop, isOver }, drop] = useDrop(
    () => ({
      accept: "WORKFLOW_ITEM",
      collect: (monitor) => ({
        isOver: monitor.isOver() && monitor.canDrop(),
        canDrop: monitor.canDrop(),
      }),
      drop: (item: (Action | Trigger) & { type: "ACTION" | "TRIGGER" }) => {
        addNode({
          action: item.key,
          childrenIds: [],
          description: item.description,
          id: uuidv4(),
          integration: "",
          label: item.label,
          params: [],
          parentIds: [itemId],
          type: "ACTION",
        });
      },
    }),
    [workflow]
  );

  const childNodes = useMemo(
    () => workflow.nodes.filter((n) => item.childrenIds.includes(n.id)),
    [workflow]
  );

  return (
    <Stack
      justifyContent="start"
      alignItems="center"
      minW="lg"
      spacing="8"
      margin="auto"
    >
      <ArcherElement
        id={`workflow-item-${item.id}`}
        relations={item.childrenIds.map((childId) => ({
          targetId: `workflow-item-${childId}`,
          sourceAnchor: "bottom",
          targetAnchor: "top",
          style: {
            strokeColor: theme.colors.gray[300],
          },
        }))}
      >
        <Box
          bg="white"
          shadow="sm"
          w="lg"
          borderRadius="md"
          _hover={{ shadow: "outline" }}
          position="relative"
        >
          <Box
            height="0.75rem"
            width="0.75rem"
            bottom="-0.375rem"
            bg="blue.500"
            shadow="0 0 10px #3182CE"
            left="calc(50% - 0.375rem)"
            borderRadius="full"
            position="absolute"
            transform={isOver ? undefined : "scale(0)"}
            transition="transform"
            transitionDuration="200ms"
          />
          <Box
            id={`${item.id}-child-drop-area`}
            ref={drop}
            position="absolute"
            top="50%"
            bottom="-50%"
            left="0"
            right="0"
          />
          <Box
            p="4"
            borderRadius="md"
            display="flex"
            justifyContent="space-between"
            alignItems="center"
          >
            <Stack direction="row" alignItems="flex-start" spacing="4">
              <Box bg="blue.50" borderRadius="md" p="2">
                <LightningBoltIcon height="24" color={theme.colors.blue[400]} />
              </Box>
              <Stack spacing="1">
                <Text color="gray.900" fontWeight="semibold" fontSize="xl">
                  {item.label}
                </Text>
                <Text color="gray.500">{item.description}</Text>
              </Stack>
            </Stack>
          </Box>
        </Box>
      </ArcherElement>
      <Stack
        direction="row"
        justifyContent="start"
        alignItems="flex-start"
        spacing="8"
        margin="auto"
      >
        {childNodes.map((n) => (
          <WorkflowItem workflow={workflow} itemId={n.id} addNode={addNode} />
        ))}
      </Stack>
    </Stack>
  );
}

function WorkflowTree({
  workflow,
  addNode,
}: {
  workflow: Workflow;
  addNode: (node: Node) => void;
}) {
  // Get the top level node in the workflow graph.
  const rootId = useMemo(
    () => workflow.nodes.find((x) => x.parentIds.length === 0)?.id,
    [workflow]
  );

  return rootId ? (
    <WorkflowItem workflow={workflow} itemId={rootId} addNode={addNode} />
  ) : null;
}
