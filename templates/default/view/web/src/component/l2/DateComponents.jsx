import { Field, Panel, ui } from '../l1/index.js';
import { cx } from '../utils.js';

export function FlatpickrPicker({ label = 'Date', onChange, value = '' }) {
  return (
    <Field hint="Flatpickr adapter surface" label={label}>
      <input
        className={cx('min-h-12 w-full rounded-md border px-3 py-2 outline-none', ui.input)}
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
    <Panel className="p-4" data-integration="FullCalendar">
      <h3 className={cx('m-0 text-base', ui.text)}>FullCalendar</h3>
      <div className={cx('mt-3 grid grid-cols-7 gap-1 text-center text-xs', ui.subtle)}>
        {Array.from({ length: 14 }).map((_, index) => (
          <span className={cx('min-h-9 rounded p-2', ui.surfaceQuiet)} key={index}>
            {index + 1}
          </span>
        ))}
      </div>
      {events.length ? <p className={cx('m-0 mt-3 text-sm', ui.subtle)}>{events.length} events loaded</p> : null}
    </Panel>
  );
}
