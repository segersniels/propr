import { useEffect, useState } from 'react';

export default function useLocalStorage<T>(
  key: string,
  initialValue: T
): [T, (value: T) => void, boolean] {
  const [isInitialized, setIsInitialized] = useState(false);
  const [storedValue, setStoredValue] = useState(initialValue);

  /**
   * Initialise the local state from localStorage
   */
  useEffect(() => {
    // Retrieve from localStorage
    const item = window.localStorage.getItem(key);
    if (item) {
      setStoredValue(JSON.parse(item));
    }

    // Notify callers that we tried fetching the storage value at least once
    setIsInitialized(true);
  }, [key]);

  /**
   * Listen for changes to the localStorage value
   */
  useEffect(() => {
    const eventName = `event:${key}`;
    const listener = () => {
      const item = window.localStorage.getItem(key);
      if (item) {
        setStoredValue(JSON.parse(item));
      }
    };

    window.addEventListener(eventName, listener);

    // Remove event listener on cleanup
    return () => {
      window.removeEventListener(eventName, listener);
    };
  }, [key]);

  const setValue = (value: T) => {
    // Save state
    setStoredValue(value);
    // Save to localStorage
    window.localStorage.setItem(key, JSON.stringify(value));
    // Notify listeners
    window.dispatchEvent(new Event(`event:${key}`));
  };

  return [storedValue, setValue, isInitialized];
}
