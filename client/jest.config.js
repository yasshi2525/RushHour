module.exports = {
  errorOnDeprecated: true,
  collectCoverage: true,
  modulePaths: ["<rootDir>/src"],

  preset: "ts-jest",
  projects: [
    {
      displayName: "unit test",
      globals: {
        "ts-jest": {
          tsConfig: "<rootDir>/tests/unit/tsconfig.json"
        }
      },
      roots: ["<rootDir>/src/", "<rootDir>/tests/unit/"],
      preset: "ts-jest"
    },
    {
      displayName: "e2e test",
      globalSetup: "jest-environment-puppeteer/setup",
      globalTeardown: "jest-environment-puppeteer/teardown",
      globals: {
        "ts-jest": {
          tsConfig: "<rootDir>/tests/e2e/tsconfig.json"
        }
      },
      roots: ["<rootDir>/src/", "<rootDir>/tests/e2e/"],
      preset: "jest-puppeteer",
      setupFilesAfterEnv: ["<rootDir>/tests/e2e/jest-setup.ts"],
      testEnvironment: "jest-environment-puppeteer",
      transform: { ".+\\.tsx?": "ts-jest" }
    }
  ],
  reporters: ["default", "jest-junit"]
};
