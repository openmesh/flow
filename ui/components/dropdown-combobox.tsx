import { Input, InputGroup, InputLeftElement } from "@chakra-ui/input";
import { Box, List, ListItem } from "@chakra-ui/layout";
import { useTheme } from "@chakra-ui/system";
import { SearchIcon } from "@heroicons/react/outline";
import { useCombobox } from "downshift";
import { useState } from "react";

// integrations
// GitHub

const items = [
  {
    label: "GitHub",
    ref: "GITHUB",
    actions: [
      {
        name: "Update a repository",
        description: "Updates information about a repository",
        ref: "REPOSITORY_UPDATE",
      },
      {
        name: "Create a repository",
        description: "Create a new repository",
        ref: "REPOSITORY_CREATE",
      },
      {
        name: "Delete a repository",
        description: "Deletes a repository",
        ref: "REPOSITORY_DELETE",
      },
    ],
  },
];

export function DropdownCombobox() {
  const [inputItems, setInputItems] = useState(items);
  const {
    isOpen,
    getToggleButtonProps,
    getLabelProps,
    getMenuProps,
    getInputProps,
    getComboboxProps,
    highlightedIndex,
    getItemProps,
  } = useCombobox({
    items: inputItems,
    itemToString: (item) => item.label,
    onInputValueChange: ({ inputValue }) => {
      setInputItems(
        items.filter((item) =>
          item.label.toLowerCase().startsWith(inputValue.toLowerCase())
        )
      );
    },
  });

  const theme = useTheme();

  return (
    <Box>
      {/* <label {...getLabelProps()}>Choose an element:</label> */}
      <div {...getComboboxProps()}>
        <InputGroup w="full">
          <InputLeftElement
            pointerEvents="none"
            children={<SearchIcon height="24" color={theme.colors.gray[300]} />}
          />
          <Input placeholder="Search applications" {...getInputProps()} />
        </InputGroup>
        {/* <button
          type="button"
          {...getToggleButtonProps()}
          aria-label="toggle menu"
        >
          &#8595;
        </button> */}
      </div>
      <List {...getMenuProps()} mt="2" borderRadius="md" shadow="md">
        {isOpen &&
          inputItems.map((item, index) => (
            <ListItem
              style={
                highlightedIndex === index ? { backgroundColor: "#bde4ff" } : {}
              }
              borderTopRadius={index === 0 ? "md" : ""}
              borderBottomRadius={index === inputItems.length - 1 ? "md" : ""}
              key={`${item}${index}`}
              {...getItemProps({ item, index })}
              p="2"
            >
              {item.label}
            </ListItem>
          ))}
      </List>
    </Box>
  );
}
