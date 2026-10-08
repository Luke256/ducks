<script setup lang="ts">
import { computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import {
  currentFestivalId,
  festivalSelectionLocked,
  festivals,
  festivalsError,
  festivalsLoading,
  loadFestivals,
  notices,
  dismissNotice,
} from "./state";
const route = useRoute();
const nav = [
  { href: "/event", label: "イベント", icon: "01" },
  { href: "/poster", label: "ポスター", icon: "02" },
  { href: "/sales", label: "物販", icon: "03" },
  { href: "/visitors", label: "来場者", icon: "04" },
];
const tabs = [
  { href: "/sales/cashier", label: "レジ" },
  { href: "/sales/orders", label: "売上履歴" },
  { href: "/sales/stocks", label: "販売商品" },
  { href: "/sales/items", label: "商品マスター" },
];
const isSales = computed(() => route.path.startsWith("/sales"));
const needsFestival = computed(() =>
  [
    "/poster",
    "/poster/new",
    "/visitors",
    "/sales/cashier",
    "/sales/orders",
    "/sales/stocks",
    "/sales/stocks/new",
  ].includes(route.path),
);
onMounted(loadFestivals);
</script>
<template>
  <a class="skip-link" href="#main">本文へ移動</a>
  <aside class="sidebar">
    <RouterLink to="/event" class="brand"
      ><img src="/duck.svg" alt="" />Ducks<span
        >FESTIVAL TOOLS</span
      ></RouterLink
    >
    <p class="sidebar-caption">学園祭を、スムーズに。</p>
    <nav aria-label="メインメニュー">
      <RouterLink
        v-for="item in nav"
        :key="item.href"
        :to="item.href"
        :class="{ active: route.path.startsWith(item.href) }"
        ><span class="nav-number" aria-hidden="true">{{ item.icon }}</span
        >{{ item.label
        }}<span class="nav-arrow" aria-hidden="true">↗</span></RouterLink
      >
    </nav>
    <div class="sidebar-footer">
      traP 工大祭<br /><span>運営管理ツール</span>
    </div>
  </aside>
  <div class="workspace">
    <header class="topbar">
      <span class="topbar-label">{{
        isSales
          ? "物販管理"
          : route.path.startsWith("/poster")
            ? "ポスター管理"
            : route.path.startsWith("/visitors")
              ? "来場者数"
              : "イベント管理"
      }}</span>
      <label v-if="needsFestival" class="festival-picker"
        >対象イベント
        <select
          v-model="currentFestivalId"
          :disabled="festivalsLoading || festivalSelectionLocked"
        >
          <option value="">
            {{ festivalsLoading ? "読み込み中…" : "イベントを選択" }}
          </option>
          <option
            v-for="festival in festivals"
            :key="festival.id"
            :value="festival.id"
          >
            {{ festival.name }}
          </option>
        </select>
      </label>
      <span v-else class="topbar-note">traP / 工大祭</span>
    </header>
    <main id="main" tabindex="-1">
      <div v-if="festivalsError" class="state error" role="alert">
        {{ festivalsError }}
        <button class="button secondary" @click="loadFestivals">
          再読み込み
        </button>
      </div>
      <nav v-if="isSales" class="tabs" aria-label="物販メニュー">
        <RouterLink
          v-for="tab in tabs"
          :key="tab.href"
          :to="tab.href"
          :class="{ active: route.path.startsWith(tab.href) }"
          >{{ tab.label }}</RouterLink
        >
      </nav>
      <RouterView :key="route.path" />
    </main>
  </div>
  <div class="notifications" aria-live="polite" aria-atomic="false">
    <div
      v-for="notice in notices"
      :key="notice.id"
      :class="['notice', notice.type]"
      :role="notice.type === 'error' ? 'alert' : 'status'"
    >
      <span>{{ notice.message }}</span
      ><button
        type="button"
        aria-label="通知を閉じる"
        @click="dismissNotice(notice.id)"
      >
        ×
      </button>
    </div>
  </div>
</template>
