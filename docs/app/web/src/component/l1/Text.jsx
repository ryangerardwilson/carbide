import { cx } from "../../lib/cx.js";
import { docsClassLayers } from "./tokens.js";

export function DocsLink({ children, className = "", href }) {
  return (
    <a className={cx(docsClassLayers.link.l1, docsClassLayers.link.l2, docsClassLayers.link.l3, className)} href={href}>
      {children}
    </a>
  );
}

export function InlineCode({ children, className = "" }) {
  return <code className={cx(docsClassLayers.code.l1, docsClassLayers.code.l2, docsClassLayers.code.l3, className)}>{children}</code>;
}
