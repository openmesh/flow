import { useMemo } from "react";
import { useDrop } from "react-dnd";
import { Action, Trigger } from "../../shared/types/integration";
import { Node, Workflow } from "../../shared/types/workflow";
import { v4 as uuidv4 } from "uuid";
import { Box, Grid, GridItem, Stack, Text } from "@chakra-ui/layout";
import { ArcherElement } from "react-archer";
import { useTheme } from "@chakra-ui/system";
import { LightningBoltIcon } from "@heroicons/react/solid";

export function WorkflowItem({
  workflow,
  itemId,
  addNode,
  selectNode,
  selectedNodeId,
}: {
  workflow: Workflow;
  itemId: string;
  addNode: (node: Node) => void;
  selectNode: (id: string) => void;
  selectedNodeId: string;
}) {
  const theme = useTheme();

  const node = useMemo(() => workflow.nodes.find((x) => x.id === itemId), [
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
    () => workflow.nodes.filter((n) => node.childrenIds.includes(n.id)),
    [workflow]
  );

  return (
    <Grid templateColumns={`repeat(${childNodes.length || 1}, 1fr)`}>
      <GridItem
        colSpan={childNodes.length || 1}
        justifyContent="center"
        display="flex"
      >
        <ArcherElement
          id={`workflow-item-${node.id}`}
          relations={node.childrenIds.map((childId) => ({
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
            shadow={selectedNodeId === itemId ? "md" : "sm"}
            borderColor={selectedNodeId === itemId ? "blue.500" : "transparent"}
            borderWidth="thin"
            w="lg"
            borderRadius="md"
            _hover={{ shadow: "lg" }}
            position="relative"
            m="8"
            role="button"
            onClick={() => selectNode(itemId)}
            transition="all"
            transitionDuration="200ms"
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
              id={`${node.id}-child-drop-area`}
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
                    {node.label}
                  </Text>
                  <Text color="gray.500">{node.description}</Text>
                </Stack>
              </Stack>
            </Box>
          </Box>
        </ArcherElement>
      </GridItem>
      {childNodes.map((n) => (
        <GridItem colSpan={1} key={n.id}>
          <WorkflowItem
            workflow={workflow}
            itemId={n.id}
            addNode={addNode}
            selectNode={selectNode}
            selectedNodeId={selectedNodeId}
          />
        </GridItem>
      ))}
    </Grid>
  );
}
