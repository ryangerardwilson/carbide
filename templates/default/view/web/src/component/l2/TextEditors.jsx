import { Field, ui } from '../l1/index.js';
import { cx } from '../utils.js';

function TextEditorAdapter({ library, label, onChange, placeholder = '', value = '' }) {
  return (
    <Field hint={`${library} adapter surface`} label={label || library}>
      <textarea
        className={cx('min-h-32 w-full rounded-md border px-3 py-2 outline-none', ui.input)}
        data-integration={library}
        onChange={(event) => onChange?.(event.target.value)}
        placeholder={placeholder}
        value={value}
      />
    </Field>
  );
}

export function TrixEditor(props) {
  return <TextEditorAdapter library="Trix" {...props} />;
}

export function QuillEditor(props) {
  return <TextEditorAdapter library="Quill" {...props} />;
}

export function SimpleMDEEditor(props) {
  return <TextEditorAdapter library="SimpleMDE" {...props} />;
}
