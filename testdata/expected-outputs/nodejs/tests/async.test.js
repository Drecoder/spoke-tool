/**
 * Async Operations Tests
 * 
 * Tests for asynchronous operations including promises, callbacks,
 * events, streams, and various async patterns.
 */

const { 
    delay,
    fetchWithTimeout,
    retry,
    parallel,
    series,
    waterfall,
    queue,
    debounce,
    throttle,
    memoize,
    EventEmitter
} = require('../src/async');

// ============================================================================
// Promise Tests
// ============================================================================

describe('Promise Operations', () => {
    
    describe('delay', () => {
        test('should resolve after specified delay', async () => {
            const start = Date.now();
            await delay(100);
            const elapsed = Date.now() - start;
            expect(elapsed).toBeGreaterThanOrEqual(90);
            expect(elapsed).toBeLessThan(150);
        });

        test('should handle zero delay', async () => {
            const start = Date.now();
            await delay(0);
            const elapsed = Date.now() - start;
            expect(elapsed).toBeLessThan(50);
        });

        test('should handle negative delay (treat as zero)', async () => {
            const start = Date.now();
            await delay(-100);
            const elapsed = Date.now() - start;
            expect(elapsed).toBeLessThan(50);
        });

        test('should not block other operations', async () => {
            let flag = false;
            setTimeout(() => { flag = true; }, 50);
            
            await delay(100);
            expect(flag).toBe(true);
        });
    });

    describe('fetchWithTimeout', () => {
        test('should resolve before timeout', async () => {
            const promise = new Promise(resolve => 
                setTimeout(() => resolve('success'), 50)
            );
            
            const result = await fetchWithTimeout(promise, 100);
            expect(result).toBe('success');
        });

        test('should reject on timeout', async () => {
            const promise = new Promise(resolve => 
                setTimeout(() => resolve('too late'), 200)
            );
            
            await expect(fetchWithTimeout(promise, 100))
                .rejects
                .toThrow('Operation timed out');
        });

        test('should handle immediate resolve', async () => {
            const promise = Promise.resolve('instant');
            const result = await fetchWithTimeout(promise, 100);
            expect(result).toBe('instant');
        });

        test('should handle immediate reject', async () => {
            const promise = Promise.reject(new Error('immediate error'));
            
            await expect(fetchWithTimeout(promise, 100))
                .rejects
                .toThrow('immediate error');
        });

        test('should handle zero timeout', async () => {
            const promise = new Promise(resolve => 
                setTimeout(() => resolve('success'), 10)
            );
            
            await expect(fetchWithTimeout(promise, 0))
                .rejects
                .toThrow('Operation timed out');
        });

        test('should handle negative timeout', async () => {
            const promise = Promise.resolve('success');
            
            await expect(fetchWithTimeout(promise, -100))
                .rejects
                .toThrow('Operation timed out');
        });
    });

    describe('retry', () => {
        test('should succeed on first attempt', async () => {
            const fn = jest.fn().mockResolvedValue('success');
            
            const result = await retry(fn, { attempts: 3 });
            
            expect(result).toBe('success');
            expect(fn).toHaveBeenCalledTimes(1);
        });

        test('should retry on failure and succeed', async () => {
            const fn = jest.fn()
                .mockRejectedValueOnce(new Error('fail 1'))
                .mockRejectedValueOnce(new Error('fail 2'))
                .mockResolvedValue('success');
            
            const result = await retry(fn, { attempts: 3, delay: 10 });
            
            expect(result).toBe('success');
            expect(fn).toHaveBeenCalledTimes(3);
        });

        test('should fail after max attempts', async () => {
            const fn = jest.fn().mockRejectedValue(new Error('persistent failure'));
            
            await expect(retry(fn, { attempts: 3, delay: 10 }))
                .rejects
                .toThrow('persistent failure');
            
            expect(fn).toHaveBeenCalledTimes(3);
        });

        test('should respect delay between retries', async () => {
            const fn = jest.fn()
                .mockRejectedValueOnce(new Error('fail'))
                .mockResolvedValue('success');
            
            const start = Date.now();
            await retry(fn, { attempts: 2, delay: 100 });
            const elapsed = Date.now() - start;
            
            expect(elapsed).toBeGreaterThanOrEqual(90);
            expect(fn).toHaveBeenCalledTimes(2);
        });

        test('should support exponential backoff', async () => {
            const fn = jest.fn()
                .mockRejectedValueOnce(new Error('fail 1'))
                .mockRejectedValueOnce(new Error('fail 2'))
                .mockResolvedValue('success');
            
            const start = Date.now();
            await retry(fn, { 
                attempts: 3, 
                delay: 50, 
                backoff: 2 
            });
            const elapsed = Date.now() - start;
            
            // First retry: 50ms, second retry: 100ms
            expect(elapsed).toBeGreaterThanOrEqual(140);
            expect(fn).toHaveBeenCalledTimes(3);
        });

        test('should handle retry on specific errors only', async () => {
            const fn = jest.fn()
                .mockRejectedValueOnce(new Error('retryable'))
                .mockRejectedValueOnce(new Error('fatal'))
                .mockResolvedValue('success');
            
            const shouldRetry = (err) => err.message === 'retryable';
            
            await expect(retry(fn, { attempts: 3, shouldRetry }))
                .rejects
                .toThrow('fatal');
            
            expect(fn).toHaveBeenCalledTimes(2);
        });

        test('should pass arguments to retried function', async () => {
            const fn = jest.fn().mockResolvedValue('success');
            
            await retry(fn, { attempts: 3 }, 'arg1', 'arg2');
            
            expect(fn).toHaveBeenCalledWith('arg1', 'arg2');
        });
    });
});

// ============================================================================
// Control Flow Tests
// ============================================================================

describe('Control Flow', () => {
    
    describe('parallel', () => {
        test('should run all tasks in parallel', async () => {
            const order = [];
            const tasks = [
                async () => {
                    await delay(30);
                    order.push(1);
                    return 1;
                },
                async () => {
                    await delay(20);
                    order.push(2);
                    return 2;
                },
                async () => {
                    await delay(10);
                    order.push(3);
                    return 3;
                }
            ];
            
            const results = await parallel(tasks);
            
            expect(results).toEqual([1, 2, 3]);
            expect(order).toEqual([3, 2, 1]); // Fastest completes first
        });

        test('should limit concurrency', async () => {
            let running = 0;
            let maxRunning = 0;
            
            const tasks = Array(10).fill().map((_, i) => async () => {
                running++;
                maxRunning = Math.max(maxRunning, running);
                await delay(50);
                running--;
                return i;
            });
            
            const results = await parallel(tasks, { concurrency: 3 });
            
            expect(maxRunning).toBe(3);
            expect(results.length).toBe(10);
        });

        test('should handle empty task array', async () => {
            const results = await parallel([]);
            expect(results).toEqual([]);
        });

        test('should propagate errors', async () => {
            const tasks = [
                async () => 'ok',
                async () => { throw new Error('task failed'); },
                async () => 'also ok'
            ];
            
            await expect(parallel(tasks)).rejects.toThrow('task failed');
        });

        test('should handle mixed sync and async tasks', async () => {
            const tasks = [
                () => 'sync',
                async () => {
                    await delay(10);
                    return 'async';
                },
                () => 'also sync'
            ];
            
            const results = await parallel(tasks);
            expect(results).toEqual(['sync', 'async', 'also sync']);
        });

        test('should respect timeout per task', async () => {
            const tasks = [
                async () => {
                    await delay(100);
                    return 'slow';
                },
                async () => 'fast'
            ];
            
            await expect(parallel(tasks, { timeout: 50 }))
                .rejects
                .toThrow('timed out');
        });
    });

    describe('series', () => {
        test('should run tasks in order', async () => {
            const order = [];
            const tasks = [
                async () => {
                    await delay(30);
                    order.push(1);
                    return 1;
                },
                async () => {
                    await delay(20);
                    order.push(2);
                    return 2;
                },
                async () => {
                    await delay(10);
                    order.push(3);
                    return 3;
                }
            ];
            
            const results = await series(tasks);
            
            expect(results).toEqual([1, 2, 3]);
            expect(order).toEqual([1, 2, 3]); // Runs in order
        });

        test('should stop on error', async () => {
            const task2 = jest.fn().mockResolvedValue('should not run');
            
            const tasks = [
                async () => 'ok',
                async () => { throw new Error('fail'); },
                task2
            ];
            
            await expect(series(tasks)).rejects.toThrow('fail');
            expect(task2).not.toHaveBeenCalled();
        });

        test('should handle empty task array', async () => {
            const results = await series([]);
            expect(results).toEqual([]);
        });

        test('should pass results between tasks', async () => {
            const tasks = [
                () => 1,
                (prev) => prev + 1,
                (prev) => prev * 2
            ];
            
            const results = await series(tasks);
            expect(results).toEqual([1, 2, 4]);
        });
    });

    describe('waterfall', () => {
        test('should pass result from one task to next', async () => {
            const result = await waterfall([
                () => 1,
                (val) => val + 1,
                (val) => val * 2,
                (val) => `result: ${val}`
            ]);
            
            expect(result).toBe('result: 4');
        });

        test('should handle async tasks', async () => {
            const result = await waterfall([
                async () => {
                    await delay(10);
                    return 1;
                },
                async (val) => {
                    await delay(10);
                    return val + 1;
                },
                (val) => val * 2
            ]);
            
            expect(result).toBe(4);
        });

        test('should stop on error', async () => {
            const task3 = jest.fn();
            
            await expect(waterfall([
                () => 1,
                () => { throw new Error('fail'); },
                task3
            ])).rejects.toThrow('fail');
            
            expect(task3).not.toHaveBeenCalled();
        });

        test('should handle single task', async () => {
            const result = await waterfall([() => 42]);
            expect(result).toBe(42);
        });

        test('should handle empty array', async () => {
            const result = await waterfall([]);
            expect(result).toBeUndefined();
        });
    });

    describe('queue', () => {
        test('should process tasks in order', async () => {
            const order = [];
            const q = queue(async (task) => {
                await delay(task.delay);
                order.push(task.id);
                return task.id;
            }, { concurrency: 1 });
            
            q.push({ id: 1, delay: 30 });
            q.push({ id: 2, delay: 20 });
            q.push({ id: 3, delay: 10 });
            
            await q.drain();
            
            expect(order).toEqual([1, 2, 3]); // Processed in order due to concurrency 1
        });

        test('should respect concurrency limit', async () => {
            let running = 0;
            let maxRunning = 0;
            
            const q = queue(async (task) => {
                running++;
                maxRunning = Math.max(maxRunning, running);
                await delay(50);
                running--;
                return task;
            }, { concurrency: 3 });
            
            for (let i = 0; i < 10; i++) {
                q.push(i);
            }
            
            await q.drain();
            
            expect(maxRunning).toBe(3);
        });

        test('should handle task errors', async () => {
            const q = queue(async (task) => {
                if (task === 3) throw new Error('task 3 failed');
                return task;
            });
            
            q.push(1);
            q.push(2);
            q.push(3);
            q.push(4);
            
            const results = [];
            q.on('task:error', (err, task) => {
                results.push({ error: err.message, task });
            });
            
            await expect(q.drain()).rejects.toThrow('task 3 failed');
            expect(results).toContainEqual({ error: 'task 3 failed', task: 3 });
        });

        test('should provide task completion order', async () => {
            const completed = [];
            const q = queue(async (task) => {
                await delay(task.delay);
                return task.id;
            }, { concurrency: 2 });
            
            q.push({ id: 1, delay: 30 });
            q.push({ id: 2, delay: 20 });
            q.push({ id: 3, delay: 10 });
            
            q.on('task:complete', (result, task) => {
                completed.push({ result, task: task.id });
            });
            
            await q.drain();
            
            // With concurrency 2, tasks 1 and 2 start together, 
            // task 3 starts when one finishes
            expect(completed.length).toBe(3);
        });

        test('should support pausing and resuming', async () => {
            const processed = [];
            const q = queue(async (task) => {
                processed.push(task);
                await delay(10);
            }, { concurrency: 2 });
            
            q.push(1);
            q.push(2);
            q.push(3);
            q.push(4);
            
            await delay(5);
            q.pause();
            
            expect(q.isPaused()).toBe(true);
            
            q.push(5);
            expect(processed.length).toBe(2); // Only first two processed
            
            q.resume();
            await q.drain();
            
            expect(processed.length).toBe(5);
        });

        test('should report queue statistics', async () => {
            const q = queue(async () => await delay(50), { concurrency: 2 });
            
            q.push(1);
            q.push(2);
            q.push(3);
            q.push(4);
            
            expect(q.length()).toBe(4);
            expect(q.running()).toBe(0);
            
            await delay(10);
            expect(q.running()).toBe(2);
            
            await q.drain();
            expect(q.length()).toBe(0);
            expect(q.running()).toBe(0);
        });
    });
});

// ============================================================================
// Rate Limiting Tests
// ============================================================================

describe('Rate Limiting', () => {
    
    describe('debounce', () => {
        test('should debounce function calls', async () => {
            const fn = jest.fn();
            const debounced = debounce(fn, 100);
            
            debounced();
            debounced();
            debounced();
            
            expect(fn).not.toHaveBeenCalled();
            
            await delay(150);
            expect(fn).toHaveBeenCalledTimes(1);
        });

        test('should pass latest arguments', async () => {
            const fn = jest.fn();
            const debounced = debounce(fn, 100);
            
            debounced(1);
            debounced(2);
            debounced(3);
            
            await delay(150);
            expect(fn).toHaveBeenCalledWith(3);
        });

        test('should support immediate execution', async () => {
            const fn = jest.fn();
            const debounced = debounce(fn, 100, true);
            
            debounced(1);
            expect(fn).toHaveBeenCalledWith(1);
            
            debounced(2);
            debounced(3);
            
            await delay(150);
            expect(fn).toHaveBeenCalledTimes(2); // First immediate, then last after delay
            expect(fn).toHaveBeenLastCalledWith(3);
        });

        test('should cancel pending execution', async () => {
            const fn = jest.fn();
            const debounced = debounce(fn, 100);
            
            debounced();
            debounced.cancel();
            
            await delay(150);
            expect(fn).not.toHaveBeenCalled();
        });

        test('should handle rapid successive calls', async () => {
            const fn = jest.fn();
            const debounced = debounce(fn, 50);
            
            for (let i = 0; i < 100; i++) {
                debounced(i);
                await delay(10);
            }
            
            await delay(100);
            expect(fn).toHaveBeenCalledTimes(1);
            expect(fn).toHaveBeenCalledWith(99);
        });
    });

    describe('throttle', () => {
        test('should throttle function calls', async () => {
            const fn = jest.fn();
            const throttled = throttle(fn, 100);
            
            throttled();
            throttled();
            throttled();
            
            expect(fn).toHaveBeenCalledTimes(1); // First call goes through
            
            await delay(150);
            throttled();
            expect(fn).toHaveBeenCalledTimes(2);
        });

        test('should respect leading edge', async () => {
            const fn = jest.fn();
            const throttled = throttle(fn, 100, { leading: true, trailing: false });
            
            throttled(1);
            throttled(2);
            throttled(3);
            
            expect(fn).toHaveBeenCalledTimes(1);
            expect(fn).toHaveBeenCalledWith(1);
            
            await delay(150);
            throttled(4);
            expect(fn).toHaveBeenCalledTimes(2);
            expect(fn).toHaveBeenCalledWith(4);
        });

        test('should respect trailing edge', async () => {
            const fn = jest.fn();
            const throttled = throttle(fn, 100, { leading: false, trailing: true });
            
            throttled(1);
            throttled(2);
            throttled(3);
            
            expect(fn).not.toHaveBeenCalled();
            
            await delay(150);
            expect(fn).toHaveBeenCalledTimes(1);
            expect(fn).toHaveBeenCalledWith(3);
        });

        test('should handle both leading and trailing', async () => {
            const fn = jest.fn();
            const throttled = throttle(fn, 100, { leading: true, trailing: true });
            
            throttled(1);
            throttled(2);
            throttled(3);
            
            expect(fn).toHaveBeenCalledTimes(1);
            expect(fn).toHaveBeenCalledWith(1);
            
            await delay(150);
            expect(fn).toHaveBeenCalledTimes(2);
            expect(fn).toHaveBeenCalledWith(3);
        });

        test('should cancel pending trailing call', async () => {
            const fn = jest.fn();
            const throttled = throttle(fn, 100, { trailing: true });
            
            throttled(1);
            throttled.cancel();
            
            await delay(150);
            expect(fn).not.toHaveBeenCalled();
        });
    });

    describe('memoize', () => {
        test('should cache function results', async () => {
            const fn = jest.fn().mockImplementation(async (x) => x * 2);
            const memoized = memoize(fn);
            
            const result1 = await memoized(5);
            const result2 = await memoized(5);
            
            expect(result1).toBe(10);
            expect(result2).toBe(10);
            expect(fn).toHaveBeenCalledTimes(1);
        });

        test('should use different cache keys', async () => {
            const fn = jest.fn().mockImplementation(async (x) => x * 2);
            const memoized = memoize(fn);
            
            await memoized(5);
            await memoized(6);
            await memoized(5);
            
            expect(fn).toHaveBeenCalledTimes(2);
        });

        test('should support custom key resolver', async () => {
            const fn = jest.fn().mockImplementation(async (obj) => obj.value);
            const keyResolver = (obj) => obj.id;
            const memoized = memoize(fn, { keyResolver });
            
            await memoized({ id: 1, value: 10 });
            await memoized({ id: 1, value: 20 }); // Same key, returns cached 10
            await memoized({ id: 2, value: 30 });
            
            expect(fn).toHaveBeenCalledTimes(2);
        });

        test('should respect TTL', async () => {
            const fn = jest.fn().mockImplementation(async (x) => x * 2);
            const memoized = memoize(fn, { ttl: 100 });
            
            await memoized(5);
            await memoized(5); // Cached
            
            await delay(150);
            await memoized(5); // Expired, should call again
            
            expect(fn).toHaveBeenCalledTimes(2);
        });

        test('should handle cache size limits', async () => {
            const fn = jest.fn().mockImplementation(async (x) => x);
            const memoized = memoize(fn, { maxSize: 2 });
            
            await memoized(1);
            await memoized(2);
            await memoized(1); // Cached
            await memoized(3); // Should evict oldest (2)
            await memoized(2); // Need to call again
            
            expect(fn).toHaveBeenCalledTimes(4); // 1,2,3,2
        });

        test('should clear cache', async () => {
            const fn = jest.fn().mockImplementation(async (x) => x);
            const memoized = memoize(fn);
            
            await memoized(1);
            await memoized(2);
            
            memoized.clear();
            
            await memoized(1);
            await memoized(2);
            
            expect(fn).toHaveBeenCalledTimes(4);
        });

        test('should handle rejected promises', async () => {
            const fn = jest.fn().mockImplementation(async (x) => {
                if (x < 0) throw new Error('negative');
                return x;
            });
            
            const memoized = memoize(fn);
            
            await expect(memoized(-1)).rejects.toThrow('negative');
            await expect(memoized(-1)).rejects.toThrow('negative'); // Should not cache rejection
            await memoized(5);
            await memoized(5); // Should cache success
            
            expect(fn).toHaveBeenCalledTimes(3); // -1, -1, 5
        });
    });
});

// ============================================================================
// Event Emitter Tests
// ============================================================================

describe('Event Emitter', () => {
    
    test('should emit and handle events', () => {
        const emitter = new EventEmitter();
        const handler = jest.fn();
        
        emitter.on('test', handler);
        emitter.emit('test', 1, 2, 3);
        
        expect(handler).toHaveBeenCalledWith(1, 2, 3);
    });

    test('should handle multiple listeners', () => {
        const emitter = new EventEmitter();
        const handler1 = jest.fn();
        const handler2 = jest.fn();
        
        emitter.on('test', handler1);
        emitter.on('test', handler2);
        emitter.emit('test');
        
        expect(handler1).toHaveBeenCalled();
        expect(handler2).toHaveBeenCalled();
    });

    test('should remove listeners', () => {
        const emitter = new EventEmitter();
        const handler = jest.fn();
        
        emitter.on('test', handler);
        emitter.off('test', handler);
        emitter.emit('test');
        
        expect(handler).not.toHaveBeenCalled();
    });

    test('should support once listeners', () => {
        const emitter = new EventEmitter();
        const handler = jest.fn();
        
        emitter.once('test', handler);
        emitter.emit('test');
        emitter.emit('test');
        
        expect(handler).toHaveBeenCalledTimes(1);
    });

    test('should handle error events', () => {
        const emitter = new EventEmitter();
        const errorHandler = jest.fn();
        
        emitter.on('error', errorHandler);
        emitter.emit('error', new Error('test error'));
        
        expect(errorHandler).toHaveBeenCalledWith(expect.any(Error));
    });

    test('should return listener count', () => {
        const emitter = new EventEmitter();
        
        emitter.on('test', () => {});
        emitter.on('test', () => {});
        
        expect(emitter.listenerCount('test')).toBe(2);
    });

    test('should handle event names with symbols', () => {
        const emitter = new EventEmitter();
        const sym = Symbol('test');
        const handler = jest.fn();
        
        emitter.on(sym, handler);
        emitter.emit(sym, 'data');
        
        expect(handler).toHaveBeenCalledWith('data');
    });

    test('should support prepend listeners', () => {
        const emitter = new EventEmitter();
        const order = [];
        
        emitter.on('test', () => order.push(1));
        emitter.prependListener('test', () => order.push(2));
        emitter.emit('test');
        
        expect(order).toEqual([2, 1]);
    });

    test('should remove all listeners', () => {
        const emitter = new EventEmitter();
        const handler1 = jest.fn();
        const handler2 = jest.fn();
        
        emitter.on('test', handler1);
        emitter.on('test', handler2);
        emitter.removeAllListeners('test');
        emitter.emit('test');
        
        expect(handler1).not.toHaveBeenCalled();
        expect(handler2).not.toHaveBeenCalled();
    });
});

// ============================================================================
// Promise Utilities Tests
// ============================================================================

describe('Promise Utilities', () => {
    
    describe('Promise.allSettled variant', () => {
        test('should handle all resolved', async () => {
            const promises = [
                Promise.resolve(1),
                Promise.resolve(2),
                Promise.resolve(3)
            ];
            
            const results = await Promise.allSettled(promises);
            
            expect(results).toEqual([
                { status: 'fulfilled', value: 1 },
                { status: 'fulfilled', value: 2 },
                { status: 'fulfilled', value: 3 }
            ]);
        });

        test('should handle mixed resolve/reject', async () => {
            const promises = [
                Promise.resolve(1),
                Promise.reject(new Error('fail')),
                Promise.resolve(3)
            ];
            
            const results = await Promise.allSettled(promises);
            
            expect(results[0]).toEqual({ status: 'fulfilled', value: 1 });
            expect(results[1].status).toBe('rejected');
            expect(results[1].reason).toBeInstanceOf(Error);
            expect(results[2]).toEqual({ status: 'fulfilled', value: 3 });
        });
    });

    describe('Promise.any', () => {
        test('should resolve with first successful promise', async () => {
            const promises = [
                Promise.reject(new Error('fail 1')),
                delay(50).then(() => 2),
                Promise.reject(new Error('fail 3'))
            ];
            
            const result = await Promise.any(promises);
            expect(result).toBe(2);
        });

        test('should reject if all promises reject', async () => {
            const promises = [
                Promise.reject(new Error('fail 1')),
                Promise.reject(new Error('fail 2')),
                Promise.reject(new Error('fail 3'))
            ];
            
            await expect(Promise.any(promises)).rejects.toThrow('All promises rejected');
        });
    });

    describe('Promise.race', () => {
        test('should resolve with fastest promise', async () => {
            const promises = [
                delay(50).then(() => 'slow'),
                delay(10).then(() => 'fast'),
                delay(30).then(() => 'medium')
            ];
            
            const result = await Promise.race(promises);
            expect(result).toBe('fast');
        });

        test('should reject with fastest rejection', async () => {
            const promises = [
                delay(50).then(() => 'slow'),
                delay(10).then(() => { throw new Error('fast error'); }),
                delay(30).then(() => 'medium')
            ];
            
            await expect(Promise.race(promises)).rejects.toThrow('fast error');
        });
    });
});

// ============================================================================
// Error Handling Tests
// ============================================================================

describe('Async Error Handling', () => {
    
    test('should handle unhandled promise rejections', (done) => {
        const handler = jest.fn();
        process.once('unhandledRejection', handler);
        
        // Create unhandled rejection
        Promise.reject(new Error('test rejection'));
        
        setTimeout(() => {
            expect(handler).toHaveBeenCalled();
            done();
        }, 10);
    });

    test('should handle multiple nested async errors', async () => {
        const asyncFn1 = async () => {
            await delay(10);
            throw new Error('error in asyncFn1');
        };
        
        const asyncFn2 = async () => {
            try {
                await asyncFn1();
            } catch (err) {
                throw new Error(`wrapped: ${err.message}`);
            }
        };
        
        await expect(asyncFn2()).rejects.toThrow('wrapped: error in asyncFn1');
    });

    test('should handle errors in promise chains', async () => {
        const result = await Promise.resolve(1)
            .then(x => x + 1)
            .then(x => { throw new Error('chain error'); })
            .then(x => x + 1)
            .catch(err => err.message);
        
        expect(result).toBe('chain error');
    });

    test('should handle finally after error', async () => {
        const finallySpy = jest.fn();
        
        await expect(Promise.reject(new Error('test'))
            .finally(finallySpy)
        ).rejects.toThrow('test');
        
        expect(finallySpy).toHaveBeenCalled();
    });
});

// ============================================================================
// Performance Tests
// ============================================================================

describe('Async Performance', () => {
    
    test('should measure async operation duration', async () => {
        const start = performance.now();
        await delay(100);
        const duration = performance.now() - start;
        
        expect(duration).toBeGreaterThanOrEqual(90);
        expect(duration).toBeLessThan(150);
    });

    test('should handle high concurrency', async () => {
        const tasks = Array(100).fill().map((_, i) => async () => {
            await delay(10);
            return i;
        });
        
        const start = performance.now();
        const results = await parallel(tasks, { concurrency: 20 });
        const duration = performance.now() - start;
        
        expect(results.length).toBe(100);
        // With concurrency 20, should take ~50ms (100/20 * 10ms)
        expect(duration).toBeLessThan(200);
    });

    test('should not block event loop', async () => {
        let intervalFired = false;
        
        const interval = setInterval(() => {
            intervalFired = true;
        }, 10);
        
        // Heavy async operation
        await parallel(Array(1000).fill().map(() => async () => {
            await delay(1);
        }), { concurrency: 100 });
        
        clearInterval(interval);
        expect(intervalFired).toBe(true);
    });
});