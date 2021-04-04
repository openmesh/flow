import { ChakraProvider, extendTheme } from "@chakra-ui/react";
import { useEffect, useRef } from "react";
import { ArcherContainer } from "react-archer";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";

const theme = extendTheme({
  colors: {
    blue: {
      "50": "#f0f5ff",
      "100": "#d6e4ff",
      "200": "#adc6ff",
      "300": "#85a5ff",
      "400": "#597ef7",
      "500": "#2f54eb",
      "600": "#1d39c4",
      "700": "#10239e",
      "800": "#061178",
      "900": "#030852",
    },
  },
});

function MyApp({ Component, pageProps }) {
  const archerContainerRef = useRef<ArcherContainer>();

  function onScroll() {
    archerContainerRef.current?.refreshScreen();
  }

  useEffect(() => {
    window.addEventListener("scroll", onScroll);
  }, []);

  return (
    <DndProvider backend={HTML5Backend}>
      <ChakraProvider theme={theme}>
        <Component {...pageProps} />
      </ChakraProvider>
    </DndProvider>
  );
}

export default MyApp;
