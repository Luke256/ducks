import { createRouter, createWebHistory } from "vue-router";
const router = createRouter({
    history: createWebHistory(),
    scrollBehavior: () => ({ top: 0 }),
    routes: [
        { path: "/", redirect: "/event" },
        {
            path: "/event",
            component: () => import("./views/EventsView.vue"),
            meta: { title: "イベント管理" },
        },
        {
            path: "/event/:eventId",
            component: () => import("./views/EventDetailView.vue"),
            meta: { title: "イベント詳細" },
        },
        {
            path: "/poster",
            component: () => import("./views/PostersView.vue"),
            meta: { title: "ポスター管理" },
        },
        {
            path: "/poster/new",
            component: () => import("./views/PosterFormView.vue"),
            meta: { title: "ポスター登録" },
        },
        {
            path: "/poster/detail/:posterId",
            component: () => import("./views/PosterDetailView.vue"),
            meta: { title: "ポスター詳細" },
        },
        { path: "/sales", redirect: "/sales/cashier" },
        {
            path: "/visitors",
            component: () => import("./views/VisitorsView.vue"),
            meta: { title: "来場者数" },
        },
        {
            path: "/sales/cashier",
            component: () => import("./views/CashierView.vue"),
            meta: { title: "レジ" },
        },
        {
            path: "/sales/orders",
            component: () => import("./views/OrdersView.vue"),
            meta: { title: "売上管理" },
        },
        {
            path: "/sales/items",
            component: () => import("./views/ItemsView.vue"),
            meta: { title: "商品管理" },
        },
        {
            path: "/sales/items/new",
            component: () => import("./views/ItemFormView.vue"),
            meta: { title: "商品登録" },
        },
        {
            path: "/sales/items/:itemId",
            component: () => import("./views/ItemFormView.vue"),
            meta: { title: "商品詳細" },
        },
        {
            path: "/sales/stocks",
            component: () => import("./views/StocksView.vue"),
            meta: { title: "販売商品" },
        },
        {
            path: "/sales/stocks/new",
            component: () => import("./views/StockFormView.vue"),
            meta: { title: "販売商品登録" },
        },
        {
            path: "/sales/stocks/:stockId",
            component: () => import("./views/StockFormView.vue"),
            meta: { title: "販売商品詳細" },
        },
        {
            path: "/:pathMatch(.*)*",
            component: () => import("./views/NotFoundView.vue"),
            meta: { title: "ページが見つかりません" },
        },
    ],
});
router.afterEach((to) => {
    document.title = `${to.meta.title || "管理"} | Ducks`;
});
export default router;
