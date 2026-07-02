import { Field } from '../l1/index.js';

export function Listbox({ label, onChange, options = [], value }) {
  return (
    <Field label={label}>
      <select
        className="min-h-12 w-full rounded-md border border-emerald-900/20 bg-white px-3 py-2 text-[#16211b] outline-none focus:border-teal-700 focus:ring-4 focus:ring-teal-700/15"
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
