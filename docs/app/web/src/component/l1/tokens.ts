const scrollbar = "[scrollbar-width:thin] [scrollbar-color:rgb(217_119_6)_transparent] hover:[scrollbar-color:rgb(180_83_9)_transparent] dark:[scrollbar-color:rgb(250_204_21)_transparent] dark:hover:[scrollbar-color:rgb(253_224_71)_transparent]";

export const ui = {
  accent: "text-amber-950 dark:text-yellow-300",
  border: "border-amber-300/45 dark:border-yellow-300/18",
  borderStrong: "border-amber-400/55 dark:border-yellow-300/35",
  code: "bg-amber-100/80 text-amber-950 dark:bg-black dark:text-neutral-50",
  focus: "focus-visible:ring-4 focus-visible:ring-amber-300/45 dark:focus-visible:ring-yellow-300/30",
  gridLines: "bg-amber-300/30 dark:bg-yellow-300/15",
  hero: "bg-amber-50 text-neutral-950 dark:bg-neutral-950 dark:text-neutral-50",
  heroMuted: "text-amber-900/80 dark:text-neutral-400",
  muted: "text-neutral-700 dark:text-neutral-400",
  page: `bg-amber-50 text-neutral-950 font-sans dark:bg-neutral-950 dark:text-neutral-50 ${scrollbar}`,
  scrollbar,
  shadowSubtle: "shadow-sm shadow-amber-950/10 dark:shadow-black/70",
  subtle: "text-neutral-500 dark:text-neutral-500",
  surface: "bg-amber-50/88 dark:bg-neutral-950",
  surfaceQuiet: "bg-amber-100/72 dark:bg-neutral-900/82",
  surfaceSoft: "bg-amber-100/55 dark:bg-neutral-900/70",
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
    l3: "text-amber-950/78 hover:text-amber-950 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-700 dark:text-neutral-300 dark:hover:text-yellow-300 dark:focus-visible:outline-yellow-300",
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
