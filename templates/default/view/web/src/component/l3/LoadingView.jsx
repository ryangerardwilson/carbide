import { Eyebrow, Heading } from '../l1/index.js';

export function LoadingView() {
  return (
    <main className="grid min-h-svh place-items-center px-8 text-center">
      <div>
        <Eyebrow>Carbide</Eyebrow>
        <Heading className="mt-2" level={1}>
          Loading app state
        </Heading>
      </div>
    </main>
  );
}
