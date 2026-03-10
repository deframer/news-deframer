import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';
import { StatusBar, StyleSheet, Text } from 'react-native';

function App() {
  return (
    <SafeAreaProvider>
      <StatusBar barStyle="dark-content" backgroundColor="#f6f3ea" />
      <SafeAreaView style={styles.container}>
        <Text style={styles.title}>News Deframer</Text>
        <Text style={styles.subtitle}>Hello world</Text>
      </SafeAreaView>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#f6f3ea',
    paddingHorizontal: 24,
  },
  title: {
    marginBottom: 8,
    color: '#1f2937',
    fontSize: 32,
    fontWeight: '700',
  },
  subtitle: {
    color: '#4b5563',
    fontSize: 18,
  },
});

export default App;
