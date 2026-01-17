import { HelloWorld } from './helloworld';

describe('HelloWorld', () => {
    it('should return "Hello World"', () => {
        const hello = new HelloWorld();
        expect(hello.message()).toBe("Hello World");
    });

    it('should return "Hello Gemini"', () => {
        const hello = new HelloWorld();
        expect(hello.message("Gemini")).toBe("Hello Gemini");
    });
});
