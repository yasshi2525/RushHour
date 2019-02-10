module.exports = {
    preset: "ts-jest",
    testMatch: ["<rootDir>/webtests/**/*.test.(ts|tsx)"],
    collectCoverage: true,
    collectCoverageFrom: [
        "**/*.(ts|tsx)"
    ],
    moduleNameMapper: {
        "^@/(.+)": "<rootDir>/web/$1",
        "\\.(css|less)$": "identity-obj-proxy"
    },
    setupFiles: ["jest-canvas-mock"],
    setupFilesAfterEnv: ["<rootDir>/setup.ts"]
};