import { Metadata } from "next";
import PosterPageClient from "./pageClient";

export const metadata: Metadata = {
  title: "ポスター一覧",
  description: "ポスター管理ページ",
}

export default function PosterPage() {
  return <PosterPageClient />;
}