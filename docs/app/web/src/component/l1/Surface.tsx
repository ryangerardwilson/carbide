import type { ReactNode } from "react";
import { cx } from "../../lib/cx";
import { docsClassLayers } from "./tokens";

interface SurfaceProps {
  children: ReactNode;
  className?: string;
}

export function Surface({ children, className = "" }: SurfaceProps) {
  return <section className={cx(docsClassLayers.surface.l1, docsClassLayers.surface.l2, docsClassLayers.surface.l3, className)}>{children}</section>;
}
