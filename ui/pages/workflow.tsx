import { Button, ButtonGroup, IconButton } from "@chakra-ui/button";
import { Editable, EditableInput, EditablePreview } from "@chakra-ui/editable";
import { Box, Center, HStack } from "@chakra-ui/layout";
import { useTheme } from "@chakra-ui/system";
import { PencilIcon, ZoomInIcon, ZoomOutIcon } from "@heroicons/react/outline";
import { useMemo, useRef } from "react";
import { ArcherContainer } from "react-archer";
import { useDrop } from "react-dnd";
import { v4 as uuidv4 } from "uuid";
import { PropertiesDrawer } from "../modules/workflow/components/properties-drawer";
import { Sidebar } from "../modules/workflow/components/sidebar";
import { WorkflowTree } from "../modules/workflow/components/workflow-tree";
import { useNodeSelector } from "../modules/workflow/hooks/use-node-selector";
import { useWorkflowBuilder } from "../modules/workflow/hooks/use-workflow-builder";
import { useZoom } from "../modules/workflow/hooks/use-zoom";
import { Action, Trigger } from "../modules/shared/types/integration";

export default () => {
  const [workflow, { addNode }] = useWorkflowBuilder();
  const [zoom, { zoomIn, zoomOut }] = useZoom();
  const [
    selectedNodeId,
    { setSelectedNodeId, deselectNode },
  ] = useNodeSelector();

  const selectedNode = useMemo(
    () => workflow.nodes.find((n) => n.id === selectedNodeId),
    [selectedNodeId]
  );

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

  const theme = useTheme();

  return (
    <Box
      bg="gray.100"
      display="flex"
      flexDirection="column"
      h="100vh"
      w="100vw"
      overflow="hidden"
    >
      <Box
        as="header"
        bg="white"
        p="4"
        borderBottomColor="gray.300"
        borderBottomWidth="1px"
        height="4rem"
        zIndex="1"
        display="flex"
        justifyContent="space-between"
        alignItems="center"
      >
        <HStack
          spacing="4"
          borderBottomWidth="1px"
          borderBottomColor="gray.300"
        >
          <PencilIcon height="20" color={theme.colors.gray[500]} />
          <Editable defaultValue="Name your workflow">
            <EditablePreview color="gray.800" />
            <EditableInput color="grey.700" />
          </Editable>
        </HStack>
        <HStack>
          <ButtonGroup isAttached variant="solid">
            <IconButton
              icon={<ZoomOutIcon height="20" />}
              aria-label="Zoom out"
              onClick={zoomOut}
            />
            <IconButton
              icon={<ZoomInIcon height="20" />}
              aria-label="Zoom out"
              onClick={zoomIn}
            />
          </ButtonGroup>
          <Button w="full" variant="solid" colorScheme="blue">
            Save workflow
          </Button>
        </HStack>
      </Box>
      <Box display="flex" flex="1" maxH="calc(100vh - 4rem)">
        <Sidebar hasTrigger={workflow.nodes.length > 0} />
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
            width: "100%",
            height: "100%",
            margin: "auto",
            position: "relative",
          }}
        >
          <Center
            ref={drop}
            role="Dustbin"
            bg={isOver && canDrop ? "blue.50" : undefined}
            onScroll={() => archerContainerRef.current?.refreshScreen()}
            overflow="auto"
            flex="1"
            position="relative"
            width="full"
            height="full"
          >
            <Box
              position="absolute"
              top="0"
              bottom="0"
              right="0"
              left="0"
              onClick={() => deselectNode()}
            />
            <Box minW="0" minH="0" maxHeight="100%">
              <WorkflowTree
                workflow={workflow}
                addNode={addNode}
                selectNode={(id: string) => setSelectedNodeId(id)}
                selectedNodeId={selectedNodeId}
              />
            </Box>
          </Center>
          <PropertiesDrawer
            selectedNode={selectedNode}
            onClose={() => deselectNode()}
          />
        </ArcherContainer>
      </Box>
    </Box>
  );
};
