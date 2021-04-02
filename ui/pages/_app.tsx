import { ChakraProvider } from "@chakra-ui/react";
import { useEffect, useRef } from "react";
import { ArcherContainer } from "react-archer";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";

function MyApp({ Component, pageProps }) {
  const archerContainerRef = useRef<ArcherContainer>();

  function onScroll(e) {
    console.log(e);
    console.log("refreshing screen");
    archerContainerRef.current?.refreshScreen();
  }

  useEffect(() => {
    window.addEventListener("scroll", onScroll);
    // return window.removeEventListener("scroll", onScroll);
  }, []);

  return (
    <DndProvider backend={HTML5Backend}>
      <ChakraProvider>
        <Component {...pageProps} />
      </ChakraProvider>
    </DndProvider>
  );
}

export default MyApp;
