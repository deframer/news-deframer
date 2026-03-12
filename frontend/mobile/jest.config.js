module.exports = {
  preset: 'react-native',
  roots: ['<rootDir>/src', '<rootDir>/../shared'],
  testMatch: ['**/*.test.ts', '**/*.test.tsx'],
  modulePaths: ['<rootDir>/node_modules'],
  moduleDirectories: ['node_modules', '<rootDir>/node_modules'],
  moduleNameMapper: {
    '\\.svg$': '<rootDir>/__mocks__/svgMock.js',
  },
};
