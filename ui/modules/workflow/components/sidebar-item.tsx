import { Box, Stack, Text } from "@chakra-ui/layout";
import { useTheme } from "@chakra-ui/system";
import { LightningBoltIcon } from "@heroicons/react/solid";
import { useDrag } from "react-dnd";
import { Action, Trigger } from "../../shared/types/integration";

export function SidebarItem({
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
