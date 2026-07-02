import { cx } from '../utils.js';

export function RadioGroup({ label, name, onChange, options = [], value }) {
  return (
    <fieldset className="grid gap-2">
      <legend className="font-bold text-[#16211b]">{label}</legend>
      <div className="grid gap-2 sm:grid-cols-2">
        {options.map((option) => (
          <Radio
            checked={value === option.value}
            key={option.value}
            label={option.label}
            name={name}
            onChange={() => onChange?.(option.value)}
            value={option.value}
          />
        ))}
      </div>
    </fieldset>
  );
}

export function Radio({ checked, label, name, onChange, value }) {
  return (
    <label
      className={cx(
        'flex min-h-11 items-center gap-3 rounded-md border px-3 text-sm font-bold',
        checked ? 'border-teal-700 bg-teal-50 text-teal-900' : 'border-emerald-950/10 bg-white text-[#16211b]'
      )}
    >
      <input checked={checked} name={name} onChange={onChange} type="radio" value={value} />
      {label}
    </label>
  );
}
