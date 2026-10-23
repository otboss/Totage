import { signal, effect, Signal } from '@angular/core';

export function debouncedSignal<T>(
  source: Signal<T>, 
  delayMs: number,
  callback?: () => void
): Signal<T> {
  // 1. Create the internal writable signal with the current source value
  const debounced = signal(source());

  // 2. Use an effect to react to source changes
  effect((onCleanup) => {
    const value = source(); // Track the source signal

    // 3. Set a timeout to update the debounced signal
    const timeout = setTimeout(() => {
      debounced.set(value);
      callback()
    }, delayMs);

    // 4. CRITICAL: Clear the previous timeout if the source changes again
    onCleanup(() => clearTimeout(timeout));
  }, { allowSignalWrites: true });

  // 5. Return as a read-only signal
  return debounced.asReadonly();
}
