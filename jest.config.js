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
    setupFilesAfterEnv: ["./setup.ts"],
    setupFiles: [ "./client.js" ],
    globalSetup: "./setup.js",
    globalTeardown: "./teardown.js",
    testEnvironment: "./puppeteer_environment.js",
};