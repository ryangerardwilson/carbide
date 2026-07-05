import { Eyebrow, Heading, ui } from '../l1';
import { cx } from '../../lib/cx';

const loadingClassLayers = {
  shell: {
    l1: 'grid place-items-center',
    l2: 'min-h-svh px-4 text-center',
    l3: ui.page
  }
};

export function LoadingView() {
  return (
    <main className={cx(loadingClassLayers.shell.l1, loadingClassLayers.shell.l2, loadingClassLayers.shell.l3)}>
      <div>
        <Eyebrow>Carbide</Eyebrow>
        <Heading className="mt-2" level={1}>
          Loading app state
        </Heading>
      </div>
    </main>
  );
}
