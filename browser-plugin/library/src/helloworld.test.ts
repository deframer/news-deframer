import { helloWorld } from './helloworld';

describe('helloWorld', () => {
    it('should return "Hello World"', () => {
        expect(helloWorld()).toBe("Hello World");
    });
});