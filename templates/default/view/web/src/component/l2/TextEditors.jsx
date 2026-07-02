import { Field } from '../l1/index.js';

function TextEditorAdapter({ library, label, onChange, placeholder = '', value = '' }) {
  return (
    <Field hint={`${library} adapter surface`} label={label || library}>
      <textarea
        className="min-h-32 w-full rounded-md border border-emerald-900/20 bg-white px-3 py-2 text-[#16211b] outline-none focus:border-teal-700 focus:ring-4 focus:ring-teal-700/15"
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
