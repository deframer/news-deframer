module.exports = {
  preset: 'react-native',
  modulePaths: ['<rootDir>/node_modules'],
  moduleNameMapper: {
    '\\.svg$': '<rootDir>/__mocks__/svgMock.js',
  },
};
