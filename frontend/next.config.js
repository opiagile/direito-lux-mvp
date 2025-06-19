/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
  images: {
    domains: ['localhost', 'direito-lux.com'],
  },
  env: {
    NEXTAUTH_URL: process.env.NEXTAUTH_URL || 'http://localhost:3000',
    NEXTAUTH_SECRET: process.env.NEXTAUTH_SECRET || 'dev-secret-key',
    API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:8090',
    AI_SERVICE_URL: process.env.AI_SERVICE_URL || 'http://localhost:8000',
    SEARCH_SERVICE_URL: process.env.SEARCH_SERVICE_URL || 'http://localhost:8086',
    REPORT_SERVICE_URL: process.env.REPORT_SERVICE_URL || 'http://localhost:8087',
  },
  typescript: {
    ignoreBuildErrors: false,
  },
  eslint: {
    ignoreDuringBuilds: false,
  },
  output: 'standalone',
  poweredByHeader: false,
  compress: true,
  generateEtags: false,
  reactStrictMode: true,
  swcMinify: true,
}

module.exports = nextConfig