const scrollbar = "[scrollbar-width:thin] [scrollbar-color:rgb(163_163_163)_transparent] hover:[scrollbar-color:rgb(115_115_115)_transparent] dark:[scrollbar-color:rgb(82_82_82)_transparent] dark:hover:[scrollbar-color:rgb(115_115_115)_transparent]";

export const ui = {
  accent: "text-neutral-950 dark:text-neutral-50",
  border: "border-neutral-200 dark:border-neutral-800",
  borderStrong: "border-neutral-300 dark:border-neutral-700",
  code: "bg-neutral-100 text-neutral-950 dark:bg-neutral-900 dark:text-neutral-50",
  focus: "focus-visible:ring-4 focus-visible:ring-neutral-300 dark:focus-visible:ring-neutral-700",
  gridLines: "bg-neutral-200 dark:bg-neutral-800",
  hero: "bg-white text-neutral-950 dark:bg-black dark:text-neutral-50",
  heroMuted: "text-neutral-600 dark:text-neutral-400",
  muted: "text-neutral-600 dark:text-neutral-400",
  page: `bg-white text-neutral-950 font-sans dark:bg-black dark:text-neutral-50 ${scrollbar}`,
  scrollbar,
  shadowSubtle: "shadow-sm shadow-neutral-950/5 dark:shadow-black/70",
  subtle: "text-neutral-500 dark:text-neutral-500",
  surface: "bg-white dark:bg-neutral-950",
  surfaceQuiet: "bg-neutral-100 dark:bg-neutral-900",
  surfaceSoft: "bg-neutral-50 dark:bg-neutral-950",
  text: "text-neutral-950 dark:text-neutral-50",
};

export const docsClassLayers = {
  page: {
    l1: "min-h-screen",
    l2: "text-sm leading-6",
    l3: ui.page,
  },
  scrollbar: {
    l1: ui.scrollbar,
    l2: "",
    l3: "",
  },
  link: {
    l1: "inline-flex items-center",
    l2: "text-sm underline-offset-4",
    l3: "text-neutral-700 hover:text-neutral-950 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-neutral-950 dark:text-neutral-300 dark:hover:text-neutral-50 dark:focus-visible:outline-neutral-50",
  },
  surface: {
    l1: "relative overflow-hidden",
    l2: "rounded-lg border p-4",
    l3: `${ui.border} ${ui.surface} ${ui.text}`,
  },
  code: {
    l1: "inline-block",
    l2: "rounded-md px-1.5 py-0.5 font-mono text-xs",
    l3: `border ${ui.border} ${ui.code}`,
  },
};
