import { docsClassLayers, ui } from "../l1";
import { docsStaticHeaders, rewriteDocsClasses } from "../l2";

export const docsSiteClassLayers = {
  shell: {
    l1: "grid min-h-screen",
    l2: "grid-cols-1 lg:grid-cols-[270px_minmax(0,1fr)_224px]",
    l3: `${docsClassLayers.page.l3} [scrollbar-width:thin]`,
  },
  article: {
    l1: "min-w-0",
    l2: "max-w-3xl",
    l3: ui.muted,
  },
};

interface DocsResponseOptions {
  cache: string;
  type: string;
}

export function docsWebContract() {
  return {
    id: "docs-web:l1-l2-l3-tailwind",
    levels: ["component/l1", "component/l2", "component/l3"],
    styling: "tailwind",
  };
}

export function docsResponseHeaders(options: DocsResponseOptions): Record<string, string> {
  return docsStaticHeaders({
    ...options,
    contract: docsWebContract().id,
  });
}

export function rewriteDocsHtml(html: string): string {
  return rewriteDocsClasses(html);
}
