import { Button, IconButton } from "@chakra-ui/button";
import { Box, Center, Flex, Text } from "@chakra-ui/layout";
import { PlusOutline } from "@graywolfai/react-heroicons";

export default () => {
  return (
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
      <Box display="flex" h="full" flex="1" position="relative">
        <Sidebar />
        <Center flex="1">
          <IconButton
            position="absolute"
            top="1rem"
            right="1rem"
            aria-label="Add node"
            icon={<PlusOutline />}
            isRound
            size="lg"
            colorScheme="blue"
          />

          <Workflow />
        </Center>
      </Box>
    </Box>
  );
};

function Sidebar() {
  return (
    <Box bg="white" flex="1" maxW="xs" p="4">
      <Text>stuffs</Text>
    </Box>
  );
}

function Workflow() {
  return (
    <Box bg="white" shadow="sm" maxW="container.sm" flex="1" borderRadius="md">
      <Box
        bg="blue.400"
        p="4"
        borderRadius="md"
        display="flex"
        justifyContent="space-between"
        alignItems="center"
      >
        <Box>
          <Text color="white" fontSize="xl" mb="1">
            Trigger
          </Text>
          <Text color="blue.50">
            A trigger is the event that starts a workflow
          </Text>
        </Box>
        <Button colorScheme="blue" variant="solid">
          Learn more
        </Button>
      </Box>
    </Box>
  );
}
