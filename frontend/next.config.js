/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  images: {
    domains: ["i.scdn.co", "via.placeholder.com"],
  },
  env: {
    NEXT_PUBLIC_GRAPHQL_URL:
      process.env.NEXT_PUBLIC_GRAPHQL_URL || "http://localhost:8080/query",
  },
};

module.exports = nextConfig;
