<script setup lang="ts">
import { computed, onScopeDispose, ref } from "vue";
import type { VisitorDay } from "@/lib/visitors";

const props = defineProps<{ days: VisitorDay[] }>();
const bucketDuration = 10 * 60 * 1000;
const currentBucketStart = ref(0);
let boundaryTimer: ReturnType<typeof setTimeout>;
function updateCurrentBucket() {
  const now = Date.now();
  currentBucketStart.value = Math.floor(now / bucketDuration) * bucketDuration;
  boundaryTimer = setTimeout(updateCurrentBucket, bucketDuration - now % bucketDuration);
}
updateCurrentBucket();
onScopeDispose(() => clearTimeout(boundaryTimer));
const colors = ["#2864a6", "#c04a25", "#27804c", "#8248ad", "#987100", "#147b87"];
function minuteOfDay(start: number) {
  const date = new Date(start);
  return (date.getUTCHours() * 60 + date.getUTCMinutes() + 540) % 1440;
}
function timeLabel(minute: number) {
  return `${String(Math.floor(minute / 60)).padStart(2, "0")}:${String(minute % 60).padStart(2, "0")}`;
}
const minutes = computed(() => props.days.flatMap((day) => day.buckets.map((b) => minuteOfDay(b.start))));
const first = computed(() => Math.min(...minutes.value));
const last = computed(() => Math.max(...minutes.value));
const maxCount = computed(() => Math.ceil(Math.max(...props.days.flatMap((day) => day.buckets.map((b) => b.count))) / 4) * 4);
const yTicks = computed(() => Array.from({ length: 5 }, (_, i) => i * maxCount.value / 4));
const viewHeight = computed(() => Math.max(280, 42 + props.days.length * 18));
function xPosition(minute: number) {
  return first.value === last.value ? 410 : 60 + (minute - first.value) * 700 / (last.value - first.value);
}
const series = computed(() => props.days.map((day, i) => ({
  date: day.date,
  pending: day.buckets[day.buckets.length - 1].start === currentBucketStart.value,
  color: colors[i % colors.length],
  dash: i < colors.length ? undefined : `${8 - (Math.floor(i / colors.length) % 3) * 2} 4`,
  points: day.buckets.map((bucket) => ({
    ...bucket,
    minute: minuteOfDay(bucket.start),
    x: xPosition(minuteOfDay(bucket.start)),
    y: 220 - bucket.count * 180 / maxCount.value,
  })),
})));
const times = computed(() => Array.from({ length: (last.value - first.value) / 10 + 1 }, (_, i) => ({
  minute: first.value + i * 10,
  x: xPosition(first.value + i * 10),
})));
const xTicks = computed(() => times.value.filter((_, i, all) =>
  i % Math.max(1, Math.ceil((all.length - 1) / 6)) === 0 || i === all.length - 1,
).filter((point, i, ticks) => i === ticks.length - 1 || ticks[i + 1].x - point.x >= 70));
const hoveredMinute = ref<number | null>(null);
const hovered = computed(() => times.value.find((point) => point.minute === hoveredMinute.value));
const hoveredRows = computed(() => series.value.flatMap((day) => {
  const point = day.points.find((p) => p.minute === hoveredMinute.value);
  return point ? [{ date: day.date, color: day.color, count: point.count, x: point.x, y: point.y }] : [];
}));
const tooltipWidth = 170;
const tooltipHeight = computed(() => 22 + hoveredRows.value.length * 18);
const cursor = ref<{ x: number; y: number } | null>(null);
const tooltipTransform = computed(() => {
  const position = cursor.value ?? (hovered.value ? { x: hovered.value.x, y: 40 } : null);
  if (!position) return "";
  const x = position.x + tooltipWidth + 12 <= 790 ? position.x + 12 : position.x - tooltipWidth - 12;
  const y = position.y + tooltipHeight.value + 12 <= viewHeight.value - 10
    ? position.y + 12 : position.y - tooltipHeight.value - 12;
  return `translate(${Math.max(10, Math.min(790 - tooltipWidth, x))}, ${Math.max(10, Math.min(viewHeight.value - tooltipHeight.value - 10, y))})`;
});
function showTooltip(minute: number, event?: PointerEvent) {
  hoveredMinute.value = minute;
  cursor.value = null;
  if (!event) return;
  const bounds = (event.currentTarget as SVGRectElement).ownerSVGElement!.getBoundingClientRect();
  if (bounds.width && bounds.height) {
    cursor.value = {
      x: (event.clientX - bounds.left) * 800 / bounds.width,
      y: (event.clientY - bounds.top) * viewHeight.value / bounds.height,
    };
  }
}
function hoverLabel(minute: number) {
  return `${timeLabel(minute)} ${series.value.flatMap((day) => {
    const point = day.points.find((p) => p.minute === minute);
    return point ? [`${day.date}：${point.count.toLocaleString()} 人`] : [];
  }).join("、")}`;
}
</script>

<template>
  <section class="panel visitor-chart">
    <ul class="visitor-chart-legend" aria-label="日付ごとの線">
      <li v-for="day in series" :key="day.date">
        <svg width="24" height="12" aria-hidden="true">
          <line x1="0" x2="24" y1="6" y2="6" :stroke="day.color" :stroke-dasharray="day.dash" stroke-width="3" />
        </svg>
        {{ day.date }}（{{ timeLabel(day.points[0].minute) }}〜{{ timeLabel(day.points[day.points.length - 1].minute) }}）
      </li>
    </ul>
    <div class="visitor-chart-scroll" tabindex="0" aria-label="日別のグラフ。横にスクロールできます。">
      <svg :viewBox="`0 0 800 ${viewHeight}`" role="img" aria-label="日別の10分ごとの来場者数（日本時間）" @pointerleave="hoveredMinute = null">
        <title>日別の10分ごとの来場者数（日本時間）</title>
        <text x="16" y="22">人数（人）</text>
        <g v-for="tick in yTicks" :key="tick">
          <line x1="60" x2="760" :y1="220 - tick * 180 / maxCount" :y2="220 - tick * 180 / maxCount" class="visitor-chart-grid" />
          <text x="50" :y="225 - tick * 180 / maxCount" text-anchor="end">{{ tick.toLocaleString() }}</text>
        </g>
        <g v-for="point in xTicks" :key="point.minute">
          <text :x="point.x" y="246" text-anchor="middle">{{ timeLabel(point.minute) }}</text>
        </g>
        <g v-for="day in series" :key="day.date">
          <polyline v-if="!day.pending || day.points.length > 1" :data-date="day.date"
            :points="(day.pending ? day.points.slice(0, -1) : day.points).map((point) => `${point.x},${point.y}`).join(' ')"
            :stroke="day.color" :stroke-dasharray="day.dash" class="visitor-chart-line" />
          <polyline v-if="day.pending" :data-date="day.date"
            :points="day.points.slice(-2).map((point) => `${point.x},${point.y}`).join(' ')"
            :stroke="day.color" stroke-dasharray="4 4" class="visitor-chart-line visitor-chart-pending-line" />
        </g>
        <rect v-for="(point, i) in times" :key="point.minute"
          :x="i === 0 ? 60 : (times[i - 1].x + point.x) / 2"
          :width="(i === times.length - 1 ? 760 : (point.x + times[i + 1].x) / 2) - (i === 0 ? 60 : (times[i - 1].x + point.x) / 2)"
          y="40" height="180" class="visitor-chart-hover" tabindex="0" :aria-label="hoverLabel(point.minute)"
          @pointerenter="showTooltip(point.minute, $event)" @pointermove="showTooltip(point.minute, $event)"
          @pointerdown="showTooltip(point.minute, $event)"
          @focus="showTooltip(point.minute)" @blur="hoveredMinute = null" />
        <circle v-for="point in hoveredRows" :key="point.date" :data-date="point.date"
          :cx="point.x" :cy="point.y" r="4" :fill="point.color" class="visitor-chart-selected-point" />
        <g v-if="hovered && hoveredRows.length" class="visitor-chart-tooltip" role="tooltip" :transform="tooltipTransform">
          <rect :width="tooltipWidth" :height="tooltipHeight" rx="4" class="visitor-chart-tooltip-background" />
          <text x="8" y="16">{{ timeLabel(hovered.minute) }}</text>
          <g v-for="(row, i) in hoveredRows" :key="row.date">
            <rect x="8" :y="28 + i * 18" width="8" height="2" :fill="row.color" />
            <text x="22" :y="32 + i * 18">{{ row.date }}：{{ row.count.toLocaleString() }} 人</text>
          </g>
        </g>
      </svg>
    </div>
  </section>
</template>
