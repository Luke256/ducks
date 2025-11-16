import { Metadata } from "next";
import NewPosterPageClient from "./pageClient";

export const metadata: Metadata = {
  title: "新規ポスター登録",
  description: "新規ポスターの登録",
};

const NewPosterPage = () => {
    return <NewPosterPageClient />;
}

export default NewPosterPage;