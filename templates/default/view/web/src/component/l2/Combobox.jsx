import { useMemo, useState } from 'react';
import { Field, TextInput } from '../l1/index.js';

export function Combobox({ label, onChange, options = [], value = '' }) {
  const [query, setQuery] = useState(value);
  const matches = useMemo(
    () => options.filter((option) => option.label.toLowerCase().includes(query.toLowerCase())).slice(0, 5),
    [options, query]
  );

  return (
    <div className="relative">
      <Field label={label}>
        <TextInput
          onChange={(event) => {
            setQuery(event.target.value);
            onChange?.(event.target.value);
          }}
          value={query}
        />
      </Field>
      {query ? (
        <div className="absolute z-20 mt-2 w-full rounded-lg border border-emerald-950/10 bg-white p-1 shadow-lg">
          {matches.map((option) => (
            <button
              className="block min-h-10 w-full rounded-md px-3 text-left text-sm hover:bg-emerald-50"
              key={option.value}
              type="button"
              onClick={() => {
                setQuery(option.label);
                onChange?.(option.value);
              }}
            >
              {option.label}
            </button>
          ))}
        </div>
      ) : null}
    </div>
  );
}
