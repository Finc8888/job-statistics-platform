const esbuild = require('esbuild');
const fs = require('fs');
const path = require('path');

const isWatch = process.argv.includes('--watch');

// Создаем директорию dist если её нет
if (!fs.existsSync('dist')) {
  fs.mkdirSync('dist', { recursive: true });
}

// Копируем index.html
if (fs.existsSync('public/index.html')) {
  fs.copyFileSync('public/index.html', 'dist/index.html');
}

const buildOptions = {
  entryPoints: ['src/index.tsx'],
  bundle: true,
  outfile: 'dist/bundle.js',
  minify: !isWatch,
  sourcemap: isWatch,
  target: 'es2020',
  loader: {
    '.tsx': 'tsx',
    '.ts': 'ts',
  },
  define: {
    'process.env.NODE_ENV': isWatch ? '"development"' : '"production"',
  },
};

async function build() {
  try {
    if (isWatch) {
      const ctx = await esbuild.context(buildOptions);
      await ctx.watch();
      console.log('👀 Watching for changes...');
      console.log('📦 Open http://localhost:3000 to view the app');
    } else {
      await esbuild.build(buildOptions);
      console.log('✅ Build completed successfully!');
    }
  } catch (error) {
    console.error('❌ Build failed:', error);
    process.exit(1);
  }
}

build();
