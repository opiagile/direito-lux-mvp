/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ['localhost', 'direito-lux.com'],
  },
  env: {
    NEXTAUTH_URL: process.env.NEXTAUTH_URL || 'http://localhost:3000',
    NEXTAUTH_SECRET: process.env.NEXTAUTH_SECRET || 'dev-secret-key',
    API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:8081',
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
  // Optimizations to prevent ChunkLoadError
  experimental: {
    // optimizeCss: true, // Disabled - causing critters MODULE_NOT_FOUND
  },
  webpack: (config, { dev, isServer }) => {
    // Optimize chunk splitting for better loading
    if (!dev && !isServer) {
      config.optimization.splitChunks = {
        ...config.optimization.splitChunks,
        cacheGroups: {
          ...config.optimization.splitChunks.cacheGroups,
          vendor: {
            test: /[\\/]node_modules[\\/]/,
            name: 'vendors',
            chunks: 'all',
            maxSize: 244000, // 244KB chunks
          },
          common: {
            name: 'common',
            minChunks: 2,
            chunks: 'all',
            maxSize: 244000,
            enforce: true,
          },
        },
      }
    }
    return config
  },
}

module.exports = nextConfig