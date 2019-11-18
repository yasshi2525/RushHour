module.exports = {
  collectCoverage: true,
  errorOnDeprecated: true,
  globalSetup: "./setup.js",
  globalTeardown: "./teardown.js",
  moduleNameMapper: {
    "\\.(png|svg|jpg|gif|woff|woff2|eot|ttf)$":
      "<rootDir>/tests/__mocks__/fileMock.js",
    "\\.(css|less)$": "identity-obj-proxy"
  },
  modulePathIgnorePatterns: ["<rootDir>/dist/"],
  modulePaths: ["<rootDir>/src/"],
  preset: "jest-runner-prettier",
  reporters: ["default", "jest-junit"],
  runner: "prettier",
  roots: ["<rootDir>/src/", "<rootDir>/tests/"],
  setupFilesAfterEnv: ["./setup.ts"],
  setupFiles: ["./client.js"],
  testEnvironment: "./puppeteer_environment.js",
  testMatch: ["<rootDir>/tests/**/*.test.ts"]
};
