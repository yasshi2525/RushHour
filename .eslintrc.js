module.exports = {
    "env": {
        "browser": true,
        "node": true
    },
    "extends": "eslint:recommended",
    "parserOptions": {
        "sourceType": "module",
        "ecmaFeatures": {
            "jsx": true
        },
        "ecmaVersion": 6
    },
    "plugins": [
        "react"
    ],
    "rules": {
        "semi": [2, "always"],
        "quotes": [2, "double"],
        // JSX で react の import を必須化
        "react/jsx-uses-react": 1,
        // 変数使用有無の検査対象に jsx記法を含める
        "react/jsx-uses-vars": 1,
        // console.log の使用を警告
        "no-console": 1
    }
};