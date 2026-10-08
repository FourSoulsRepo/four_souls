// Strict ESLint config with security rules (PR-03).
import eslintReact from '@eslint-react/eslint-plugin';
import js from '@eslint/js';
import { defineConfig } from 'eslint/config';
import noUnsanitized from 'eslint-plugin-no-unsanitized';
import reactHooks from 'eslint-plugin-react-hooks';
import security from 'eslint-plugin-security';
import globals from 'globals';
import tseslint from 'typescript-eslint';

export default defineConfig(
  {
    // Generated Wails bindings and build output are not linted.
    ignores: ['dist/', 'wailsjs/'],
  },
  {
    files: ['src/**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      tseslint.configs.strictTypeChecked,
      tseslint.configs.stylisticTypeChecked,
      eslintReact.configs['strict-type-checked'],
      reactHooks.configs.flat.recommended,
      security.configs.recommended,
      noUnsanitized.configs.recommended,
    ],
    languageOptions: {
      globals: globals.browser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
    rules: {
      // Code execution from strings.
      'no-eval': 'error',
      'no-implied-eval': 'error',
      'no-new-func': 'error',
      'no-script-url': 'error',
      // Raw HTML injection.
      'no-restricted-syntax': [
        'error',
        {
          selector: "JSXAttribute[name.name='dangerouslySetInnerHTML']",
          message: 'Raw HTML is not allowed.',
        },
      ],
      // Only the bridge may talk to Wails (A-14).
      'no-restricted-imports': [
        'error',
        {
          patterns: [
            {
              group: ['**/wailsjs/**'],
              message: 'Import Wails only through src/bridge.',
            },
          ],
        },
      ],
    },
  },
  {
    files: ['src/bridge/**'],
    rules: {
      'no-restricted-imports': 'off',
    },
  },
  {
    // Config files run in Node and are not part of the typed project.
    files: ['*.config.{js,ts}'],
    extends: [js.configs.recommended],
    languageOptions: {
      globals: globals.node,
    },
  },
);
