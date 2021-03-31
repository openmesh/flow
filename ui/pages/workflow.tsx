import { Button, IconButton } from "@chakra-ui/button";
import { Input, InputGroup, InputLeftElement } from "@chakra-ui/input";
import { Box, Center, Flex, Stack, Text } from "@chakra-ui/layout";
import { Select } from "@chakra-ui/select";
import { useTheme } from "@chakra-ui/system";
import { Tab, TabList, TabPanel, TabPanels, Tabs } from "@chakra-ui/tabs";
import { PlusIcon, SearchIcon } from "@heroicons/react/outline";
import { DropdownCombobox } from "../components/dropdown-combobox";

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
            icon={<PlusIcon />}
            isRound
            size="lg"
            colorScheme="blue"
          />
          <Box>
            <Workflow />
            <Button>Add node</Button>
          </Box>
        </Center>
      </Box>
    </Box>
  );
};

function Sidebar() {
  const theme = useTheme();
  return (
    <Box bg="white" flex="1" maxW="xs">
      <Stack px="4" py="6" spacing="4">
        {/* <InputGroup>
          <InputLeftElement
            pointerEvents="none"
            children={<SearchIcon height="24" color={theme.colors.gray[300]} />}
          />
          <Input placeholder="Search nodes" />
        </InputGroup> */}
        <Tabs isFitted>
          <TabList>
            <Tab>Triggers</Tab>
            <Tab>Actions</Tab>
          </TabList>
          <TabPanels>
            <TabPanel px="0">
              {/* <InputGroup w="full">
                <InputLeftElement
                  pointerEvents="none"
                  children={
                    <SearchIcon height="24" color={theme.colors.gray[300]} />
                  }
                />
                <Input placeholder="Search applications" />
              </InputGroup> */}
              <DropdownCombobox />
              {/* <Select placeholder="Integration"></Select> */}
            </TabPanel>
            <TabPanel>
              <p>two!</p>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Stack>
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
