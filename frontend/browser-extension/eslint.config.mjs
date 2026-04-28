import js from "@eslint/js";
import tseslint from "typescript-eslint";
import eslintReact from "@eslint-react/eslint-plugin";
import globals from "globals";
import simpleImportSort from "eslint-plugin-simple-import-sort";

export default tseslint.config(
    {
        ignores: ["dist", "node_modules"],
    },
    {
        files: ["**/*.{ts,tsx}"],
        extends: [
            js.configs.recommended,
            ...tseslint.configs.recommended,
            eslintReact.configs["recommended-typescript"],
        ],
        plugins: {
            'simple-import-sort': simpleImportSort,
        },
        languageOptions: {
            globals: {
                ...globals.browser,
                ...globals.webextensions,
            },
        },
        rules: {
            'simple-import-sort/imports': 'error',
            'simple-import-sort/exports': 'error',
        },
    },
    {
        files: ["**/*.test.{ts,tsx}"],
        rules: {
            '@eslint-react/no-unnecessary-use-prefix': 'off',
        },
    }
);
