import type { ReactNode } from "react";
import { cx } from "../../lib/cx";
import { docsClassLayers } from "./tokens";

interface DocsLinkProps {
  children: ReactNode;
  className?: string;
  href: string;
}

interface InlineCodeProps {
  children: ReactNode;
  className?: string;
}

export function DocsLink({ children, className = "", href }: DocsLinkProps) {
  return (
    <a className={cx(docsClassLayers.link.l1, docsClassLayers.link.l2, docsClassLayers.link.l3, className)} href={href}>
      {children}
    </a>
  );
}

export function InlineCode({ children, className = "" }: InlineCodeProps) {
  return <code className={cx(docsClassLayers.code.l1, docsClassLayers.code.l2, docsClassLayers.code.l3, className)}>{children}</code>;
}
