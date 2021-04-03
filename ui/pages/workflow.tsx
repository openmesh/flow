import { ButtonGroup, IconButton } from "@chakra-ui/button";
import { Input, InputGroup, InputLeftElement } from "@chakra-ui/input";
import {
  Box,
  Center,
  Flex,
  Grid,
  GridItem,
  Stack,
  Text,
} from "@chakra-ui/layout";
import { useTheme } from "@chakra-ui/system";
import { Tab, TabList, TabPanel, TabPanels, Tabs } from "@chakra-ui/tabs";
import { FilterIcon, MinusIcon, PlusIcon } from "@heroicons/react/outline";
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
  const [zoom, { zoomIn, zoomOut }] = useZoom();
  console.log(zoom);

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
    <Box bg="gray.100" display="flex" flexDirection="column" h="100vh">
      <Box
        as="header"
        bg="white"
        p="4"
        borderBottomColor="gray.300"
        borderBottomWidth="1px"
        height="4rem"
        zIndex="1"
      >
        <Flex justify="space-between">
          <h1>Create workflow</h1>
        </Flex>
      </Box>
      <Box display="flex" flex="1" maxH="calc(100vh - 4rem)">
        <Sidebar hasTrigger={workflow.nodes.length > 0} />
        <Center
          ref={drop}
          role="Dustbin"
          bg={isOver && canDrop ? "blue.50" : undefined}
          onScroll={() => archerContainerRef.current?.refreshScreen()}
          overflow="auto"
          flex="1"
          position="relative"
        >
          <ButtonGroup
            position="absolute"
            top="2rem"
            left="2rem"
            isAttached
            variant="solid"
          >
            <IconButton
              icon={<MinusIcon height="20" />}
              aria-label="Zoom out"
              onClick={zoomOut}
            />
            <IconButton
              icon={<PlusIcon height="20" />}
              aria-label="Zoom out"
              onClick={zoomIn}
            />
          </ButtonGroup>
          <ArcherContainer
            endShape={{
              arrow: {
                arrowLength: 3,
              },
            }}
            ref={archerContainerRef}
            style={{
              zoom: `${zoom}%`,
              minWidth: "0",
              minHeight: "0",
              maxHeight: "100%",
              maxWidth: "100%",
              margin: "auto",
            }}
          >
            <WorkflowTree workflow={workflow} addNode={addNode} />
          </ArcherContainer>
        </Center>
      </Box>
    </Box>
  );
};

function useZoom(): [number, typeof handlers] {
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

function Sidebar({ hasTrigger }: { hasTrigger: boolean }) {
  const integrations = useIntegrations();
  const [selectedIntegration, setSelectedIntegration] = useState<Integration>();

  const theme = useTheme();

  return (
    <Box bg="white" maxW="xs" zIndex="1" flex="1" minH="0">
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

  console.log(childNodes);

  return (
    <Grid templateColumns={`repeat(${childNodes.length || 1}, 1fr)`}>
      <GridItem
        colSpan={childNodes.length || 1}
        justifyContent="center"
        display="flex"
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
            m="8"
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
                  <LightningBoltIcon
                    height="24"
                    color={theme.colors.blue[400]}
                  />
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
      </GridItem>
      {childNodes.map((n) => (
        <GridItem colSpan={1}>
          <WorkflowItem workflow={workflow} itemId={n.id} addNode={addNode} />
        </GridItem>
      ))}
    </Grid>
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
