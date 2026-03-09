module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  roots: ['<rootDir>/src', '<rootDir>/../shared'],
  testMatch: ['**/*.test.ts', '**/*.test.tsx'],
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
  moduleDirectories: ['node_modules', '<rootDir>/node_modules'],
  moduleNameMapper: {
    '^@browser/(.*)$': '<rootDir>/src/$1',
    '^@frontend-shared/(.*)$': '<rootDir>/../shared/$1',
  },
  transform: {
    '^.+\\.tsx?$': 'ts-jest',
  },
};
