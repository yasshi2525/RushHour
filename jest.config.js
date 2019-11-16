module.exports = {
    preset: "ts-jest",
    testMatch: ["<rootDir>/client/test/**/*.test.(ts|tsx)"],
    collectCoverage: true,
    moduleNameMapper: {
        "^@/(.+)": "<rootDir>/client/src/$1",
        "\\.(css|less)$": "identity-obj-proxy"
    },
    setupFilesAfterEnv: ["./setup.ts"],
    setupFiles: [ "./client.js" ],
    globalSetup: "./setup.js",
    globalTeardown: "./teardown.js",
    testEnvironment: "./puppeteer_environment.js",
    reporters: ["default", "jest-junit"]
};