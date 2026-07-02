import { Badge, Panel, ui } from '../l1/index.js';
import { cx } from '../utils.js';

const toneMap = {
  info: 'neutral',
  success: 'good',
  warning: 'warn',
  error: 'danger'
};

export function Notifications({ items = [] }) {
  return (
    <div className="grid gap-2" role="status">
      {items.map((item) => (
        <Panel className="flex items-start justify-between gap-4 p-4" key={item.id || item.title}>
          <div>
            <h3 className={cx('m-0 text-base', ui.text)}>{item.title}</h3>
            {item.detail ? <p className={cx('m-0 mt-1 text-sm', ui.subtle)}>{item.detail}</p> : null}
          </div>
          <Badge tone={toneMap[item.tone] || 'neutral'}>{item.tone || 'info'}</Badge>
        </Panel>
      ))}
    </div>
  );
}

export const NotificationStack = Notifications;
