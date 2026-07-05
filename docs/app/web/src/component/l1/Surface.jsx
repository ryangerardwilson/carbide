import { cx } from "../../lib/cx.js";
import { docsClassLayers } from "./tokens.js";

export function Surface({ children, className = "" }) {
  return <section className={cx(docsClassLayers.surface.l1, docsClassLayers.surface.l2, docsClassLayers.surface.l3, className)}>{children}</section>;
}
