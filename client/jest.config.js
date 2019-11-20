module.exports = {
  errorOnDeprecated: true,
  collectCoverage: true,

  preset: "ts-jest",
  projects: [
    {
      displayName: "unit test",
      globals: {
        "ts-jest": {
          tsConfig: "<rootDir>/tests/unit/tsconfig.json"
        }
      },
      modulePaths: ["<rootDir>/src"],
      preset: "ts-jest",
      roots: ["<rootDir>/src/", "<rootDir>/tests/unit/"],
      setupFiles: ["jest-canvas-mock", "jest-webgl-canvas-mock"]
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
      modulePaths: ["<rootDir>/src"],
      preset: "jest-puppeteer",
      roots: ["<rootDir>/src/", "<rootDir>/tests/e2e/"],
      setupFilesAfterEnv: ["<rootDir>/tests/e2e/jest-setup.ts"],
      testEnvironment: "jest-environment-puppeteer",
      transform: { ".+\\.tsx?": "ts-jest" }
    }
  ],
  reporters: ["default", "jest-junit"],
  verbose: true
};
