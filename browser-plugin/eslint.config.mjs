import js from "@eslint/js";
import tseslint from "typescript-eslint";
import react from "eslint-plugin-react";
import globals from "globals";

export default tseslint.config(
    {
        ignores: ["dist", "node_modules"],
    },
    js.configs.recommended,
    ...tseslint.configs.recommended,
    {
        files: ["**/*.{ts,tsx}"],
        plugins: {
            react,
        },
        languageOptions: {
            globals: {
                ...globals.browser,
                ...globals.webextensions,
            },
        },
        rules: {
            ...react.configs.recommended.rules,
            "react/react-in-jsx-scope": "off",
        },
        settings: {
            react: {
                version: "detect",
            },
        },
    }
);