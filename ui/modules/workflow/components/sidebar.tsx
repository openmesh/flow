import { Input, InputGroup, InputLeftElement } from "@chakra-ui/input";
import { Box, Stack } from "@chakra-ui/layout";
import { useTheme } from "@chakra-ui/system";
import { Tab, TabList, TabPanel, TabPanels, Tabs } from "@chakra-ui/tabs";
import { FilterIcon } from "@heroicons/react/outline";
import { useState } from "react";
import { DropdownCombobox } from "../../shared/components/dropdown-combobox";
import { useIntegrations } from "../../shared/hooks/use-integrations";
import { Integration } from "../../shared/types/integration";
import { SidebarItem } from "./sidebar-item";

export function Sidebar({ hasTrigger }: { hasTrigger: boolean }) {
  const integrations = useIntegrations();
  const [selectedIntegration, setSelectedIntegration] = useState<Integration>();

  const theme = useTheme();

  return (
    <Box
      bg="white"
      minW="xs"
      zIndex="1"
      flex="1"
      display="flex"
      flexDirection="column"
      borderRightColor="gray.300"
      borderRightWidth="1px"
    >
      <Stack px="4" py="6" spacing="4" flex="1">
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
                  <SidebarItem
                    item={trigger}
                    type="TRIGGER"
                    key={trigger.key}
                  />
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
                  <SidebarItem item={action} type="ACTION" key={action.key} />
                ))}
              </Stack>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Stack>
    </Box>
  );
}
