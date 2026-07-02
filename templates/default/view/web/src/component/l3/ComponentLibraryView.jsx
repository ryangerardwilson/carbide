import { useState } from 'react';
import { Button, CodeText, Metric, Muted, Panel, ui } from '../l1/index.js';
import {
  Accordion,
  ApexChartsPanel,
  ChartJsPanel,
  ChoicesSelect,
  Combobox,
  DateRangePicker,
  Dropdown,
  FlatpickrPicker,
  FullCalendarPanel,
  GlideCarousel,
  Lessons,
  Listbox,
  Modal,
  NotificationStack,
  Popover,
  QuillEditor,
  RadioGroup,
  Select2Select,
  SimpleMDEEditor,
  Slideover,
  SplideCarousel,
  Tabs,
  Toggle,
  Tooltip,
  TrixEditor
} from '../l2/index.js';
import { cx } from '../utils.js';

const options = [
  { label: 'Frontend', value: 'frontend' },
  { label: 'Backend', value: 'backend' },
  { label: 'Database', value: 'database' }
];

const carouselItems = [
  { title: 'Shell', detail: 'Bun serves the React app and same-origin proxy.' },
  { title: 'API', detail: 'Go handles auth, sessions, validation, and JSON.' },
  { title: 'Data', detail: 'Postgres keeps durable app state.' }
];

export function ComponentLibraryView() {
  const [modalOpen, setModalOpen] = useState(false);
  const [panelOpen, setPanelOpen] = useState(false);
  const [toggle, setToggle] = useState(true);
  const [choice, setChoice] = useState('frontend');
  const [dates, setDates] = useState({ start: '', end: '' });

  return (
    <div className="grid gap-6">
      <div className="grid gap-4 lg:grid-cols-3">
        <Panel className="lg:col-span-2">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 className={cx('m-0 text-2xl', ui.text)}>Component library</h2>
              <Muted className="mt-1">
                L3 screens compose L2 patterns. L2 patterns compose L1 primitives.
              </Muted>
            </div>
            <Dropdown
              align="right"
              items={[
                { label: 'Open modal', value: 'modal' },
                { label: 'Open slideover', value: 'slideover' }
              ]}
              label="Actions"
              onSelect={(item) => {
                if (item.value === 'modal') setModalOpen(true);
                if (item.value === 'slideover') setPanelOpen(true);
              }}
            />
          </div>
          <div className="mt-5 grid gap-4 sm:grid-cols-3">
            <Metric label="Primitive layer" value="L1" detail="controls, text, surface" />
            <Metric label="Pattern layer" value="L2" detail="dialogs, tabs, inputs" />
            <Metric label="Product layer" value="L3" detail="starter app views" />
          </div>
        </Panel>
        <NotificationStack
          items={[
            { title: 'Auth ready', detail: 'Same-origin API proxy is active.', tone: 'success' },
            { title: 'UI contract', detail: 'Component boundaries are scaffolded.', tone: 'info' }
          ]}
        />
      </div>

      <Tabs
        tabs={[
          {
            label: 'Patterns',
            value: 'patterns',
            content: (
              <div className="grid gap-4 lg:grid-cols-2">
                <Accordion
                  items={[
                    { title: 'Disclosure', content: 'Accordion and Disclosure share one L2 pattern.' },
                    { title: 'Dialog', content: 'Modal and Slideover own temporary focus surfaces.' }
                  ]}
                />
                <Panel className="grid gap-4">
                  <Toggle checked={toggle} label="Switch / Toggle" onChange={setToggle} />
                  <RadioGroup label="Radio Group" name="layer" onChange={setChoice} options={options} value={choice} />
                  <Tooltip text="Tooltips are L2 because they own interaction state.">
                    <Button variant="secondary">Tooltip target</Button>
                  </Tooltip>
                  <Popover label="Popover">
                    <Muted>
                      Popovers keep short contextual controls close to the trigger.
                    </Muted>
                  </Popover>
                </Panel>
                <Lessons
                  active={1}
                  items={[
                    { title: 'Create the app', detail: 'Use carbide new or carbide init.' },
                    { title: 'Run the stack', detail: 'Use carbide run dev.' },
                    { title: 'Build product UI', detail: 'Start in L3 and move reusable patterns to L2.' }
                  ]}
                />
                <GlideCarousel items={carouselItems} />
              </div>
            )
          },
          {
            label: 'Inputs',
            value: 'inputs',
            content: (
              <div className="grid gap-4 lg:grid-cols-2">
                <Listbox label="Listbox / Select" onChange={setChoice} options={options} value={choice} />
                <Combobox label="Combobox / Autocomplete" options={options} />
                <Select2Select label="Select2" onChange={setChoice} options={options} value={choice} />
                <ChoicesSelect label="Choices.js" onChange={setChoice} options={options} value={choice} />
                <DateRangePicker end={dates.end} onChange={setDates} start={dates.start} />
                <FlatpickrPicker label="Flatpickr" />
              </div>
            )
          },
          {
            label: 'Integrations',
            value: 'integrations',
            content: (
              <div className="grid gap-4 xl:grid-cols-2">
                <TrixEditor label="Trix" placeholder="Write rich text..." />
                <QuillEditor label="Quill" placeholder="Write rich text..." />
                <SimpleMDEEditor label="SimpleMDE" placeholder="Write Markdown..." />
                <ChartJsPanel title="Chart.js" />
                <ApexChartsPanel title="ApexCharts" />
                <FullCalendarPanel events={[{ title: 'Launch' }]} />
                <SplideCarousel items={carouselItems} />
              </div>
            )
          }
        ]}
      />

      <Modal
        description="Dialog (Modal) is an L2 interaction pattern."
        onClose={() => setModalOpen(false)}
        open={modalOpen}
        title="Dialog"
      >
        <Muted>
          A generated Carbide app can use <CodeText>Modal</CodeText> without adding a second UI kit.
        </Muted>
      </Modal>

      <Slideover onClose={() => setPanelOpen(false)} open={panelOpen} title="Slideover">
        <Muted>
          Slideovers use the same L2 boundary as modals and stay independent from app-specific data.
        </Muted>
      </Slideover>
    </div>
  );
}
