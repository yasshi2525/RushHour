module.exports = {
    preset: "ts-jest",
    testMatch: ["<rootDir>/webtests/**/*.test.ts"],
    collectCoverage: true,
    collectCoverageFrom: [
        "<rootDir>/web/**/*.(ts|tsx)"
    ],
    moduleNameMapper: {
        "^@/(.+)": "<rootDir>/web/$1"
    },
    setupFiles: ["jest-canvas-mock"]
};