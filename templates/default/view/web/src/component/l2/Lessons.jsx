import { Badge, Panel } from '../l1/index.js';

export function Lessons({ active = 0, items = [] }) {
  return (
    <Panel className="grid gap-4">
      {items.map((item, index) => (
        <div className="grid grid-cols-[auto_minmax(0,1fr)] gap-3" key={item.title}>
          <Badge tone={index <= active ? 'good' : 'neutral'}>{index + 1}</Badge>
          <div>
            <h3 className="m-0 text-base text-[#16211b]">{item.title}</h3>
            {item.detail ? <p className="m-0 mt-1 text-sm text-[#66786e]">{item.detail}</p> : null}
          </div>
        </div>
      ))}
    </Panel>
  );
}
