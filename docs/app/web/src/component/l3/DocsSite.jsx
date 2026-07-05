import { docsStaticHeaders } from "../l2/index.js";

export const docsSiteClassLayers = {
  shell: {
    l1: "grid min-h-screen",
    l2: "grid-cols-1 lg:grid-cols-[270px_minmax(0,1fr)_224px]",
    l3: "bg-neutral-950 text-neutral-100",
  },
  article: {
    l1: "min-w-0",
    l2: "max-w-3xl",
    l3: "text-neutral-200",
  },
};

export function docsWebContract() {
  return {
    id: "docs-web:l1-l2-l3-tailwind",
    levels: ["component/l1", "component/l2", "component/l3"],
    styling: "tailwind",
  };
}

export function docsResponseHeaders(options) {
  return docsStaticHeaders({
    ...options,
    contract: docsWebContract().id,
  });
}
