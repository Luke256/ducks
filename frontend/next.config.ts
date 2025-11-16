import type { NextConfig } from "next";

const isDev = process.env.NODE_ENV !== "production";

const nextConfig: NextConfig = {
  /* config options here */
  images: {
    remotePatterns: [
      new URL(`${process.env.NEXT_PUBLIC_API_URL}/**` || "http://localhost:8080/**"),
    ],
    dangerouslyAllowLocalIP: isDev,
  }
};

export default nextConfig;
