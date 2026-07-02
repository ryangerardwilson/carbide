import { Field, ui } from '../l1/index.js';
import { cx } from '../utils.js';

export function Listbox({ label, onChange, options = [], value }) {
  return (
    <Field label={label}>
      <select
        className={cx('min-h-12 w-full rounded-md border px-3 py-2 outline-none', ui.input)}
        onChange={(event) => onChange?.(event.target.value)}
        value={value}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </Field>
  );
}
