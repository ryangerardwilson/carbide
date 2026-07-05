import { docsClassLayers } from "../l1/index.js";

export const docsChromeClassLayers = {
  topbar: {
    l1: "sticky top-0 z-30 flex items-center",
    l2: "min-h-16 border-b px-5",
    l3: "border-neutral-800 bg-black/90 text-neutral-100 backdrop-blur",
  },
  sidebar: {
    l1: "sticky top-16 self-start",
    l2: "max-h-[calc(100vh-4rem)] border-r px-5 py-6",
    l3: "border-neutral-800 text-neutral-400",
  },
  content: {
    l1: "min-w-0",
    l2: "px-5 py-8",
    l3: docsClassLayers.page.l3,
  },
};

export function docsStaticHeaders({ cache, type, contract }) {
  return {
    "cache-control": cache,
    "content-type": type,
    "x-carbide-component-contract": contract,
  };
}
