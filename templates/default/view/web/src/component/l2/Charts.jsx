import { Panel, ui } from '../l1/index.js';
import { cx } from '../utils.js';

function ChartAdapter({ library, points = [18, 42, 31, 58, 46, 69], title = library }) {
  const max = Math.max(...points, 1);

  return (
    <Panel data-integration={library}>
      <h3 className={cx('m-0 text-base', ui.text)}>{title}</h3>
      <div className="mt-4 flex h-32 items-end gap-2" aria-label={`${library} chart preview`}>
        {points.map((point, index) => (
          <div className="flex flex-1 items-end" key={`${point}-${index}`}>
            <span
              className="cb-chart-bar block w-full rounded-t"
              style={{ height: `${Math.max(12, (point / max) * 100)}%` }}
            />
          </div>
        ))}
      </div>
    </Panel>
  );
}

export function ChartJsPanel(props) {
  return <ChartAdapter library="Chart.js" {...props} />;
}

export function ApexChartsPanel(props) {
  return <ChartAdapter library="ApexCharts" {...props} />;
}
