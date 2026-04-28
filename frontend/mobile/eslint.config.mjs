import js from '@eslint/js';
import eslintReact from '@eslint-react/eslint-plugin';
import globals from 'globals';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  {
    ignores: [
      'android/**',
      'ios/**',
      'node_modules/**',
      '.bundle/**',
      'eslint.config.js',
      'eslint.config.mjs',
      '.eslintrc.js',
      '.prettierrc.js',
      'babel.config.js',
      'jest.config.js',
      'metro.config.js',
      'webpack.web.js',
      'scripts/**',
      'index.js',
      '**/*.js',
      '**/*.jsx',
    ],
  },
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      ...tseslint.configs.recommended,
      eslintReact.configs['recommended-typescript'],
    ],
    languageOptions: {
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
      },
      globals: {
        ...globals.browser,
        ...globals.node,
        ...globals.jest,
      },
    },
  },
  {
    files: ['**/*.test.{ts,tsx}', '**/__tests__/**/*.{ts,tsx}'],
    rules: {
      '@eslint-react/no-unnecessary-use-prefix': 'off',
    },
  },
);
