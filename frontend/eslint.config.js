// Strict ESLint config with security rules (PR-03).
import eslintReact from '@eslint-react/eslint-plugin';
import js from '@eslint/js';
import { defineConfig } from 'eslint/config';
import noUnsanitized from 'eslint-plugin-no-unsanitized';
import reactHooks from 'eslint-plugin-react-hooks';
import security from 'eslint-plugin-security';
import globals from 'globals';
import tseslint from 'typescript-eslint';

// rawHTML bans raw HTML injection.
const rawHTML = {
  selector: "JSXAttribute[name.name='dangerouslySetInnerHTML']",
  message: 'Raw HTML is not allowed.',
};

// Text the player reads comes from en.json through t() (ADR 021): no
// letters in JSX text, in a text attribute, or in a string or template
// shown as a child. Numbers and symbols (·, ¢, …, ×) may stay.
const uiText = 'Put UI text in src/locales/en.json and show it with t() (ADR 021).';
const textAttr = 'JSXAttribute[name.name=/^(title|aria-label|aria-description|placeholder|alt|label)$/]';
const shown = ':matches(JSXElement, JSXFragment) > JSXExpressionContainer';
const either = ':matches(ConditionalExpression, LogicalExpression)';
// A text value: a string, a template, or one picked by a condition, up
// to two conditions deep; deeper code shows no text without t().
const values = (at) => [
  `${at} > Literal[value=/[A-Za-z]/]`,
  `${at} > TemplateLiteral > TemplateElement[value.raw=/[A-Za-z]/]`,
  `${at} > ${either} > Literal[value=/[A-Za-z]/]`,
  `${at} > ${either} > TemplateLiteral > TemplateElement[value.raw=/[A-Za-z]/]`,
  `${at} > ${either} > ${either} > Literal[value=/[A-Za-z]/]`,
  `${at} > ${either} > ${either} > TemplateLiteral > TemplateElement[value.raw=/[A-Za-z]/]`,
];
const hardText = [
  'JSXText[value=/[A-Za-z]/]',
  `${textAttr} > Literal[value=/[A-Za-z]/]`,
  ...values(`${textAttr} > JSXExpressionContainer`),
  ...values(shown),
].map((selector) => ({ selector, message: uiText }));

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
      // Raw HTML injection, and UI text outside en.json.
      'no-restricted-syntax': ['error', rawHTML, ...hardText],
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
    // The card gallery is a dev-only test screen (?screen=cards).
    files: ['src/screens/CardGallery.tsx'],
    rules: {
      'no-restricted-syntax': ['error', rawHTML],
    },
  },
  {
    // make i18n reads and writes translation files (16.7): fixed folders,
    // and LOCALE is checked as a language code before any path is made.
    files: ['src/translations.test.ts'],
    rules: {
      'security/detect-non-literal-fs-filename': 'off',
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
