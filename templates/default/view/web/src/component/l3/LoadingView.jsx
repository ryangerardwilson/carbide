import { Eyebrow, Heading, ui } from '../l1/index.js';
import { cx } from '../utils.js';

export function LoadingView() {
  return (
    <main className={cx('grid min-h-svh place-items-center px-8 text-center', ui.page)}>
      <div>
        <Eyebrow>Carbide</Eyebrow>
        <Heading className="mt-2" level={1}>
          Loading app state
        </Heading>
      </div>
    </main>
  );
}
