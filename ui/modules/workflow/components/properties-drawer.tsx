import { Button, IconButton } from "@chakra-ui/button";
import { Checkbox } from "@chakra-ui/checkbox";
import { Input } from "@chakra-ui/input";
import { Box, Divider, Flex, HStack, Stack, Text } from "@chakra-ui/layout";
import {
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
} from "@chakra-ui/number-input";
import { chakra, HTMLChakraProps, useTheme } from "@chakra-ui/system";
import { Tooltip } from "@chakra-ui/tooltip";
import { InformationCircleIcon, XIcon } from "@heroicons/react/outline";
import { HTMLMotionProps, motion } from "framer-motion";
import { Node, Param } from "../../shared/types/workflow";

type Merge<P, T> = Omit<P, keyof T> & T;
type MotionBoxProps = Merge<HTMLChakraProps<"div">, HTMLMotionProps<"div">>;
const MotionBox: React.FC<MotionBoxProps> = motion(chakra.div);

export function PropertiesDrawer({
  selectedNode,
  onClose,
  updateParamValue,
}: {
  selectedNode: Node;
  onClose: () => void;
  updateParamValue: (nodeId: string, key: string, value: string) => void;
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
      transition={{
        duration: "0.2",
        bounce: 0,
      }}
      borderLeftColor="gray.300"
      borderLeftWidth="thin"
      display="flex"
      overflow="auto"
    >
      <Flex direction="column" justifyContent="space-between" w="full" p="4">
        <Stack>
          <Flex justify="space-between" align="center" mb="4">
            <Text fontWeight="semibold">Properties</Text>
            <IconButton
              variant="ghost"
              aria-label="Close"
              icon={<XIcon height="20" color={theme.colors.gray[500]} />}
              onClick={onClose}
            />
          </Flex>
          <Stack spacing="6">
            {selectedNode?.params.map(
              (p) =>
                ({
                  string: (
                    <ParamInput param={p} nodeId={selectedNode.id} key={p.key}>
                      <Input
                        id={`${selectedNode.id}-${p.key}`}
                        value={p.value}
                        placeholder={p.input.label}
                        onChange={(e) =>
                          updateParamValue(
                            selectedNode.id,
                            p.key,
                            e.target.value
                          )
                        }
                      />
                    </ParamInput>
                  ),
                  number: (
                    <ParamInput param={p} nodeId={selectedNode.id} key={p.key}>
                      <NumberInput
                        id={`${selectedNode.id}-${p.key}`}
                        value={p.value}
                        placeholder={p.input.label}
                        onChange={(v) =>
                          updateParamValue(selectedNode.id, p.key, v)
                        }
                      >
                        <NumberInputField />
                        <NumberInputStepper>
                          <NumberIncrementStepper />
                          <NumberDecrementStepper />
                        </NumberInputStepper>
                      </NumberInput>
                    </ParamInput>
                  ),
                  boolean: (
                    <Checkbox
                      isChecked={p.value === "true"}
                      key={p.key}
                      onChange={(e) =>
                        updateParamValue(
                          selectedNode.id,
                          p.key,
                          e.target.checked.toString()
                        )
                      }
                    >
                      <Text
                        fontSize="sm"
                        fontWeight="semibold"
                        textColor="gray.600"
                      >
                        {p.input.label}
                      </Text>
                    </Checkbox>
                  ),
                }[p.input.type])
            )}
          </Stack>
        </Stack>
        <Divider my="4" />
        <Flex>
          <Button w="full" variant="outline" colorScheme="gray">
            Remove node
          </Button>
        </Flex>
        <Box minH="4" />
      </Flex>
    </MotionBox>
  );
}

function ParamInput({
  param,
  nodeId,
  children,
}: {
  param: Param;
  nodeId: string;
  children: React.ReactNode;
}) {
  const id = `${nodeId}-${param.key}`;
  const theme = useTheme();
  return (
    <Stack spacing="1">
      <HStack>
        <Text
          as="label"
          htmlFor={id}
          fontSize="sm"
          fontWeight="semibold"
          textColor="gray.600"
        >
          {param.input.label}
        </Text>
        <Tooltip
          label={param.input.description}
          fontSize="md"
          p="2"
          px="4"
          borderRadius="md"
        >
          <Box>
            <InformationCircleIcon height="16" color={theme.colors.blue[300]} />
          </Box>
        </Tooltip>
      </HStack>
      {/*  */}
      {children}
    </Stack>
  );
}
