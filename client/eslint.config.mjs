import tsPlugin from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';
import pluginVue from 'eslint-plugin-vue';
import globals from 'globals';
import vueParser from 'vue-eslint-parser';

export default [
    {
        ignores: [
            'dist/**',
            'src-capacitor/**',
            'src-cordova/**',
            '.quasar/**',
            'node_modules/**',
            'src-ssr/**',
            'src/stores/store-flag.d.ts',
        ],
    },

    ...pluginVue.configs['flat/essential'],

    {
        files: ['**/*.{js,mjs,cjs,ts,vue}'],
        languageOptions: {
            parser: vueParser,
            parserOptions: {
                parser: tsParser,
                ecmaVersion: 'latest',
                sourceType: 'module',
                extraFileExtensions: ['.vue'],
            },
            globals: {
                ...globals.browser,
                ...globals.node,
                ga: 'readonly',
                cordova: 'readonly',
                __statics: 'readonly',
                __QUASAR_SSR__: 'readonly',
                __QUASAR_SSR_SERVER__: 'readonly',
                __QUASAR_SSR_CLIENT__: 'readonly',
                __QUASAR_SSR_PWA__: 'readonly',
                process: 'readonly',
                Capacitor: 'readonly',
                chrome: 'readonly',
            },
        },
        plugins: {
            '@typescript-eslint': tsPlugin,
        },
        rules: {
            ...tsPlugin.configs.recommended.rules,
            'prefer-promise-reject-errors': 'off',
            quotes: ['warn', 'single', { avoidEscape: true }],
            '@typescript-eslint/explicit-function-return-type': 'off',
            '@typescript-eslint/no-var-requires': 'off',
            'no-unused-vars': 'off',
            'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
        },
    },
];
