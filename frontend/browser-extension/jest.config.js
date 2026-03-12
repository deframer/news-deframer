module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  roots: ['<rootDir>/src', '<rootDir>/../shared'],
  testMatch: ['**/*.test.ts', '**/*.test.tsx'],
  moduleDirectories: ['node_modules', '<rootDir>/node_modules'],
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
  transform: {
    '^.+\\.tsx?$': 'ts-jest',
  },
};
