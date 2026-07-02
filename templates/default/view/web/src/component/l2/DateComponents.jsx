import { Field } from '../l1/index.js';

export function FlatpickrPicker({ label = 'Date', onChange, value = '' }) {
  return (
    <Field hint="Flatpickr adapter surface" label={label}>
      <input
        className="min-h-12 w-full rounded-md border border-emerald-900/20 bg-white px-3 py-2 text-[#16211b] outline-none focus:border-teal-700 focus:ring-4 focus:ring-teal-700/15"
        data-integration="Flatpickr"
        onChange={(event) => onChange?.(event.target.value)}
        type="date"
        value={value}
      />
    </Field>
  );
}

export function DateRangePicker({ end = '', onChange, start = '' }) {
  return (
    <div className="grid gap-3 sm:grid-cols-2" data-integration="Date Range Picker">
      <FlatpickrPicker label="Start" onChange={(next) => onChange?.({ start: next, end })} value={start} />
      <FlatpickrPicker label="End" onChange={(next) => onChange?.({ start, end: next })} value={end} />
    </div>
  );
}

export function FullCalendarPanel({ events = [] }) {
  return (
    <section className="rounded-lg border border-emerald-950/10 bg-white p-4" data-integration="FullCalendar">
      <h3 className="m-0 text-base text-[#16211b]">FullCalendar</h3>
      <div className="mt-3 grid grid-cols-7 gap-1 text-center text-xs text-[#66786e]">
        {Array.from({ length: 14 }).map((_, index) => (
          <span className="min-h-9 rounded bg-emerald-50 p-2" key={index}>
            {index + 1}
          </span>
        ))}
      </div>
      {events.length ? <p className="m-0 mt-3 text-sm text-[#66786e]">{events.length} events loaded</p> : null}
    </section>
  );
}
