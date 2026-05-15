# Thinker Queue

The worker uses two queues:

- `thinker` for normal processing
- `thinker-fixer` for exhausted items

## Rule

An item is eligible when `think_result IS NULL`, its error count is inside the queue's range, and it is unlocked.

## Queue Windows

- `thinker`: `think_error_count` `0..3`, ordered by `updated_at ASC`
- `thinker-fixer`: `think_error_count` `4..6`, ordered by `created_at ASC`
- `7+`: dead-letter

## Algorithm

1. Load up to `N` eligible items for the selected queue.
2. Skip rows locked by another worker.
3. Mark the selected rows with a temporary lock window.
4. Run the model.
5. On success, store `think_result` and reset `think_error_count`.
6. On failure, store the error message and increment `think_error_count`.
7. Re-queue the item only while it still fits the queue window.

## Why this works

- Fresh items are processed first.
- Retryable failures stay in the thinker lane.
- Exhausted failures move to the fixer lane.
- The fixer starts with the oldest matching item by creation time.
- There is no gap where an item can fail once or twice and then disappear.

## Practical effect

This keeps the queues bounded and makes retry behavior deterministic:

- `0` errors: normal first pass
- `1..3` errors: thinker retryable
- `4..6` errors: fixer retryable
- `7+` errors: stop retrying
