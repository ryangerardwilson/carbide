import { Listbox } from './Listbox.jsx';

function EnhancedSelectAdapter({ library, ...props }) {
  return (
    <div data-integration={library}>
      <Listbox {...props} />
    </div>
  );
}

export function Select2Select(props) {
  return <EnhancedSelectAdapter library="Select2" {...props} />;
}

export function ChoicesSelect(props) {
  return <EnhancedSelectAdapter library="Choices.js" {...props} />;
}
