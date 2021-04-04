import { IconButton } from "@chakra-ui/button";
import { Flex, Text } from "@chakra-ui/layout";
import { chakra, HTMLChakraProps, useTheme } from "@chakra-ui/system";
import { XIcon } from "@heroicons/react/outline";
import { HTMLMotionProps, motion } from "framer-motion";
import { Node } from "../../shared/types/workflow";

type Merge<P, T> = Omit<P, keyof T> & T;
type MotionBoxProps = Merge<HTMLChakraProps<"div">, HTMLMotionProps<"div">>;
const MotionBox: React.FC<MotionBoxProps> = motion(chakra.div);

export function PropertiesDrawer({
  selectedNode,
  onClose,
}: {
  selectedNode: Node;
  onClose: () => void;
}) {
  const theme = useTheme();

  return (
    <MotionBox
      w="xs"
      position="absolute"
      top="0"
      bg="white"
      height="100%"
      variants={{
        hidden: { right: `-${theme.sizes.xs}` },
        visible: { right: "0rem" },
      }}
      animate={selectedNode ? "visible" : "hidden"}
      transitionDuration="0.5"
      p="4"
    >
      <Flex justify="space-between" align="center">
        <Text fontWeight="semibold">Properties</Text>
        <IconButton
          variant="ghost"
          aria-label="Close"
          icon={<XIcon height="20" color={theme.colors.gray[500]} />}
          onClick={onClose}
        />
      </Flex>
    </MotionBox>
  );
}
